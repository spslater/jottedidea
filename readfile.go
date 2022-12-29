package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "regexp"
)

func blockfile(filename string) [][]string{
    var blocks [][]string

    file, _ := os.Open(filename)
    defer file.Close()

    var fragments []string
    var scanner *bufio.Scanner = bufio.NewScanner(file)
    for scanner.Scan() {
        var line string = scanner.Text()
        if len(line) == 0 && len(fragments) == 0 {
            continue
        }
        if len(line) == 0 {
            blocks = append(blocks, fragments)
            fragments = []string{}
            continue
        }
        fragments = append(fragments, line)
    }
    if len(fragments) > 0 {
        blocks = append(blocks, fragments)
    }

    return blocks
}

func writefile(filename string, output string) {
    file, _ := os.Create(filename)
    defer file.Close()

    writer := bufio.NewWriter(file)
    writer.WriteString(output)
    writer.Flush()
}

func getline(reader *bufio.Reader, msg string) (string, error) {
    fmt.Print(msg)
    input, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("err:", err)
        return "", err
    }

    return strings.TrimSpace(input), nil
}

func getmultiline(reader *bufio.Reader, msg string) (string, error) {
    var tmp []string
    var skip bool
    fmt.Print(msg)
    for {
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("err:", err)
            return "", err
        }
        input = strings.TrimSpace(input)
        if input == "" && skip {
            break
        } else if input == "" && !skip {
            skip = true
        } else if input != "" && skip {
            skip = false
        }
        tmp = append(tmp, input)
    }

    return strings.Join(tmp, "\n"), nil
}

type Jot struct {
    Title string
    Desc string
    Section string
    Tags string
    Doc string
}

func savejot(jot Jot) {
    if jot.Title == "" { return }
    var re *regexp.Regexp = regexp.MustCompile(`\s+`)
    var filename string = re.ReplaceAllString(jot.Title, "_") + ".ji"

    if _, err := os.Stat(filename); err == nil {
        fmt.Println("jot name already exists")
    }

    var sb strings.Builder
    sb.WriteString("~~~\n")
    sb.WriteString(fmt.Sprintf("title: %s\n", jot.Title))
    sb.WriteString(fmt.Sprintf("description: %s\n", jot.Desc))
    sb.WriteString(fmt.Sprintf("section: %s\n", jot.Section))
    sb.WriteString(fmt.Sprintf("tags: %s\n", jot.Tags))
    sb.WriteString("~~~\n\n")
    sb.WriteString(jot.Doc)

    writefile(filename, sb.String())
}

func newjot() {
    var title, desc, sect, tags, doc string
    var err error

    var reader *bufio.Reader = bufio.NewReader(os.Stdin)
    title, err = getline(reader, "title: ")
    if err != nil { return }

    desc, err = getline(reader, "description: ")
    if err != nil { return }

    sect, err = getline(reader, "section: ")
    if err != nil { return }

    tags, err = getline(reader, "tags: ")
    if err != nil { return }

    doc, err = getmultiline(reader, "doc:\n")
    if err != nil { return }

    savejot(Jot{
        Title: title,
        Desc: desc,
        Section: sect,
        Tags: tags,
        Doc: doc,
    })
}