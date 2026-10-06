package repository

import (
	"petshelter/internal/models"
)

func CreateClinics() map[int]models.Clinic {
	Clinics := map[int]models.Clinic{
		0: {
			Address:     "Мира 1",
			PhoneNumber: "+37529-566-13-54",
			WorkingTime: "10:00-23:00",
		},
		1: {
			Address:     "Ленина 133",
			PhoneNumber: "+37529-644-55-71",
			WorkingTime: "09:00-24:00",
		},
	}
	return Clinics
}

func CreateShelters() map[int]models.Shelter {
	Shelters := map[int]models.Shelter{
		0: {
			Address:     "Шелтор 0 , Пятруся Глебки 17",
			PhoneNumber: "+375 29 511-22-13",
			WorkingTime: "10:00 - 22:00",
		},
		1: {
			Address:     "Шелтор 1, Мстислава Чудотворца 4/1",
			PhoneNumber: "+375 12 544-65-45",
			WorkingTime: "11:00 - 21:15",
		},
	}
	return Shelters
}

func CreateDogs(shelterData map[int]models.Shelter, Сlinic map[int]models.Clinic) map[int]models.Dog {
	dogs := map[int]models.Dog{
		0: {
			ID:          0,
			Nickname:    "Чарли",
			Age:         "12",
			WeightKg:    "135",
			CheckInDate: "05.02.2025",
			Shelter:     shelterData[0],
			Сlinic:      Сlinic[0],
		},
		1: {
			ID:          1,
			Nickname:    "Спайси",
			Age:         "13",
			WeightKg:    "15",
			CheckInDate: "15.12.2025",
			Shelter:     shelterData[1],
			Сlinic:      Сlinic[0],
		},
		2: {
			ID:          2,
			Nickname:    "Кайман",
			Age:         "643",
			WeightKg:    "12",
			CheckInDate: "05.07.2025",
			Shelter:     shelterData[0],
			Сlinic:      Сlinic[1],
		},
	}
	return dogs
}
