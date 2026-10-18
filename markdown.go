package main

import (
	"regexp"
)

func parseMarkdown(input string) ([]Node, error) {
	lexed, err := lexMarkdown(input)
	if err != nil {
		return nil, err
	}

	output, err := parseMarkdownSub(lexed, true)

	if err != nil {
		return nil, err
	}
	return output, nil
}

type MDParserState int

const (
	MDPParagraph = iota
)

func parseMarkdownSub(lexed []MarkdownToken, toplevel bool) ([]Node, error) {
	output := make([]Node, 0)

	var state MDParserState = MDParserState(MDPParagraph)

	for i := range lexed {
		lexeme := lexed[i]

		switch state {
		case MDPParagraph:

		}
	}

	return output, nil
}

func findParaTerminator(tokens []MarkdownToken) int {
	for i := 0; i < len(tokens)-1; i++ {
		if tokens[i].tokenType == MDTokenNewLine && tokens[i+1].tokenType == MDTokenNewLine {
			return i
		}
	}
	return len(tokens)
}

func findBoldItalTerminator(tokens []MarkdownToken) int {
	const (
		Italic = iota + 1
		Bold
		BoldItalic
	)

	stack := []int{}
	for i := range tokens {
		token := tokens[i]
		stack_top := -1
		if len(stack) > 0 {
			stack_top = stack[len(stack)-1]
		}
		switch token.tokenType {
		case MDTokenBold:
			if stack_top == Bold {
				stack = stack[:len(stack)-1]
			} else if stack_top == BoldItalic {
				stack[len(stack)-1] = Italic
			} else {
				stack = append(stack, Bold)
			}
		case MDTokenItalic:
			if stack_top == Italic {
				stack = stack[:len(stack)-1]
			} else if stack_top == BoldItalic {
				stack[len(stack)-1] = Bold
			} else {
				stack = append(stack, Italic)
			}
		case MDTokenBoldItalic:
			if stack_top == BoldItalic {
				stack = stack[:len(stack)-1]
			} else if len(stack) > 1 && stack_top+stack[len(stack)-2] == BoldItalic {
				stack = stack[:len(stack)-2]
			} else {
				stack = append(stack, BoldItalic)
			}
		}
		if len(stack) == 0 {
			return i
		}
	}
	return -1
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

	MDTokenEscape
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

	regexp.MustCompile("(?m)(?:[^\\s`]|^)(`)[^`]|[^`](`)(?:[^\\s`]|$)"),
	nil,

	nil,
	regexp.MustCompile(`(?m)^((- )|(\* ))`),

	nil,

	nil,

	nil,
	nil,

	nil,
	nil,

	nil,

	regexp.MustCompile(`(\\)[-\\_#*\x60]`),
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
