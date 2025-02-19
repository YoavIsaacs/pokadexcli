package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/YoavIsaacs/pokadexcli/internal/pokecache"
)

type Config struct {
	next         string
	prev         string
	Cache        pokecache.Cache
	PokemonCache pokecache.Cache
	Pokedex      PlayerPokedex
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

type PlayerPokedex struct {
	Pokemon []string
}

func commandHelp(c *Config, _ string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("\nUsage:\n\n")
	for _, value := range Commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Println()
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
		description: "    Exit the Pokedex",
		callback:    commandExit,
	},
	"map": {
		name:        "map",
		description: "     Show the next 20 locations",
		callback:    commandMap,
	},
	"nmap": {
		name:        "nmap",
		description: "    Show the previous 20 locations",
		callback:    commandNmap,
	},
	"explore": {
		name:        "explore",
		description: " List all Pokemon in this area",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "   Attempt to catch a Pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: " Get the stats of a caught Pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: " List all caught Pokemon",
		callback:    commandPokedex,
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
		description: "    Displays a help message",
		callback:    commandHelp,
	}
}

func commandExplore(c *Config, loc string) error {
	data, ok := c.PokemonCache.Cache[loc]
	if ok {
		fmt.Println()
		for _, pokemon := range data.Val {
			fmt.Println(pokemon)
		}
	}
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

type pokemonChance struct {
	Chance int `json:"base_experience"`
}

func commandCatch(c *Config, name string) error {
	fmt.Println("Throwing a Pokeball at " + name + "...")
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("invalid response: %v\n", res.Status)
	}

	defer res.Body.Close()
	var chance pokemonChance
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&chance)
	if err != nil {
		return err
	}

	playerChance := rand.IntN(chance.Chance * 2)
	if playerChance > chance.Chance {
		fmt.Println(name + " was caught!")
		fmt.Println("You man now inspect it with the inspect command.")
		c.Pokedex.Pokemon = append(c.Pokedex.Pokemon, name)
	} else {
		fmt.Println(name + " escaped!")
	}
	return nil
}

type statInfo struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type typeInfo struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokemonInfoForInspect struct {
	Name   string     `json:"name"`
	Height int        `json:"height"`
	Weight int        `json:"weight"`
	Stats  []statInfo `json:"stats"`
	Types  []typeInfo `json:"types"`
}

func existsInPokedex(c *Config, name string) bool {
	for _, n := range c.Pokedex.Pokemon {
		if name == n {
			return true
		}
	}
	return false
}

func commandInspect(c *Config, name string) error {
	if !existsInPokedex(c, name) {
		fmt.Println("You have not caught that pokemon")
		return nil
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("invalid response: %v\n", res.Status)
	}

	defer res.Body.Close()
	var info PokemonInfoForInspect
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&info)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", info.Name)
	fmt.Printf("Height: %d\n", info.Height)
	fmt.Printf("Weight: %d\n", info.Weight)

	fmt.Println("Stats:")
	for _, stat := range info.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range info.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(c *Config, _ string) error {
	if len(c.Pokedex.Pokemon) == 0 {
		fmt.Println("Your Pokedex is empty...")
		return nil
	}
	fmt.Println("Your Pokemon:")
	for _, p := range c.Pokedex.Pokemon {
		fmt.Println("  - " + p)
	}
	return nil
}
