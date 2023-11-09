package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarselBissengaliyev/cats/internal/api"
	"github.com/MarselBissengaliyev/cats/internal/models"
)

var (
	url = "https://catfact.ninja/breeds"
)

func main() {
	// Создание канала для принятия сигналов завершения
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Создание контекста для отслеживания сигналов завершения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск главного приложения в горутине
	go func() {
		// Получение данных о породах кошек с API
		catBreeds, err := api.FetchCatBreeds(url)
		if err != nil {
			fmt.Println("Ошибка при получении данных:", err)
			return
		}

		// Обработка и сохранение данных
		models.ProcessAndSaveData(catBreeds)
	}()

	// Ожидание сигналов завершения или завершение по истечении времени
	select {
	case <-stopChan:
		fmt.Println("Принят сигнал завершения, ожидание завершения задач...")
		cancel() // Отмена контекста при получении сигнала
	case <-time.After(5 * time.Second):
		fmt.Println("Прошло 5 секунд, завершение приложения...")
		cancel() // Отмена контекста по таймауту
	}

	// Дождитесь завершения всех задач
	<-ctx.Done()

	fmt.Println("Graceful Shutdown завершен.")
}
