package main

import (
	"petshelter/cli"
	"petshelter/internal"
	"petshelter/network"
	"log"
)

func main() {
	conn, err := network.Connect()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()

	cli.Init(conn, conn)

	Shelters := internal.CreateShelters()
	Policlinics := internal.CreatePoliclinics()
	Dogs := internal.CreateDogs(Shelters, Policlinics)

	flag := true

	for flag {
		choice := cli.ReadMenuChoice("1. Выбрать собаку\n2. Добавить собаку\n3. Выход\n ", 1, 3)
		switch choice {
		case 1:
			cli.Println("Выбрать собаку, я пользователь")
			flag = cli.ScenarioTakeDog(Dogs)
		case 2:
			cli.Println("Добавить собаку, я администратор")
			flag = cli.ScenarioAddDog(Dogs, Shelters, Policlinics)
		case 3:
			cli.Println("Выход")
			return
		}
	}
}
