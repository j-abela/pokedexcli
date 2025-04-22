package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
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
		"explore": {
			name:        "explore",
			description: "Format: 'explore <area_name>'. Displays a list of all Pokémon in a given area.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Format: 'catch <pokemon>'. Tries to catch a Pokémon to add to your Pokédex",
			callback:    commandCatch,
		},
	}
}

func commandExit(cfg *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, params []string) error {
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

func commandMap(cfg *config, params []string) error {
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

func commandMapb(cfg *config, params []string) error {
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

func mapHelper(body []byte, cfg *config) error {
	locationArea := LocationArea{}
	unmarshalErr := json.Unmarshal(body, &locationArea)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	cfg.Next = locationArea.Next
	if locationArea.Previous != nil {
		cfg.Previous = *locationArea.Previous
	} else {
		cfg.Previous = ""
	}

	for _, area := range locationArea.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(cfg *config, params []string) error {
	fullURL := baseURL + "/location-area/" + params[0]

	body, err := pokemonGET(fullURL)
	if err != nil {
		return err
	}

	locationArea := LocationAreaDetailed{}
	unmarshalErr := json.Unmarshal(body, &locationArea)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, params []string) error {
	fullURL := baseURL + "/pokemon/" + params[0]

	body, err := pokemonGET(fullURL)
	if err != nil {
		return err
	}

	pokemon := Pokemon{}
	unmarshalErr := json.Unmarshal(body, &pokemon)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	name := pokemon.Name
	baseExperience := pokemon.BaseExperience
	randomChance := rnd.Intn(100)
	catchThreshold := 70 - (baseExperience / 20)

	// Ensure there's always a minimum chance to catch
	if catchThreshold < 25 {
		catchThreshold = 25
	}

	caught := false

	fmt.Printf("Throwing a Pokeball at %v...\n", name)
	if randomChance <= catchThreshold {
		fmt.Printf("%v was caught!\n", name)
		caught = true
	} else {
		fmt.Printf("%v escaped!\n", name)
	}

	if caught {

	}

	return nil
}
