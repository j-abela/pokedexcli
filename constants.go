package main

const baseURL = "https://pokeapi.co/api/v2"

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// contains urls to paginate through locations
type config struct {
	Next     string
	Previous string
}

type LocationArea struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
