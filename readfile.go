package main

import (
    "os"
    "bufio"
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