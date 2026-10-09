package main

import (
	"fmt"
	"os"
	"strings"
)

type NodeType int

const (
	NodeText = iota
	NodeH1
	NodeH2
	NodeH3
	NodeH4
	NodeH5
	NodeH6
	NodeBold
	NodeItalic
	NodeParagraph
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

func parseMarkdown(input string) ([]Node, error) {
	output := make([]Node, 0)
	for line := range strings.Lines(input) {
		output = append(output, Node{
			nodeType: NodeParagraph,
			leaves: []Node{
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": line},
				},
			},
			properties: nil,
		})
	}
	return output, nil
}
