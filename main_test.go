package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	if hello() != "hi" {
		t.Errorf("L")
	}
}
