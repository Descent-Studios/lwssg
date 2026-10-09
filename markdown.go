package main

import "strings"

func parseMarkdown(input string) ([]Node, error) {
	_, err := lexMarkdown(input)
	if err != nil {
		return nil, err
	}
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

	leaves := []Node{Node{
		nodeType:   NodeText,
		leaves:     nil,
		properties: map[string]string{"content": trimmedLine},
	}}

	levels := []byte{}
	for i := 0; i < len(trimmedLine); i++ {
		char := trimmedLine[i]

		if char == '*' || char == '_' {
			if len(levels) != 0 && levels[len(levels)-1] == char {
				levels = levels[:len(levels)-1]
				continue
			}
			i++
			char := trimmedLine[i]
			if char == ' ' {
				continue
			}
			levels = append(levels, trimmedLine[i-1])
		}
	}

	return Node{
		nodeType:   nodeType,
		leaves:     leaves,
		properties: nil,
	}, nil
}

type MarkdownTokenType int

const (
	MDTokenText = iota
	MDTokenNewLine
	MDTokenItalic
	MDTokenBold
	MDTokenPre
	MDTokenCode
	MDTokenNumberList
	MDTokenUndorderedList
	MDTokenHorizontalRule
	MDTokenBlockquote
	MDTokenLinkBOpen
	MDTokenLinkBClose
	MDTokenLinkAOpen
	MDTokenLinkAClose
	MDTokenImageMark
	MDTokenHeader1
	MDTokenHeader2
	MDTokenHeader3
	MDTokenHeader4
	MDTokenHeader5
	MDTokenHeader6
)

type MarkdownToken struct {
	tokenType MarkdownTokenType
	text      string
}

func lexMarkdown(text string) ([]MarkdownToken, error) {
	return nil, nil
}
