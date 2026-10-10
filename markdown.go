package main

import (
	"errors"
	"regexp"

	"descentdev.cc/lwssg/utils"
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

func parseMarkdownSub(lexed []MarkdownToken, toplevel bool) ([]Node, error) {
	output := make([]Node, 0)

	for token_i := 0; token_i < len(lexed); token_i++ {
		token := lexed[token_i]
		switch token.tokenType {
		case MDTokenHeader1:
			end_i := utils.Find(lexed, func(val MarkdownToken) bool { return val.tokenType == MDTokenNewLine })
			if end_i == -1 {
				end_i = len(lexed) - token_i
			}
			result, err := parseMarkdownSub(lexed[token_i+1:token_i+end_i], false)
			if err != nil {
				return nil, err
			}
			output = append(output,
				Node{
					nodeType:   NodeH1,
					leaves:     result,
					properties: nil,
				},
			)
			token_i += end_i - 1
		case MDTokenText:
			if !toplevel {
				output = append(output, Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": token.text},
				})
			} else {
				end_i := findParaTerminator(lexed[token_i:])
				result, err := parseMarkdownSub(lexed[token_i:token_i+end_i], false)
				if err != nil {
					return nil, err
				}

				output = append(output,
					Node{
						nodeType:   NodeParagraph,
						leaves:     result,
						properties: nil,
					},
				)
				token_i += end_i - 1
			}
		case MDTokenBold:
			end_i := findBoldItalTerminator(lexed[token_i:])
			if end_i == -1 {
				return nil, errors.New("Unclosed bold tag")
			}
			result, err := parseMarkdownSub(lexed[token_i+1:token_i+end_i], false)

			if err != nil {
				return nil, err
			}

			output = append(output,
				Node{
					nodeType:   NodeBold,
					leaves:     result,
					properties: nil,
				},
			)
			token_i += end_i
		case MDTokenItalic:
			end_i := findBoldItalTerminator(lexed[token_i:])
			if end_i == -1 {
				return nil, errors.New("Unclosed italic tag")
			}
			result, err := parseMarkdownSub(lexed[token_i+1:token_i+end_i], false)

			if err != nil {
				return nil, err
			}

			output = append(output,
				Node{
					nodeType:   NodeItalic,
					leaves:     result,
					properties: nil,
				},
			)

			token_i += end_i

		case MDTokenBoldItalic:
			end_i := findBoldItalTerminator(lexed[token_i:])
			if end_i == -1 {
				return nil, errors.New("Unclosed bold/italic")
			}
			result, err := parseMarkdownSub(lexed[token_i+1:token_i+end_i], false)

			if err != nil {
				return nil, err
			}

			output = append(output,
				Node{
					nodeType:   NodeItalic,
					leaves:     result,
					properties: nil,
				},
			)

			token_i += end_i
		case MDTokenNewLine:
			if !toplevel {
				if len(output) > 0 {
					output = append(output, Node{
						nodeType:   NodeSoftBreak,
						leaves:     nil,
						properties: nil,
					})
				} else {
					continue
				}
			} else {
				continue
			}
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
		Italic = iota
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
			} else if len(stack) > 2 && stack_top+stack[len(stack)-2] == BoldItalic {
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
