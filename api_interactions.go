package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/j-abela/pokedexcli/internal/pokecache"
)

var pokeCache *pokecache.Cache

func init() {
	// Create a cache with a 5-minute (or whatever duration you prefer) expiration
	pokeCache = pokecache.NewCache(5 * time.Minute)
}

func pokemonGET(url string) ([]byte, error) {
	// Check if the response is in the cache
	if cachedData, found := pokeCache.Get(url); found {
		return cachedData, nil
	}

	// Not in cache, perform the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %v", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Add the response to the cache
	pokeCache.Add(url, data)

	return data, nil
}
