package main

import (
	"fmt"
	"os"
)

type NodeType int

const (
	NodeText = iota
	NodeBold
	NodeItalic
	NodePre
	NodeCode
	NodeLink
	NodeImage
	NodeQuote
	NodeUL
	NodeOL
	NodeLI
	NodeHR
	NodeParagraph
	NodeH1
	NodeH2
	NodeH3
	NodeH4
	NodeH5
	NodeH6
)

type Node struct {
	nodeType   NodeType
	leaves     []Node
	properties map[string]string
}

func main() {
	file, err := os.ReadFile("test.md")
	if err != nil {
		panic(err)
	}

	str := string(file)

	fmt.Println(parseMarkdown(str))
}

func hello() string {
	return "hi"
}
