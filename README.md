# PokédexCLI

A command-line Pokémon exploration tool built in Go that interfaces with the [PokéAPI](https://pokeapi.co/).

## Features

- **Map Navigation**: Explore different Pokémon habitats and locations
- **Pokémon Discovery**: Find which Pokémon live in each area
- **Catching System**: Try your luck at catching Pokémon with a probability-based system
- **Pokémon Collection**: Keep track of your caught Pokémon in your Pokédex
- **Detailed Information**: Inspect detailed stats for Pokémon you've caught
- **Caching System**: Efficient data retrieval with automatic cache expiration

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.20 or higher recommended)

## Installation

### Linux/macOS

```bash
# Clone the repository
git clone https://github.com/YoavIsaacs/pokadexcli.git

# Navigate to the project directory
cd pokadexcli

# Build the project
go build
```

### Windows

```powershell
# Clone the repository
git clone https://github.com/YoavIsaacs/pokadexcli.git

# Navigate to the project directory
cd pokadexcli

# Build the project
go build
```

## Usage

After building, run the executable:

```bash
./pokadexcli
```

You'll be greeted with a prompt:

```
Pokedex >
```

## Commands

| Command  | Description                           |
|----------|---------------------------------------|
| `help`   | Displays a help message               |
| `exit`   | Exit the Pokédex                      |
| `map`    | Show the next 20 locations            |
| `nmap`   | Show the previous 20 locations        |
| `explore`| List all Pokémon in a specific area   |
| `catch`  | Attempt to catch a Pokémon            |
| `inspect`| Get the stats of a caught Pokémon     |
| `pokedex`| List all caught Pokémon               |

### Command Examples

#### Exploring an Area

```
Pokedex > explore canalave-city-area

tentacool
tentacruel
...
```

#### Catching a Pokémon

```
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
You may now inspect it with the inspect command.
```

#### Inspecting a Caught Pokémon

```
Pokedex > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  - hp: 35
  - attack: 55
  - defense: 40
  - special-attack: 50
  - special-defense: 50
  - speed: 90
Types:
  - electric
```

## Architecture

The application is structured with the following components:

1. **Main CLI Loop**: Handles user input and command execution
2. **Command System**: Modular command handlers with unified interface
3. **HTTP Client**: Communicates with the PokéAPI
4. **Cache System**: Implements an efficient in-memory cache with automatic expiration
5. **Pokédex Storage**: Keeps track of caught Pokémon

## Technical Details

- **Cache Implementation**: Thread-safe with mutex locks and automatic cleanup
- **Random Encounters**: Uses Go's math/rand/v2 package for catch probability
- **JSON Parsing**: Structured decoding of API responses
- **Error Handling**: Robust error propagation and user-friendly messages

## License

[MIT License](LICENSE)

## Acknowledgements

- [PokéAPI](https://pokeapi.co/) for providing the Pokémon data
- The Go community for excellent libraries and documentation

