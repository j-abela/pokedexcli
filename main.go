package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		input.Scan()
		inputText := input.Text()
		words := cleanInput(inputText)
		fmt.Printf("Your command was: %s\n", words[0])
	}
}
