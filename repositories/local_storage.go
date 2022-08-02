package repositories

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"GoConcurrency-Bootcamp-2022/models"
)

type LocalStorage struct{}

const filePath = "resources/pokemons.csv"

func (l LocalStorage) Write(pokemons []models.Pokemon) error {
	file, fErr := os.Create(filePath)
	defer file.Close()
	if fErr != nil {
		return fErr
	}

	w := csv.NewWriter(file)
	records := buildRecords(pokemons)
	if err := w.WriteAll(records); err != nil {
		return err
	}

	return nil
}

func (l LocalStorage) Read() (<-chan models.Pokemon, int, error) {
	file, fErr := os.Open(filePath)
	defer file.Close()
	if fErr != nil {
		return nil, 0, fErr
	}

	r := csv.NewReader(file)
	records, rErr := r.ReadAll()
	if rErr != nil {
		return nil, 0, rErr
	}

	pokemons, size, err := parseCSVData(records)

	if err != nil {
		return nil, 0, err
	}

	return pokemons, size, nil
}

func buildRecords(pokemons []models.Pokemon) [][]string {
	headers := []string{"id", "name", "height", "weight", "flat_abilities"}
	records := [][]string{headers}
	for _, p := range pokemons {
		record := fmt.Sprintf("%d,%s,%d,%d,%s",
			p.ID,
			p.Name,
			p.Height,
			p.Weight,
			p.FlatAbilityURLs)
		records = append(records, strings.Split(record, ","))
	}

	return records
}

func parseCSVData(records [][]string) (<-chan models.Pokemon, int, error) {
	csvRecords := make(chan models.Pokemon)
	wg := sync.WaitGroup{}
	var size = 0

	for i, record := range records {
		if i == 0 {
			continue
		}
		size++
		wg.Add(1)
		go func(record []string) {
			defer wg.Done()

			id, _ := strconv.Atoi(record[0])
			height, _ := strconv.Atoi(record[2])

			weight, _ := strconv.Atoi(record[3])

			pokemon := models.Pokemon{
				ID:              id,
				Name:            record[1],
				Height:          height,
				Weight:          weight,
				Abilities:       nil,
				FlatAbilityURLs: record[4],
				EffectEntries:   nil,
			}

			csvRecords <- pokemon
		}(record)
	}

	go func() {
		wg.Wait()
		close(csvRecords)
	}()

	return csvRecords, size, nil
}
