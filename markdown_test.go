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
	text := "this is a paragraph *italics* **bold** _alt italics_ __alt bold__"
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
				Node{
					nodeType:   NodeText,
					leaves:     nil,
					properties: map[string]string{"content": " "},
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
							properties: map[string]string{"content": "alt italics"},
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
					nodeType: NodeBold,
					leaves: []Node{
						Node{
							nodeType:   NodeText,
							leaves:     nil,
							properties: map[string]string{"content": "alt bold"},
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
