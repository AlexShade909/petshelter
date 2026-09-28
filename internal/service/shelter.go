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
	//if shelterNumber >= len(sheltersData) {
	//	return models.Shelter{}, errors.New("dog not found")
	//}
	v := sheltersData[shelterNumber]
	return v, nil
}
