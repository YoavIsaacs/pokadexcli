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
	InitCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		command, exists := Commands[input]
		if !exists {
			fmt.Print("Unknown command\n")
			continue
		}
		if err := command.callback(); err != nil {
			fmt.Print(fmt.Errorf("error executing command %s: %v", command.name, err))
		}
	}
}
