package service

import (
	"errors"
	"petshelter/internal/models"
)

func SheltersList(Shelters []models.Shelter) ([]string, error) {
	list := make([]string, 0, len(Shelters))
	if len(Shelters) == 0 {
		return nil, errors.New("Empty list Shelters")
	}

	for _, v := range Shelters {
		list = append(list, v.NumberShelter)
	}
	return list, nil
}
