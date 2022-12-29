package main

import (
    "fmt"
    // "os"
    "strings"
)

func convert(args []string) {
    for _, filename := range args {
        var blocks [][]string = blockfile(filename)
        var head, body string
        for _, block := range blocks {
            var btype BlockType = getBlockType(block[0])
            switch btype {
                case BlockHeader:
                    th, tb := doBlockHeader(block)
                    head += th + "\n"
                    body = tb + "\n" + body
                case BlockList:
                    body += doBlockList(block) + "\n"
                case BlockCode:
                    body += doBlockCode(block) + "\n"
                case BlockQuote:
                    body += doBlockQuote(block) + "\n"
                case BlockText:
                    body += doBlockText(block) + "\n"
                default:
                    fmt.Println("unknown block type", btype, block)
            }
        }
        var output string = "<head>\n" + head + "</head>\n<body>\n" + body + "</body>"
        var outfile string = strings.ReplaceAll(filename, "ji", "html")
        writefile(outfile, output)
    }
}
