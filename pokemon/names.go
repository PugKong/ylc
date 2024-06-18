package pokemon

import (
	"crypto/rand"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
)

//go:embed names.txt
var pokemonsTxt string

type Names struct {
	available []string
	rand      io.Reader
}

func NewNames() *Names {
	available := []string{}
	for _, pokemon := range strings.Split(pokemonsTxt, "\n") {
		if strings.Trim(pokemon, " ") != "" {
			available = append(available, pokemon)
		}
	}

	return &Names{
		available: available,
		rand:      rand.Reader,
	}
}

var ErrNoPokemons = errors.New("no pokemons left")

func (g *Names) Generate() (string, error) {
	if len(g.available) == 0 {
		return "", ErrNoPokemons
	}

	r, err := rand.Int(g.rand, big.NewInt(int64(len(g.available))))
	if err != nil {
		return "", fmt.Errorf("select random pokemon: %w", err)
	}

	name := g.available[r.Int64()]
	g.Occupy(name)

	return name, nil
}

func (g *Names) Occupy(name string) {
	available := make([]string, 0, len(g.available))
	for _, availableName := range g.available {
		if name == availableName {
			continue
		}
		available = append(available, availableName)
	}

	g.available = available
}
