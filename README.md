# Pokédex CLI 🎮

A command-line Pokédex built in Go that lets you explore the Pokémon world, catch Pokémon, and manage your collection. Uses the [PokéAPI](https://pokeapi.co/) as its data source.

## Features

- Browse location areas in the Pokémon world
- Explore areas to find wild Pokémon
- Catch Pokémon with a chance-based mechanic
- Inspect your caught Pokémon's stats and types
- Built-in cache to avoid redundant API calls

## Project Structure

```
pokedex/
├── main.go
└── internal/
    ├── pokecache/
    │   └── cache.go       
    └── repl/
        └── repl.go        
```

## Getting Started

### Prerequisites

- Go 1.22+

### Installation

```bash
git clone https://github.com/Brenda2310/pokedex
cd pokedex
go build -o pokedex
./pokedex
```

## Commands

| Command | Description |
|---|---|
| `help` | Display available commands |
| `map` | Show the next 20 location areas |
| `mapb` | Show the previous 20 location areas |
| `explore <area>` | List all Pokémon in a location area |
| `catch <pokemon>` | Try to catch a Pokémon |
| `inspect <pokemon>` | View stats of a caught Pokémon |
| `pokedex` | List all your caught Pokémon |
| `exit` | Exit the Pokédex |

## Usage Example

```
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - shellos
 - gyarados
 - floatzel
 - quagsire

Pokedex > catch gyarados
Throwing a Pokeball at gyarados...
gyarados escaped!

Pokedex > catch shellos
Throwing a Pokeball at shellos...
shellos was caught!

Pokedex > inspect shellos
Name: shellos
Height: 3
Weight: 63
Stats:
  -hp: 76
  -attack: 48
  -defense: 48
  -special-attack: 57
  -special-defense: 62
  -speed: 34
Types:
  - water

Pokedex > pokedex
- shellos
```

## How Catching Works

Each Pokémon has a `base_experience` value. When you throw a Pokéball, a random number between 0 and 400 is generated. If the random number is higher than the Pokémon's base experience, you catch it — otherwise it escapes. Higher-level Pokémon are harder to catch.

## Cache

API responses are cached in memory with a 5-minute TTL to reduce redundant network requests. The cache runs a background cleanup loop that removes expired entries automatically.

## Running Tests

```bash
go test ./...
```
