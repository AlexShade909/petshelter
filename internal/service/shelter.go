package service

import (
	"errors"
	"petshelter/internal/models"
)

type shelter struct{}

func NewShelter() shelter {
	return shelter{}
}

func (s shelter) FullInfo() (map[int]models.Shelter, error) {
	return sheltersData, nil
}

func (s shelter) Info(shelterNumber int) (models.Shelter, error) {
	//if shelterNumber >= len(sheltersData) && shelterNumber < 0 {
	//	return models.Shelter{}, errors.New("number not correct")
	//}
	v := sheltersData[shelterNumber]

	return v, nil
}

func (s shelter) ListDogs(shelterNumber int) ([]string, error) {
	var listDogs []string
	if shelterNumber >= len(sheltersData) || shelterNumber < 0 {
		return []string{}, errors.New("number not correct")
	}
	for i, v := range dogsData {
		if v.Shelter == sheltersData[shelterNumber] {
			listDogs = append(listDogs, i)
		}
	}
	return listDogs, nil
}

func (s shelter) Create(shelterCreate models.Shelter) error {
	sheltersData[nextShelterID] = shelterCreate
	nextShelterID++
	return nil
}

func (s shelter) Delete(shelterNumber int) error {
	delete(sheltersData, shelterNumber)
	return nil
}

func (s shelter) Update(shelterNumber int, patch models.ShelterPatch) (models.Shelter, error) {
	sh, ok := sheltersData[shelterNumber]
	if !ok {
		return models.Shelter{}, errors.New("shelter not found")
	}

	if patch.Address != "" {
		sh.Address = patch.Address
	}
	if patch.PhoneNumber != "" {
		sh.PhoneNumber = patch.PhoneNumber
	}
	if patch.WorkingTime != "" {
		sh.WorkingTime = patch.WorkingTime
	}

	sheltersData[shelterNumber] = sh
	return sheltersData[shelterNumber], nil
}

func (s shelter) Replace(shelterNumber int, patch models.Shelter) (models.Shelter, error) {
	sheltersData[shelterNumber] = patch
	return sheltersData[shelterNumber], nil
}
