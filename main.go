package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.ReadFile("test.md")
	if err != nil {
		panic(err)
	}

	str := string(file)

	fmt.Println(str)
}
