package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func pokemonGET(url string) (marshaled []byte, err error) {
	result, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if result.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", result.StatusCode, body)
		return []byte{}, err
	}
	if err != nil {
		return []byte{}, err
	}
	return body, nil
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
