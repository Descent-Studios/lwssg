package main

import (
	"fmt"
	"os"
)

type NodeType int

const (
	NodeText = iota
	NodeHeader
	NodeBold
	NodeItalic
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

	fmt.Println(str)
}

func hello() string {
	return "hi"
}
