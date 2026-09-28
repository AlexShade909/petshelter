package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"petshelter/internal/controller"
	"petshelter/internal/service"
	"time"
)

func RunPetshelter(ctx context.Context) error {
	dogService := service.NewDog()
	dogController := controller.NewDog(dogService)
	shelterService := service.NewShelter()
	shelterController := controller.NewShelter(shelterService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dogs", dogController.NicknamesHandler)
	mux.HandleFunc("GET /dogs/{dogName}", dogController.InfoHandler)
	mux.HandleFunc("DELETE /dogs/{dogName}", dogController.DeleteHandler)
	mux.HandleFunc("POST /dogs", dogController.CreateHandler)
	mux.HandleFunc("PATCH /dogs", dogController.UpdateHandler)
	mux.HandleFunc("PUT /dogs", dogController.ReplaceHandler)

	mux.HandleFunc("GET /shelters", shelterController.ListHandler)
	mux.HandleFunc("POST /shelters", shelterController.CreateHandler)
	mux.HandleFunc("GET /shelters/{NumberShelter}", shelterController.InfoHandler)
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
