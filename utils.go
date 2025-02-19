package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	next string
	prev string
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, value := range Commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config) error
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

func commandMap(c *Config) error {
	res, err := http.Get(c.next)
	fmt.Print(c.next + "\n")
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
