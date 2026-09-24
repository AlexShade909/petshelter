package app

import (
	"context"
	"log"
	"net/http"
	"petshelter/internal/controller"
	"petshelter/internal/repository"
)

type App struct {
	server     any // HTTP
	controller any
	service    any
}

func New() *App {
	return &App{}
}

func (a *App) RunPetshelter(ctx context.Context) error {

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
	mux.HandleFunc("GET /dogs/", controller.NicknamesHandler(Dogs))

	go func() {
		log.Println("server is running...")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Println(err.Error())
		}
	}()
	// a.run()
	<-ctx.Done()
	// Тут происходит завершение сервера, контроллеров*, сервисов*
	return nil
}

func (a *App) run() {

	//conn, err := network.Connect()
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}
	//defer conn.Close()
	//client.Init(conn, conn)
	//flag := true
	//for flag {
	//	choice := client.ReadMenuChoice("1. Выбрать собаку\n2. Добавить собаку\n3. Выход\n ", 1, 3)
	//	switch choice {
	//	case 1:
	//		client.Println("Выбрать собаку, я пользователь")
	//		flag = client.ScenarioTakeDog(Dogs)
	//	case 2:
	//		client.Println("Добавить собаку, я администратор")
	//		flag = client.ScenarioAddDog(Dogs, Shelters, Policlinics)
	//	case 3:
	//		client.Println("Выход")
	//		return
	//	}
	//}
}
