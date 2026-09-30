package service

import (
	"errors"
	"petshelter/internal/models"
)

type shelter struct{}

func NewShelter() shelter {
	return shelter{}
}

func (s shelter) FullInfoHandler() ([]models.Shelter, error) {
	return sheltersData, nil
}

func (s shelter) Info(shelterNumber int) (models.Shelter, error) {
	if shelterNumber >= len(sheltersData) && shelterNumber < 0 {
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

func (s shelter) Create(shelterCreate models.Shelter) error {
	sheltersData = append(sheltersData, shelterCreate)
	return nil
}

func (s shelter) Delete(shelterNumber int) error {
	delete(sheltersData, shelterNumber)
	return nil
}

func (s shelter) Update(shelterNumber int, patch models.ShelterPatch) (models.Shelter, error) {
	if patch.Number != nil {
		sheltersData[shelterNumber].Number = *patch.Number
	}
	if patch.Address != nil {
		sheltersData[shelterNumber].Address = *patch.Address
	}
	if patch.WorkingTime != nil {
		sheltersData[shelterNumber].WorkingTime = *patch.WorkingTime
	}
	return sheltersData[shelterNumber], nil
}

func (s shelter) Replace(shelterNumber int, shelterUpdateData models.Shelter) (models.Shelter, error) {
	sheltersData[shelterNumber] = shelterUpdateData
	sheltersData[shelterNumber].NumberShelter = shelterNumber
	return sheltersData[shelterNumber], nil
}
