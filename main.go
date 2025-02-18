package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	var ret []string
	text = strings.TrimSpace(text)
	ret = strings.Fields(text)
	return ret
}

func main() {
	p := cleanInput("poo")
	fmt.Print(p)
}
