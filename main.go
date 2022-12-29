package main

import (
    "os"
)

func main() {
    if len(os.Args) > 1 {
        convert(os.Args[1:])
    } else {
        newjot()
    }
}
