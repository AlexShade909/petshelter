package service

import (
	"errors"
	"petshelter/internal/models"
)

type shelter struct{}

func NewShelter() shelter {
	return shelter{}
}

func (s shelter) SheltersList() ([]int, error) {
	if len(sheltersData) == 0 {
		return nil, errors.New("empty list Shelters")
	}
	list := make([]int, 0, len(sheltersData))
	for _, v := range sheltersData {
		list = append(list, v.NumberShelter)
	}
	return list, nil
}

func (s shelter) Create(shelterCreate models.Shelter) error {
	sheltersData = append(sheltersData, shelterCreate)
	return nil
}

func (s shelter) Info(shelterNumber int) (models.Shelter, error) {
	if shelterNumber >= len(sheltersData) || shelterNumber < 0 {
		return models.Shelter{}, errors.New("number not correct")
	}
	v := sheltersData[shelterNumber]
	return v, nil
}

func (s shelter) ListDogs(shelterNumber int) ([]string, error) {
	var listDogs []string
	if shelterNumber >= len(sheltersData) || shelterNumber < 0 {
		return []string{}, errors.New("number not correct")
	}
	for i, v := range dogsData {
		if shelterNumber == v.Shelter.NumberShelter {
			listDogs = append(listDogs, i)
		}
	}
	return listDogs, nil
}

func (s shelter) Delete(shelterNumber int) error {
	sheltersData[shelterNumber].Number = ""
	sheltersData[shelterNumber].NumberShelter = 0
	sheltersData[shelterNumber].Address = ""
	sheltersData[shelterNumber].WorkingTime = ""
	return nil
}

//TODO: переписать шелтер и клиник на мапу

func (s shelter) Update(shelterNumber int, shelterUpdateData models.Shelter) (models.Shelter, error) {
	sheltersData[shelterNumber] = shelterUpdateData
	return sheltersData[shelterNumber], nil
}

func (s shelter) Replace(shelterNumber int, shelterUpdateData models.Shelter) (models.Shelter, error) {
	sheltersData[shelterNumber] = shelterUpdateData
	return sheltersData[shelterNumber], nil
}
