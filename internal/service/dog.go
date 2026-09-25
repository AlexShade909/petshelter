package service

import (
	"errors"
	"petshelter/internal/models"
	"sort"
)

func DogsListNicknames(dogs map[string]models.Dog) []string {
	nicknames := make([]string, 0, len(dogs))
	for nickname := range dogs {
		nicknames = append(nicknames, nickname)
	}
	sort.Strings(nicknames)
	return nicknames
}

func DogInfo(dogs map[string]models.Dog, dogName string) (models.Dog, error) {
	if dogName == "" {
		return models.Dog{}, errors.New("dog name is required")
	}
	dog, ok := dogs[dogName]
	if !ok {
		return models.Dog{}, errors.New("dog not found")
	}
	return dog, nil
}

func DogDelete(dogs map[string]models.Dog, dogName string) error {
	_, ok := dogs[dogName]
	if !ok {
		return errors.New("dog not found")
	}
	delete(dogs, dogName)
	return nil
}
