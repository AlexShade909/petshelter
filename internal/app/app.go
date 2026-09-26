package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"petshelter/internal/controller"
	"petshelter/internal/repository"
	"petshelter/internal/service"
	"time"
)

func RunPetshelter(ctx context.Context) error {
	dogService := service.NewDog()

	dogController := controller.NewDog(dogService)
	// Тут происходит инициализация сервера, контроллеров, сервисов
	// 1) Инициализируем сервисы, которые занимаются бизнес-логикой
	// 2) Инициализируем контроллеры, которые обрабатывают входящие запросы
	// 3) Подводим эти контроллеры к интерфейсу http.Handler
	// 4) Инициализируем и запускаем сервер http
	// 5) Обрабатывает выход из приложения
	shelters := repository.CreateShelters()
	policlinics := repository.CreatePoliclinics()
	dogs := repository.CreateDogs(shelters, policlinics)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dogs/", dogController.NicknamesHandler())
	mux.HandleFunc("GET /dogs/{dogName}", dogController.InfoHandler(dogs))
	mux.HandleFunc("DELETE /dogs/{dogName}", dogController.DeleteHandler(dogs))
	mux.HandleFunc("POST /dogs/", dogController.CreateHandler(dogs))
	mux.HandleFunc("PATCH /dogs/", dogController.UpdateHandler(dogs))
	mux.HandleFunc("PUT /dogs/", dogController.ReplaceHandler(dogs))

	mux.HandleFunc("GET /shelters/", controller.SheltersListHandler(Shelters))
	//mux.HandleFunc(("POST /shelters/", controller.ShelterCreate(Shelters)))
	//mux.HandleFunc(("GET /shelters/{idShelter}", controller.ShelterInfo(Shelters)))
	//mux.HandleFunc(("GET /shelters/{idShelter}/dogs", controller.ShelterListDogs(Shelters)))
	//mux.HandleFunc(("DELETE /shelters/{idShelter}", controller.ShelterDelete(Shelters)))
	//mux.HandleFunc(("PATCH /shelters/{idShelter}", controller.ShelterUpdate(Shelters)))
	//mux.HandleFunc(("PUT /shelters/{idShelter}", controller.ShelterReplace(Shelters)))

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
