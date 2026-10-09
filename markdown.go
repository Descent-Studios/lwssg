package main

import (
	"regexp"
	"strings"
)

func parseMarkdown(input string) ([]Node, error) {
	// _, err := lexMarkdown(input)
	// if err != nil {
	// 	return nil, err
	// }
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

// +++ Lexer
const (
	MDTokenHeader1 = iota
	MDTokenHeader2
	MDTokenHeader3
	MDTokenHeader4
	MDTokenHeader5
	MDTokenHeader6

	MDTokenItalic
	MDTokenBold
	MDTokenBoldItalic

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

	MDTokenNewLine
	MDTokenText
)

type MarkdownToken struct {
	tokenType MarkdownTokenType
	text      string
}

var mdTokenMatcher = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^(# )`),
	regexp.MustCompile(`(?m)^(## )`),
	regexp.MustCompile(`(?m)^(### )`),
	regexp.MustCompile(`(?m)^(#### )`),
	regexp.MustCompile(`(?m)^(##### )`),
	regexp.MustCompile(`(?m)^(###### )`),

	// markdown is evil btw
	regexp.MustCompile(`(?m)(?:[^\s*]|^)(\*)[^*]|(?:[^\s_]|^)(_)[^_]|[^*](\*)(?:[^\s*]|$)|[^_](_)(?:[^\s_]|$)`),
	regexp.MustCompile(`(?m)(?:[^\s*]|^)(\*\*)[^*]|(?:[^\s_]|^)(__)[^_]|[^*](\*\*)(?:[^\s*]|$)|[^_](__)(?:[^\s_]|$)`),
	regexp.MustCompile(`(?m)(?:[^\s*]|^)(\*\*\*)[^*]|(?:[^\s_]|^)(___)[^_]|[^*](\*\*\*)(?:[^\s*]|$)|[^_](___)(?:[^\s_]|$)`),

	nil,
	nil,

	nil,
	nil,

	nil,

	nil,

	nil,
	nil,

	nil,
	nil,

	nil,

	regexp.MustCompile(`(\n)`),
}

func lexMarkdown(text string) ([]MarkdownToken, error) {

	babo := [][][]int{}

	// run through every single matcher
	for m_i := range mdTokenMatcher {
		matcher := mdTokenMatcher[m_i]
		if matcher == nil {
			babo = append(babo, nil)
			continue
		}
		matched := matcher.FindAllStringSubmatchIndex(text, -1)
		babo = append(babo, matched)
	}

	beeba := [][]MarkdownToken{}

	for i := range babo {
		holder := []MarkdownToken{}
		for j := range babo[i] {
			for k := 2; k < len(babo[i][j]); k += 2 {
				rangestart := babo[i][j][k]
				rangeend := babo[i][j][k+1]
				if rangestart == -1 || rangeend == -1 {
					continue
				}
				babo[i][j] = []int{rangestart, rangeend}
				holder = append(holder, MarkdownToken{
					tokenType: MarkdownTokenType(i),
					text:      text[rangestart:rangeend],
				})
				break
			}
		}
		beeba = append(beeba, holder)
	}

	output := []MarkdownToken{}
	builderstring := ""

	for i := 0; i < len(text); {
		x, y := findStartIndex(babo, i)
		if x != -1 {
			if len(builderstring) > 0 {
				output = append(output, MarkdownToken{
					tokenType: MDTokenText,
					text:      builderstring,
				})
				builderstring = ""
			}
			output = append(output, beeba[x][y])
			i = babo[x][y][1]
		} else {
			builderstring += string(text[i])
			i++
		}
	}

	if len(builderstring) > 0 {
		output = append(output, MarkdownToken{
			tokenType: MDTokenText,
			text:      builderstring,
		})
	}

	return output, nil
}

func findStartIndex(babo [][][]int, index int) (int, int) {
	for i := range babo {
		for j := range babo[i] {
			if babo[i][j][0] == index {
				return i, j
			}
		}
	}

	return -1, -1
}

// ---
