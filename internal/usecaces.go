package internal

import "petshelter/cli"

func ScenarioTakeDog(dogs map[string]Dog) bool {
	cli.PrintDogList(dogs)
	cli.Println("Введите кличку собаки: ")
	nickname := cli.InputValidString()
	d, ok := FindDog(dogs, nickname)
	if !ok {
		cli.Println("Собака с такой кличкой не найдена")
		return true
	}
	cli.PrintDogInfo(dogs, nickname)
	cli.Println("Забрать из приюта? (да/нет): ")
	if cli.ValidYesNo() {
		RemoveDog(dogs, nickname)
		cli.Println("Собака удалена из общего списка, приюта и поликлиники")
		cli.PrintShelterInfo(d.Shelter)
		cli.PrintPoliclinicInfo(d.Policlinic)
	}
	cli.Println("Смотреть ещё? (да/нет): ")
	return cli.ValidYesNo()
}

func ScenarioAddDog(dogs map[string]Dog, shelters []Shelter, policlinics []Policlinic) bool {
	cli.Println("Введите кличку: ")
	nickname := cli.InputValidString()
	cli.Println("Введите возраст: ")
	age := cli.InputValidString()
	cli.Println("Введите вес: ")
	weight := cli.InputValidString()
	cli.Println("Введите дату поступления: ")
	date := cli.InputValidString()
	cli.Println("Выберите приют: ")
	shelterChoice := cli.ValidationReadMenuChoice(0, len(shelters)-1)
	shelter := &shelters[shelterChoice]
	cli.Println("Выберите поликлинику: ")
	clinicChoice := cli.ValidationReadMenuChoice(0, len(policlinics)-1)
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
	cli.Println("Собака добавлена:")
	cli.Println(dog.Nickname)
	cli.Println("Смотреть ещё? (да/нет): ")
	return cli.ValidYesNo()
}
