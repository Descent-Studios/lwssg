package main

import (
	"reflect"
	"testing"
)

func TestHeader(t *testing.T) {
	text := "# Header"
	result, err := parseMarkdown(text)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []Node{
		Node{
			nodeType: NodeH1,
			leaves: []Node{
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": "Header"},
				},
			},
			properties: nil,
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParagraph(t *testing.T) {
	text := "hello i am a paragraph. Hi!"
	result, err := parseMarkdown(text)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []Node{
		Node{
			nodeType: NodeParagraph,
			leaves: []Node{
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": "hello i am a paragraph. Hi!"},
				},
			},
			properties: nil,
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMultilineParagraph(t *testing.T) {
	text := "hello i am a paragraph. Hi!\nThis is the same paragraph"
	result, err := parseMarkdown(text)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []Node{
		Node{
			nodeType: NodeParagraph,
			leaves: []Node{
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": "hello i am a paragraph. Hi!"},
				},
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": " This is the same paragraph"},
				},
			},
			properties: nil,
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestMultipleParagraphs(t *testing.T) {
	text := "hello i am a paragraph. Hi!\n\nThis is a new paragraph!"
	result, err := parseMarkdown(text)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []Node{
		Node{
			nodeType: NodeParagraph,
			leaves: []Node{
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": "hello i am a paragraph. Hi!"},
				},
			},
			properties: nil,
		},
		Node{
			nodeType: NodeParagraph,
			leaves: []Node{
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": "This is a new paragraph!"},
				},
			},
			properties: nil,
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFormatting(t *testing.T) {
	text := "this is a paragraph *bold* _italics_"
	result, err := parseMarkdown(text)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []Node{
		Node{
			nodeType: NodeParagraph,
			leaves: []Node{
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": "this is a paragraph "},
				},
				Node{
					nodeType: NodeBold,
					leaves: []Node{
						Node{
							nodeType:   NodeText,
							leaves:     nil,
							properties: map[string]string{"content": "bold"},
						},
					},
					properties: nil,
				},
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": " "},
				},
				Node{
					nodeType: NodeItalic,
					leaves: []Node{
						Node{
							nodeType:   NodeText,
							leaves:     nil,
							properties: map[string]string{"content": "italics"},
						},
					},
					properties: nil,
				},
			},
			properties: nil,
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestLexer1(t *testing.T) {
	text := "# Header\nparagraph *with* **formatting** _alt_ __alt2__"
	result, err := lexMarkdown(text)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []MarkdownToken{
		{tokenType: MDTokenHeader1, text: "#"},
		{tokenType: MDTokenText, text: "Header"},
		{tokenType: MDTokenNewLine, text: "\n"},
		{tokenType: MDTokenText, text: "paragraph "},
		{tokenType: MDTokenItalic, text: "*"},
		{tokenType: MDTokenText, text: "with"},
		{tokenType: MDTokenItalic, text: "*"},
		{tokenType: MDTokenText, text: " "},
		{tokenType: MDTokenBold, text: "**"},
		{tokenType: MDTokenText, text: "formatting"},
		{tokenType: MDTokenBold, text: "**"},
		{tokenType: MDTokenText, text: " "},
		{tokenType: MDTokenItalic, text: "_"},
		{tokenType: MDTokenText, text: "alt"},
		{tokenType: MDTokenItalic, text: "_"},
		{tokenType: MDTokenText, text: " "},
		{tokenType: MDTokenBold, text: "__"},
		{tokenType: MDTokenText, text: "alt2"},
		{tokenType: MDTokenBold, text: "__"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, Got %v", expected, result)
	}
}

func TestLexer2(t *testing.T) {
	text := "#Not header\n*nest **formatting***\nthis * is just text"
	result, err := lexMarkdown(text)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []MarkdownToken{
		{tokenType: MDTokenText, text: "#Not header"},
		{tokenType: MDTokenNewLine, text: "\n"},
		{tokenType: MDTokenItalic, text: "*"},
		{tokenType: MDTokenText, text: "nest "},
		{tokenType: MDTokenBold, text: "**"},
		{tokenType: MDTokenText, text: "formatting"},
		{tokenType: MDTokenBold, text: "**"},
		{tokenType: MDTokenItalic, text: "*"},
		{tokenType: MDTokenNewLine, text: "\n"},
		{tokenType: MDTokenText, text: "this * is just text"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, Got %v", expected, result)
	}
}
