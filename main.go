package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/YoavIsaacs/pokadexcli/internal/pokecache"
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
	config := new(Config)
	config.Cache = *pokecache.NewCache(10 * time.Second)
	config.PokemonCache = *pokecache.NewCache(10 * time.Second)
	config.next = "https://pokeapi.co/api/v2/location-area/"
	config.Pokedex.Pokemon = []string{}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		split := strings.Fields(input)
		command, exists := Commands[split[0]]
		if len(split) == 1 {
			split = append(split, "")
		}

		if !exists {
			fmt.Print("Unknown command\n")
			continue
		}
		if err := command.callback(config, split[1]); err != nil {
			fmt.Print(fmt.Errorf("error executing command %s: %v", command.name, err))
		}
	}
}
