package main

import (
	"reflect"
	"testing"
)

func TestHello(t *testing.T) {
	if hello() != "hi" {
		t.Errorf("L")
	}
}

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
					leaves:     []Node{},
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
