package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"petshelter/internal/controller"
	"petshelter/internal/repository"
	"time"
)

func RunPetshelter(ctx context.Context) error {

	// Тут происходит инициализация сервера, контроллеров, сервисов
	// 1) Инициализируем сервисы, которые занимаются бизнес-логикой
	// 2) Инициализируем контроллеры, которые обрабатывают входящие запросы
	// 3) Подводим эти контроллеры к интерфейсу http.Handler
	// 4) Инициализируем и запускаем сервер http
	// 5) Обрабатывает выход из приложения
	Shelters := repository.CreateShelters()
	Policlinics := repository.CreatePoliclinics()
	Dogs := repository.CreateDogs(Shelters, Policlinics)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dogs/", controller.Nicknames(Dogs))
	mux.HandleFunc("GET /dogs/{dogName}", controller.DogInfo(Dogs))
	mux.HandleFunc("DELETE /dogs/{dogName}", controller.DeleteDog(Dogs))
	mux.HandleFunc("POST /dogs/", controller.CreateDog(Dogs))

	go func() {
		log.Println("server is running... \nWait request")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Println(err.Error())
		}
	}()
	<-ctx.Done()
	fmt.Println("Server stoped. \nGraceful shutdown start...")
	ticker := time.NewTicker(time.Second)
	for i := 1; i < 2; i++ {
		fmt.Println(i)
		<-ticker.C
	}
	fmt.Print("Graceful shutdown finish")
	return nil
}
