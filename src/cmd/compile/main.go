package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"unicode/utf8"
)

var (
	filename = flag.String("f", "", "file name")
)

func showUsage() {
	fmt.Printf("Usage:%v -f FILENAME\n", os.Args[0])
}
func main() {
	flag.Parse()
	if *filename == "" {
		showUsage()
		return
	}

	in := bufio.NewReader(os.Stdin)
	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			log.Fatalf("WriteString: %s", err)
		}
		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}
		exprParse(&exprLex{line: line})
	}
}
