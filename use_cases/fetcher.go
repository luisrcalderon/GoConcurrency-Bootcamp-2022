package use_cases

import (
	"strings"
	"sync"

	"GoConcurrency-Bootcamp-2022/models"
)

type api interface {
	FetchPokemon(id int) (models.Pokemon, error)
}

type writer interface {
	Write(pokemons []models.Pokemon) error
}

type Fetcher struct {
	api     api
	storage writer
}

func NewFetcher(api api, storage writer) Fetcher {
	return Fetcher{api, storage}
}

func (f Fetcher) Fetch(from, to int) error {
	var pokemons []models.Pokemon
	var pokemonChannel = generatePokemon(from, to, f)

	for pokemon := range pokemonChannel {
		var flatAbilities []string
		for _, t := range pokemon.Abilities {
			flatAbilities = append(flatAbilities, t.Ability.URL)
		}
		pokemon.FlatAbilityURLs = strings.Join(flatAbilities, "|")

		pokemons = append(pokemons, pokemon)
	}

	return f.storage.Write(pokemons)
}

//Method to generate IDs channel
func generatePokemon(from, to int, f Fetcher) <-chan models.Pokemon {
	var pokemonChannel = make(chan models.Pokemon)
	wg := sync.WaitGroup{}

	for id := from; id <= to; id++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			pokemon, err := f.api.FetchPokemon(id)
			if err != nil {
				close(pokemonChannel)
			}

			pokemonChannel <- pokemon

		}(id)
	}

	go func() {
		wg.Wait()
		close(pokemonChannel)
	}()

	return pokemonChannel
}
