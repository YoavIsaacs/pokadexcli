package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	var ret []string
	text = strings.TrimSpace(text)
	ret = strings.Fields(text)
	return ret
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokadex > ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.ToLower(input)
		first := strings.Fields(input)[0]
		fmt.Printf("Your command was: %s\n", first)
	}
}
