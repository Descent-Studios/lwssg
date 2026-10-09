package main

import "strings"

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
	iterator := 0
	for iterator = 0; line[iterator] == '#'; iterator++ {
		headerLevel++
	}
	if line[iterator] != ' ' {
		headerLevel = 0
	}

	var nodeType NodeType = NodeParagraph
	if headerLevel > 0 {
		nodeType = NodeType(NodeH1 + headerLevel - 1)
		headerLevel++
	}

	trimmedLine := line[headerLevel:]

	leaves := []Node{}

	prevLevels := []int
	for i := 0; i < len(trimmedLine); i++ {
		char := trimmedLine[i]

		if char == '*' || char == '_' {
			i++
			char := trimmedLine[i]
			if char == '*' || char == '_' {
			} else if char == ' ' {
				continue
			} else {

			}
		}
	}

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
