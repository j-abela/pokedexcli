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
		"inspect": {
			name:        "inspect",
			description: "Format: 'inspect <pokemon>'. Displays information on caught Pokémon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays the name of all caught Pokémon",
			callback:    commandPokedex,
		},
	}
}

func commandExit(game *Game, cfg *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(game *Game, cfg *config, params []string) error {
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

func commandMap(game *Game, cfg *config, params []string) error {
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

func commandMapb(game *Game, cfg *config, params []string) error {
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

func commandExplore(game *Game, cfg *config, params []string) error {
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

func commandCatch(game *Game, cfg *config, params []string) error {
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
	catchThreshold := 50 - (baseExperience / 10)

	// Ensure there's always a minimum chance to catch
	if catchThreshold < 20 {
		catchThreshold = 20
	}

	caught := false

	fmt.Printf("Throwing a Pokeball at %v...\n", name)
	if randomChance <= catchThreshold {
		fmt.Printf("%v was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command")
		caught = true
	} else {
		fmt.Printf("%v escaped!\n", name)
	}

	if caught {
		game.Pokedex[pokemon.Name] = pokemon
	}
	return nil
}

func commandInspect(game *Game, cfg *config, params []string) error {
	pokemonName := params[0]
	pokemon, exists := game.Pokedex[pokemonName]

	if exists {
		fmt.Printf("Height: %v\nWeight: %v\nStats:\n", pokemon.Height, pokemon.Weight)
		for _, stat := range pokemon.Stats {
			fmt.Printf("  - %v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, pkType := range pokemon.Types {
			fmt.Printf("  - %v\n", pkType.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func commandPokedex(game *Game, cfg *config, params []string) error {
	if len(game.Pokedex) == 0 {
		fmt.Println("you haven't caught any pokemon yet")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range game.Pokedex {
		fmt.Printf(" - %v\n", name)
	}
	return nil
}
