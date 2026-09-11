package cli

import (
	"petshelter/internal"
)

func PrintDogInfo(dogs map[string]internal.Dog, nickname string) {
	d := dogs[nickname]
	Println("Кличка: ", d.Nickname)
	Print("Возраст, лет: ")
	Println(d.Age)
	Print("Вес, кг: ")
	Println(d.WeightKg)
	Print("Когда попал в приют: ")
	Println(d.CheckInDate)
	Print("К какому шелтеру относится: ")
	Println(d.Shelter.Address)
	Print("К какой поликлинике относится: ")
	Println(d.Policlinic.Address)
}

func PrintDogList(dogs map[string]internal.Dog) {
	Println("Список собак: ")
	for nickname := range dogs {
		Println(nickname)
	}
}

func PrintShelterInfo(s *internal.Shelter) {
	Print("Приют №: ")
	Println(s.NumberShelter)
	Print("Адрес: ")
	Println(s.Address)
	Print("Телефон: ")
	Println(s.Number)
	Print("Время работы: ")
	Println(s.WorkingTime)
}

func PrintPoliclinicInfo(p *internal.Policlinic) {
	Print("Название: ")
	Println(p.NumberClinic)
	Print("адрес: ")
	Println(p.Address)
	Print("телефон: ")
	Println(p.PhoneNumber)
	Print("рабочее время: ")
	Println(p.WorkingTime)
}
