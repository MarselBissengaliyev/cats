package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/MarselBissengaliyev/cats/internal/models"
	"github.com/MarselBissengaliyev/cats/internal/utils"
)

const catBreedsURL = "https://catfact.ninja/breeds"

// ConvertToCatBreeds преобразует данные в структуру CatBreed
func FetchCatBreeds(url string) ([]models.CatBreed, error) {
	// Получение данных из API
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Чтение данных из ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Конвертируем данные в мапу
	var breedsData map[string]interface{}
	err = json.Unmarshal(body, &breedsData)
	if err != nil {
		return nil, err
	}

	// Берем данные по ключу data
	breeds, ok := breedsData["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("ошибка при получении данных о породах кошек")
	}

	// Создание канала для данных о породах кошек и канала ошибок
	catBreedsChan := make(chan models.CatBreed)
	errChan := make(chan error)

	// Создание WaitGroup для дожидания завершения всех горутин
	var wg sync.WaitGroup

	// Запуск горутины для каждой породы кошек
	for _, breed := range breeds {
		wg.Add(1)
		go func(breedData interface{}) {
			defer wg.Done()

			breedMap, ok := breedData.(map[string]interface{})
			if !ok {
				errChan <- fmt.Errorf("ошибка при преобразовании данных о породе кошки")
				return
			}

			catBreed := models.CatBreed{
				Breed:   utils.GetStringValue(breedMap, "breed"),
				Country: utils.GetStringValue(breedMap, "country"),
				Origin:  utils.GetStringValue(breedMap, "origin"),
				Coat:    utils.GetStringValue(breedMap, "coat"),
				Pattern: utils.GetStringValue(breedMap, "pattern"),
			}

			catBreedsChan <- catBreed
		}(breed)
	}

	go func() {
		wg.Wait()
		close(catBreedsChan)
		close(errChan)
	}()

	// Сбор данных о породах кошек и ошибок
	var catBreeds []models.CatBreed

	for {
		select {
		case breed, ok := <-catBreedsChan:
			if !ok {
				// Канал закрыт, все горутины завершились
				return catBreeds, nil
			}
			catBreeds = append(catBreeds, breed)
		case err := <-errChan:
			// Обработка ошибок
			return nil, err
		}
	}
}
