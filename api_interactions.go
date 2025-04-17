package main

import (
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
