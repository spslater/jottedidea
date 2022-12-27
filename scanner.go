package main

import (
    "regexp"
    "strings"
)

type BlockType int
const (
    BlockUnknown BlockType = iota // 0
    BlockHeader                   // 1
    BlockList                     // 2
    BlockCode                     // 3
    BlockQuote                    // 4
    BlockText                     // 5
)

type FragmentType int
const (
    FragmentUnknown FragmentType = iota // 0
    FragmentCode                        // 3
    FragmentImage                       // 1
    FragmentLink                        // 2
    FragmentBold                        // 4
    FragmentItalic                      // 5
    FragmentUnderline                   // 6
)

type ListType int
const (
    ListUnknown ListType = iota // 0
    ListOrdered                 // 1
    ListUnordered               // 2

)
type List struct {
    indent int
    typeId ListType
}

var tab int = 4

var listLookup map[ListType]string

var reBlockHeader *regexp.Regexp
var reBlockList *regexp.Regexp
var reBlockCode *regexp.Regexp
var reBlockQuote *regexp.Regexp

var reFragBlock *regexp.Regexp

var reFragMeta *regexp.Regexp
var reFragCode *regexp.Regexp
var reFragImage *regexp.Regexp
var reFragLink *regexp.Regexp
var reFragBold *regexp.Regexp
var reFragItalic *regexp.Regexp
var reFragUnderline *regexp.Regexp
var reFragListNum *regexp.Regexp
var reFragListUn *regexp.Regexp

func init() {
    reBlockHeader = regexp.MustCompile(`^~+$`)
    reBlockList = regexp.MustCompile(`^\s*([0-9]+\.|[\*\-\+])\s+.*`)
    reBlockCode = regexp.MustCompile(`^\s*`+"```")
    reBlockQuote = regexp.MustCompile(`^\s*>`)

    reFragBlock = regexp.MustCompile(`(>\s+)?(.*)$`) // indent, val

    reFragMeta = regexp.MustCompile(`^\s*(.+?):(.*)$`) // key, value
    reFragCode = regexp.MustCompile(`^(.*)`+"`"+`(.*?)`+"`"+`(.*)$`) // front, code, back
    reFragImage = regexp.MustCompile(`^(.*)\!\[(.*?)\]\((.*?)\)(.*)$`) // front, title, img, back
    reFragLink = regexp.MustCompile(`^(.*)\[(.*?)\]\((.*?)\)(.*)$`) // front, title, src, back
    reFragBold = regexp.MustCompile(`^(.*)\*\*(.*?)\*\*(.*)$`) // front, bold, back
    reFragItalic = regexp.MustCompile(`^(.*)\*(.*?)\*(.*)$`) // front, italic, back
    reFragUnderline = regexp.MustCompile(`^(.*)_(.*?)_(.*)$`) // front, under, back

    reFragListNum = regexp.MustCompile(`^(\s*)[0-9]+\.\s+(.*)$`) // indent, value
    reFragListUn = regexp.MustCompile(`^(\s*)[\*\-\+]\s+(.*)$`) // indent, value

    listLookup = map[ListType]string{ListOrdered: "ol", ListUnordered: "ul"}
}

func getBlockType(line string) BlockType {
    if reBlockHeader.MatchString(line) {
        return BlockHeader
    }
    if reBlockList.MatchString(line) {
        return BlockList
    }
    if reBlockCode.MatchString(line) {
        return BlockCode
    }
    if reBlockQuote.MatchString(line) {
        return BlockQuote
    }
    return BlockText
}


func doBlockHeader(block []string) (string, string) { // map[string]string {
    var header map[string]string = make(map[string]string)
    for _, b := range block {
        if match := reFragMeta.FindStringSubmatch(b); match != nil {
            header[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
        }
    }

    var headout, bodyout string
    headout = "<title>" + header["title"] + "</title>\n"
    bodyout = "<h2>" + header["title"] + "</h2>\n"

    headout += "<meta name='description' content='" + header["description"] + "' />\n"
    bodyout += "<b>" + header["description"] + "</b><br />\n"

    headout += "<meta name='section' content='" + header["section"] + "' />\n"
    bodyout += "<b>Section</b>: " + strings.ReplaceAll(header["section"], "/", " > ") + "<br />\n"

    headout += "<meta name='tags' content='" + header["section"] + "' />\n"
    bodyout += "<b>Tags</b>: <i>" + header["tags"] + "</i><br />\n"

    bodyout += "<hr />\n"

    delete(header, "title")
    delete(header, "description")
    delete(header, "section")
    delete(header, "tags")

    for key, val := range header {
        headout += "<meta name='" + key + "' content='" + val + "' />\n"
    }

    return strings.TrimSpace(headout), strings.TrimSpace(bodyout)
}

func listEntry(line string) (List, string) {
    var match []string
    var typeId ListType
    if reFragListNum.MatchString(line) {
        match = reFragListNum.FindStringSubmatch(line)
        typeId = ListOrdered
    } else if reFragListUn.MatchString(line) {
        match = reFragListUn.FindStringSubmatch(line)
        typeId = ListUnordered
    }
    var tin int = len(strings.ReplaceAll(match[1], "\t", strings.Repeat(" ", tab)))
    var cur string = match[2]
    return List{indent: tin, typeId: typeId}, cur
}

func doBlockList(block []string) string {
    var tmp, cur string
    var lists []List
    var clist List
    for _, line := range block {
        clist, cur = listEntry(line)
        if tmp == "" && clist.typeId == ListOrdered {
            tmp = "<ol>"
        } else if tmp == "" && clist.typeId == ListUnordered {
            tmp = "<ul>"
        }
        var last int = len(lists)-1
        if last < 0 {
            tmp += "<li>" + cur + "</li>"
            lists = append(lists, clist)
        } else if lists[last].indent == clist.indent {
            if lists[last].typeId == clist.typeId {
                tmp += "<li>" + cur + "</li>"
            } else {
                tmp += "</" + listLookup[lists[last].typeId] + ">"
                tmp += "<" + listLookup[lists[last].typeId] + "><li>" + cur + "</li>"
                lists[last] = clist
            }
        } else if lists[last].indent < clist.indent {
            tmp += "<" + listLookup[clist.typeId] + "><li>" + cur + "</li>"
            lists = append(lists, clist)
        } else if lists[last].indent > clist.indent {
            for {
                tmp += "</" + listLookup[lists[last].typeId] + ">"
                // fmt.Println(lists)
                lists = lists[:last]
                last--
                // fmt.Println(lists, last)
                if last < 0 || lists[last].indent == clist.indent {
                    break
                }
            }
            if lists[last].typeId != clist.typeId {
                tmp += "</" + listLookup[lists[last].typeId] + ">"
                tmp += "<" + listLookup[clist.typeId] + ">"
                lists[last] = clist
            }
            tmp += "<li>" + cur + "</li>"
        }
    }

    for _, list := range lists {
        tmp += "</" + listLookup[list.typeId] + ">"
    }

    return tmp
}

func doBlockCode(block []string) string {
    var tmp string = strings.Join(block, "\n")
    return "<code><pre>" + tmp + "</pre></code>"
}

func doBlockQuote(block []string) string {
    var tmp string // = strings.Join(block, " ")
    var match []string
    for _, line := range block {
        match = reFragBlock.FindStringSubmatch(line)
        tmp += match[2] + " "
    }
    return "<blockquote>" + tmp + "</blockquote>"
}

// reFragCode
// reFragImage
// reFragLink
// reFragBold
// reFragItalic
// reFragUnderline

func doBlockText(block []string) string {
    var tmp string // = strings.Join(block, " ")
    for _, line := range block {
        for {
            match := reFragCode.FindStringSubmatch(line)
            if len(match) == 0 { break }
            front, code, back := match[1], match[2], match[3]
            line = front + " <code>" + code + "</code> " + back
        }
        for {
            match := reFragImage.FindStringSubmatch(line)
            if len(match) == 0 { break }
            front, title, img, back := match[1], match[2], match[3], match[4]
            line = front + " <img src='" + img + "' alt='" + title + "'></img> " + back
        }
        for {
            match := reFragLink.FindStringSubmatch(line)
            if len(match) == 0 { break }
            front, title, src, back := match[1], match[2], match[3], match[4]
            line = front + " <a href='" + src + "'>" + title + "</a> " + back
        }
        for {
            match := reFragBold.FindStringSubmatch(line)
            if len(match) == 0 { break }
            front, bold, back := match[1], match[2], match[3]
            line = front + " <b>" + bold + "</b> " + back
        }
        for {
            match := reFragItalic.FindStringSubmatch(line)
            if len(match) == 0 { break }
            front, italic, back := match[1], match[2], match[3]
            line = front + " <i>" + italic + "</i> " + back
        }
        for {
            match := reFragUnderline.FindStringSubmatch(line)
            if len(match) == 0 { break }
            front, under, back := match[1], match[2], match[3]
            line = front + " <u>" + under + "</u> " + back
        }
        tmp += line + " "
    }

    return "<p>" + tmp + "</p>"
}

