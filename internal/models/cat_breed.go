package models

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

// CatBreed структура для хранения информации о породе кошек
type CatBreed struct {
	Breed   string `json:"breed"`
	Country string `json:"country"`
	Origin  string `json:"origin"`
	Coat    string `json:"coat"`
	Pattern string `json:"pattern"`
}

// ProcessAndSaveData обрабатывает и сохраняет данные
func ProcessAndSaveData(catBreeds []CatBreed) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		// Группировка пород по стране происхождения
		groupedBreeds := make(map[string][]string)
		for _, breed := range catBreeds {
			groupedBreeds[breed.Country] = append(groupedBreeds[breed.Country], breed.Breed)
		}

		// Сортировка названий пород по длине
		for _, breeds := range groupedBreeds {
			sort.Slice(breeds, func(i, j int) bool {
				return len(breeds[i]) < len(breeds[j])
			})
		}

		// Создание JSON-структуры
		result := struct {
			GroupedBreeds map[string][]string `json:"groupedBreeds"`
		}{
			GroupedBreeds: groupedBreeds,
		}

		// Запись данных в JSON-файл
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println("Ошибка при маршалинге данных в JSON:", err)
			return
		}

		err = os.WriteFile("out.json", jsonData, 0644)
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}

		fmt.Println("Данные успешно записаны в out.json")
	}()

	wg.Wait()
}
