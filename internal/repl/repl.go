package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Brenda2310/pokedex/internal/pokecache"
)

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *Config, args ...string) error
}

type Config struct {
	Next          *string
	Previous      *string
	Cache         *pokecache.Cache
	CaughtPokemon map[string]Pokemon
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Minute)
	cfg := &Config{
		Cache:         cache,
		CaughtPokemon: make(map[string]Pokemon),
	}

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)

		if len(words) == 0 {
			continue
		}

		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		commandName := words[0]

		registry := getCommands()

		command, ok := registry[commandName]

		if ok {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <location_area_name>",
			description: "Displays a list of all the Pokémon located there",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Catch a Pokemon and add it to the user's Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "It takes the name of a Pokemon and prints the name, height, weight, stats and type(s) of the Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of all the names of the Pokemon the user has caught",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func cleanInput(text string) []string {
	edited := strings.ToLower(text)
	result := strings.Fields(edited)
	return result
}

func commandExit(cfg *Config, args ...string) error {
	println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args ...string) error {
	println("Welcome to the Pokedex!")
	println("Usage:")
	println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	println()
	return nil
}

func commandMap(cfg *Config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	res, err := fetchLocationAreas(url, cfg.Cache)
	if err != nil {
		return err
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous

	for _, area := range res.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(cfg *Config, args ...string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url := *cfg.Previous

	res, err := fetchLocationAreas(url, cfg.Cache)
	if err != nil {
		return err
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous

	for _, area := range res.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a location area name")
	}

	areaName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	fmt.Printf("Exploring %s...\n", areaName)

	var locationArea LocationArea

	if val, ok := cfg.Cache.Get(url); ok {
		err := json.Unmarshal(val, &locationArea)
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cfg.Cache.Add(url, body)

		err = json.Unmarshal(body, &locationArea)
		if err != nil {
			return err
		}
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}

	pokemon := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	var poke Pokemon

	if val, ok := cfg.Cache.Get(url); ok {
		err := json.Unmarshal(val, &poke)
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("response failed with status code: %d", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cfg.Cache.Add(url, body)

		err = json.Unmarshal(body, &poke)
		if err != nil {
			return err
		}
	}

	resistencia := rand.Intn(400)

	if resistencia > poke.BaseExperience {
		fmt.Printf("%s was caught!\n", poke.Name)
		cfg.CaughtPokemon[poke.Name] = poke
	} else {
		fmt.Printf("%s escaped!\n", poke.Name)
	}

	return nil
}

func commandInspect(cfg *Config, args ...string) error {
	pokemons := cfg.CaughtPokemon
	if len(pokemons) == 0 {
		return fmt.Errorf("There is not pokemons in your pokedex")
	}

	if len(args) != 1 {
		return fmt.Errorf("You must provide a pokemon name")
	}

	inspectPoke := args[0]

	poke, ok := pokemons[inspectPoke]

	if !ok {
		return fmt.Errorf("You have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", poke.Name)
	fmt.Printf("Height: %d\n", poke.Height)
	fmt.Printf("Weight: %d\n", poke.Weight)

	fmt.Println("Stats:")
	for _, stat := range poke.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range poke.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *Config, args ...string) error {
	pokemons := cfg.CaughtPokemon
	if len(pokemons) == 0 {
		return fmt.Errorf("There is not pokemons in your pokedex")
	}

	for poke := range pokemons {
		fmt.Printf("- %s\n", poke)
	}

	return nil
}

func fetchLocationAreas(url string, cache *pokecache.Cache) (LocationAreasResponse, error) {
	var locationAreas LocationAreasResponse

	if val, ok := cache.Get(url); ok {
		err := json.Unmarshal(val, &locationAreas)
		if err != nil {
			return locationAreas, err
		}
		return locationAreas, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return locationAreas, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return locationAreas, fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locationAreas, err
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return locationAreas, err
	}

	return locationAreas, nil
}
