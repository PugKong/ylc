package pokemon

import (
	"bytes"
	"crypto/rand"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNames(t *testing.T) {
	names := NewNames()

	t.Run("it initializes generator", func(t *testing.T) {
		require.NotNil(t, names.available)
		require.Greater(t, len(names.available), 1_000)

		require.NotNil(t, names.rand)
		require.Equal(t, rand.Reader, names.rand)
	})

	t.Run("it filters pokemons.txt content for empty lines", func(t *testing.T) {
		emptymons := []string{}
		for _, pokemon := range names.available {
			if strings.Trim(pokemon, " ") == "" {
				emptymons = append(emptymons, pokemon)
			}
		}

		require.Empty(t, emptymons, 0)
	})
}

func newNames() (*Names, []string) {
	pokemons := []string{"woobat", "xatu"}

	names := NewNames()
	names.available = pokemons

	return names, pokemons
}

func TestNames_Generate(t *testing.T) {
	t.Run("it returns different pokemons names", func(t *testing.T) {
		names, pokemons := newNames()

		poke1, err := names.Generate()
		require.NoError(t, err)
		require.Contains(t, pokemons, poke1)

		poke2, err := names.Generate()
		require.NoError(t, err)
		require.Contains(t, pokemons, poke2)

		require.NotEqual(t, poke1, poke2)
	})

	t.Run("it returns ErrNoPokemons when no available names left", func(t *testing.T) {
		names, _ := newNames()

		_, err := names.Generate()
		require.NoError(t, err)

		_, err = names.Generate()
		require.NoError(t, err)

		_, err = names.Generate()
		require.ErrorIs(t, err, ErrNoPokemons)
		require.EqualError(t, err, "no pokemons left")
	})

	t.Run("it handles select random pokemon error", func(t *testing.T) {
		names, _ := newNames()
		names.rand = &bytes.Buffer{}

		_, err := names.Generate()
		require.Error(t, err)
		require.ErrorIs(t, err, io.EOF)
		require.EqualError(t, err, "select random pokemon: EOF")
	})
}

func TestNames_Occupy(t *testing.T) {
	t.Run("it removes pokemon from available list", func(t *testing.T) {
		names, _ := newNames()

		names.Occupy("woobat")
		require.Equal(t, []string{"xatu"}, names.available)
	})
}
