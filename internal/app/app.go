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

	clinicService := service.NewClinic()
	clinicController := controller.NewClinic(clinicService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dogs", dogController.NicknamesHandler)
	mux.HandleFunc("GET /dogs/{dogID}", dogController.InfoHandler)
	mux.HandleFunc("DELETE /dogs/{dogID}", dogController.DeleteHandler)
	mux.HandleFunc("POST /dogs", dogController.CreateHandler)
	mux.HandleFunc("PATCH /dogs/{dogID}", dogController.UpdateHandler)
	mux.HandleFunc("PUT /dogs/{dogID}", dogController.ReplaceHandler)

	mux.HandleFunc("GET /shelters", shelterController.FullInfoHandler)
	mux.HandleFunc("GET /shelters/{NumberShelter}", shelterController.InfoHandler)
	mux.HandleFunc("POST /shelters", shelterController.CreateHandler)
	mux.HandleFunc("DELETE /shelters/{NumberShelter}", shelterController.DeleteHandler)
	mux.HandleFunc("PATCH /shelters/{NumberShelter}", shelterController.UpdateHandler)
	mux.HandleFunc("PUT /shelters/{NumberShelter}", shelterController.ReplaceHandler)

	mux.HandleFunc("GET /clinics", clinicController.FullInfoHandler)
	mux.HandleFunc("GET /clinics/{NumberClinic}", clinicController.InfoHandler)
	mux.HandleFunc("POST /clinics", clinicController.CreateHandler)
	mux.HandleFunc("DELETE /clinics/{NumberClinic}", clinicController.DeleteHandler)
	mux.HandleFunc("PATCH /clinics/{NumberClinic}", clinicController.UpdateHandler)
	mux.HandleFunc("PUT /clinics/{NumberClinic}", clinicController.ReplaceHandler)

	go func() {
		log.Println("server is running... \nwait request")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Println(err.Error())
		}
	}()
	<-ctx.Done()
	fmt.Println("server stoped. \ngraceful shutdown start...")
	ticker := time.NewTicker(time.Second)
	for i := 1; i < 2; i++ {
		fmt.Println(i)
		<-ticker.C
	}
	fmt.Print("graceful shutdown finish")
	return nil
}
