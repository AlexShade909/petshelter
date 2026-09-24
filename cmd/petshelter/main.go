package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"petshelter/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app := app.New()

	if err := app.RunPetshelter(ctx); err != nil {
		log.Fatalln(err.Error())
	}
}
