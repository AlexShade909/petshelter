package service

import (
	"fmt"
	"petshelter/internal/models"
)

// TODO: Написать сервис собак

//type Dog struct{}

// не dogPrint а dogInfo  GET /dog/{nickname}
func (d *Dog) Print() {
	fmt.Println("Кличка:", d.Nickname)
	fmt.Println("Возраст, лет:", d.Age)
	fmt.Println("Вес, кг:", d.WeightKg)
	fmt.Println("Когда попал в приют:", d.CheckInDate)
	fmt.Println("К какому шелтеру относится:", d.Shelter.Address)
	fmt.Println("К какой поликлинике относится:", d.Policlinic.Address)
}

// GET /dog
func (d *models.Dog) dogList() [string]Dog {

}
