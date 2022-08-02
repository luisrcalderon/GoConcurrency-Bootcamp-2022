package use_cases

import (
	"context"
	"strings"

	"GoConcurrency-Bootcamp-2022/models"
)

type reader interface {
	Read() (<-chan models.Pokemon, int, error)
}

type saver interface {
	Save(context.Context, []models.Pokemon) error
}

type fetcher interface {
	FetchAbility(string) (models.Ability, error)
}

type Refresher struct {
	reader
	saver
	fetcher
}

func NewRefresher(reader reader, saver saver, fetcher fetcher) Refresher {
	return Refresher{reader, saver, fetcher}
}

func (r Refresher) Refresh(ctx context.Context) error {
	var pokeAbilities []models.Pokemon

	pokemons, size, _ := r.Read()

	abilities1 := getAbilities(pokemons, r)
	abilities2 := getAbilities(pokemons, r)
	abilities3 := getAbilities(pokemons, r)

	for i := 0; i < size; i++ {
		select {
		case value, ok := <-abilities1:
			if ok {
				pokeAbilities = append(pokeAbilities, value)
			}
		case value, ok := <-abilities2:
			if ok {
				pokeAbilities = append(pokeAbilities, value)
			}
		case value, ok := <-abilities3:
			if ok {
				pokeAbilities = append(pokeAbilities, value)
			}
		}
	}

	if err := r.Save(ctx, pokeAbilities); err != nil {
		return err
	}

	return nil
}

func getAbilities(pokemons <-chan models.Pokemon, r Refresher) <-chan models.Pokemon {
	abilitiesChan := make(chan models.Pokemon)

	go func() {
		defer close(abilitiesChan)

		for p := range pokemons {
			urls := strings.Split(p.FlatAbilityURLs, "|")
			var abilities []string
			for _, url := range urls {
				ability, _ := r.FetchAbility(url)

				for _, ee := range ability.EffectEntries {
					abilities = append(abilities, ee.Effect)
				}
			}

			p.EffectEntries = abilities
			abilitiesChan <- p

		}

	}()

	return abilitiesChan

}
