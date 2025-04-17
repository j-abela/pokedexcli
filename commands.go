package main

import (
	"fmt"
	"os"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations in the Pokémon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations in the Pokémon world",
			callback:    commandMapb,
		},
	}
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *config) error {
	var fullURL string

	if cfg.Next != "" {
		fullURL = cfg.Next
	} else {
		fullURL = baseURL + "/location-area/"
	}

	body, err := pokemonGET(fullURL)
	if err != nil {
		return err
	}

	unmarshalErr := mapHelper(body, cfg)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	return nil
}

func commandMapb(cfg *config) error {
	var fullURL string

	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		fullURL = cfg.Previous
	}

	body, err := pokemonGET(fullURL)
	if err != nil {
		return err
	}

	unmarshalErr := mapHelper(body, cfg)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	return nil
}
