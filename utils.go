package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YoavIsaacs/pokadexcli/internal/pokecache"
)

type Config struct {
	next         string
	prev         string
	Cache        pokecache.Cache
	PokemonCache pokecache.Cache
}

type PokemonInfo struct {
	Name string `json:"name"`
}

type PokemonEncounter struct {
	Pokemon PokemonInfo `json:"pokemon"`
}

type PokemonResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func commandHelp(c *Config, _ string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, value := range Commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandExit(c *Config, _ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config, loc string) error
}

var Commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"map": {
		name:        "map",
		description: "Show the next 20 locations",
		callback:    commandMap,
	},
	"nmap": {
		name:        "nmap",
		description: "Show the previous 20 locations",
		callback:    commandNmap,
	},
	"explore": {
		name:        "explore",
		description: "List all pokemon in this area",
		callback:    commandExplore,
	},
}

type MapInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type APIResponse struct {
	Count    int       `json:"count"`
	NextUrl  string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []MapInfo `json:"results"`
}

func commandMap(c *Config, _ string) error {
	data, ok := c.Cache.Cache[c.next]
	if ok {
		fmt.Println()
		for _, location := range data.Val {
			fmt.Println(location)
		}
	}

	res, err := http.Get(c.next)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("invalid response: %v\n", res.Status)
	}
	defer res.Body.Close()
	var apiResponse APIResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&apiResponse)
	if err != nil {
		return err
	}
	c.next = apiResponse.NextUrl
	c.prev = apiResponse.Previous

	fmt.Println()
	for _, location := range apiResponse.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandNmap(c *Config, _ string) error {
	if c.prev == "" {
		fmt.Println("At the beginning, no previous maps...")
		return nil
	}

	data, ok := c.Cache.Cache[c.prev]
	if ok {
		fmt.Println()
		for _, location := range data.Val {
			fmt.Println(location)
		}
	}

	res, err := http.Get(c.prev)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("invalid response: %v\n", res.Status)
	}
	defer res.Body.Close()
	var apiResponse APIResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&apiResponse)
	if err != nil {
		return err
	}
	c.next = apiResponse.NextUrl
	c.prev = apiResponse.Previous

	fmt.Println()
	for _, location := range apiResponse.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func InitCommands() {
	Commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
}

func commandExplore(c *Config, loc string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + loc
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("invalid response: %v\n", res.Status)
	}

	defer res.Body.Close()
	var apiResponse PokemonResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&apiResponse)
	if err != nil {
		return err
	}
	fmt.Println()
	for _, pokemons := range apiResponse.PokemonEncounters {
		fmt.Println(pokemons.Pokemon.Name)
	}
	return nil
}
