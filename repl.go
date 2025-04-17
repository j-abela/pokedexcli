package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func startRepl() {
	reader := bufio.NewScanner(os.Stdin)
	cfg := &config{
		Next:     "",
		Previous: "",
	}

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()
		inputText := reader.Text()

		words := cleanInput(inputText)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
