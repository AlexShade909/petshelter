package internal

import (
	cli2 "petshelter/internal/cli"
)

func ScenarioTakeDog(dogs map[string]Dog) bool {
	cli2.PrintDogList(dogs)
	cli2.Println("Введите кличку собаки: ")
	nickname := cli2.InputValidString()
	d, ok := FindDog(dogs, nickname)
	if !ok {
		cli2.Println("Собака с такой кличкой не найдена")
		return true
	}
	cli2.PrintDogInfo(dogs, nickname)
	cli2.Println("Забрать из приюта? (да/нет): ")
	if cli2.ValidYesNo() {
		RemoveDog(dogs, nickname)
		cli2.Println("Собака удалена из общего списка, приюта и поликлиники")
		cli2.PrintShelterInfo(d.Shelter)
		cli2.PrintPoliclinicInfo(d.Policlinic)
	}
	cli2.Println("Смотреть ещё? (да/нет): ")
	return cli2.ValidYesNo()
}

func ScenarioAddDog(dogs map[string]Dog, shelters []Shelter, policlinics []Policlinic) bool {
	cli2.Println("Введите кличку: ")
	nickname := cli2.InputValidString()
	cli2.Println("Введите возраст: ")
	age := cli2.InputValidString()
	cli2.Println("Введите вес: ")
	weight := cli2.InputValidString()
	cli2.Println("Введите дату поступления: ")
	date := cli2.InputValidString()
	cli2.Println("Выберите приют: ")
	shelterChoice := cli2.ValidationReadMenuChoice(0, len(shelters)-1)
	shelter := &shelters[shelterChoice]
	cli2.Println("Выберите поликлинику: ")
	clinicChoice := cli2.ValidationReadMenuChoice(0, len(policlinics)-1)
	policlinic := &policlinics[clinicChoice]
	dog := AddDog(
		dogs,
		nickname,
		age,
		weight,
		date,
		shelter,
		policlinic,
	)
	cli2.Println("Собака добавлена:")
	cli2.Println(dog.Nickname)
	cli2.Println("Смотреть ещё? (да/нет): ")
	return cli2.ValidYesNo()
}
