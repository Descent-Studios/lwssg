package main

import (
	"fmt"
	"os"
	"strings"
)

type NodeType int

const (
	NodeText = iota
	NodeBold
	NodeItalic
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

	fmt.Println(str)
}

func hello() string {
	return "hi"
}

func parseMarkdown(input string) ([]Node, error) {
	output := make([]Node, 0)
	appendToPrev := false
	for line := range strings.Lines(input) {
		text := strings.TrimSpace(line)

		if len(text) == 0 {
			appendToPrev = false
			continue
		}
		if !appendToPrev {
			parsed, err := parseLine(text)

			if err != nil {
				return nil, err
			}

			output = append(output, parsed)
		} else {
			lastNode := len(output) - 1
			output[lastNode].leaves =
				append(output[lastNode].leaves,
					Node{
						nodeType:   NodeText,
						leaves:     nil,
						properties: map[string]string{"content": " " + text},
					},
				)
		}

		appendToPrev = true
	}
	return output, nil
}

func parseLine(line string) (Node, error) {
	// get header amt
	headerLevel := 0
	i := 0
	for i = 0; line[i] == '#'; i++ {
		headerLevel++
	}
	if line[i] != ' ' {
		headerLevel = 0
	}

	var nodeType NodeType = NodeParagraph
	if headerLevel > 0 {
		nodeType = NodeType(NodeH1 + headerLevel - 1)
		headerLevel++
	}

	trimmedLine := line[headerLevel:]

	return Node{
		nodeType: nodeType,
		leaves: []Node{
			Node{
				nodeType:   NodeText,
				leaves:     nil,
				properties: map[string]string{"content": trimmedLine},
			},
		},
		properties: nil,
	}, nil
}
