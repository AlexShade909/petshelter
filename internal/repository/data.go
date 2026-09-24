package repository

import "petshelter/internal/models"

func CreatePoliclinics() []models.Policlinic {
	policlinics := []models.Policlinic{
		{
			NumberClinic: "Поликлиника 0",
			Address:      "Мира 1",
			PhoneNumber:  "+37529-566-13-54",
			WorkingTime:  "10:00-23:00",
		},
		{
			NumberClinic: "Поликлиника 1",
			Address:      "Ленина 133",
			PhoneNumber:  "+37529-644-55-71",
			WorkingTime:  "09:00-24:00",
		},
	}
	return policlinics
}

func CreateShelters() []models.Shelter {
	Shelters := []models.Shelter{
		{
			NumberShelter: "Шелтер 0",
			Address:       "Пятруся Глебки 17",
			Number:        "+375 29 511-22-13",
			WorkingTime:   "10:00 - 22:00",
		},
		{
			NumberShelter: "Шелтер 1",
			Address:       "Мстислава Чудотворца 4/1",
			Number:        "+375 12 544-65-45",
			WorkingTime:   "11:00 - 21:15",
		},
	}
	return Shelters
}

func CreateDogs(shelter []models.Shelter, policlinic []models.Policlinic) map[string]models.Dog {
	dogs := map[string]models.Dog{
		"Чарли": {
			Nickname:    "Чарли",
			Age:         "12",
			WeightKg:    "135",
			CheckInDate: "05.02.2025",
			Shelter:     &shelter[0],
			Policlinic:  &policlinic[0],
		},
		"Спайси": {
			Nickname:    "Спайси",
			Age:         "13",
			WeightKg:    "15",
			CheckInDate: "15.12.2025",
			Shelter:     &shelter[1],
			Policlinic:  &policlinic[0],
		},
		"Кайман": {
			Nickname:    "Кайман",
			Age:         "643",
			WeightKg:    "12",
			CheckInDate: "05.07.2025",
			Shelter:     &shelter[0],
			Policlinic:  &policlinic[1],
		},
	}
	return dogs
}
