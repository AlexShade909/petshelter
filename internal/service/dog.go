package service

import (
	"errors"
	"fmt"
	"petshelter/internal/models"
	"sort"
)

//type DogService interface {
//	ListNicknames() []string
//	Info(nickname string) (models.Dog, error)
//	Delete(nickname string) error
//	Create(dog models.Dog) error
//	Update(dog models.Dog) error
//	Replace(dog models.Dog) error
//}

type Dog struct {
}

func NewDog() Dog {
	return Dog{}
}

func (d Dog) ListNicknames() []string {
	nicknames := make([]string, 0, len(dogs))
	for nickname := range dogs {
		nicknames = append(nicknames, nickname)
	}
	sort.Strings(nicknames)
	return nicknames
}

func (d Dog) Info(nickname string) (models.Dog, error) {
	if dogName == "" {
		return models.Dog{}, errors.New("dog name is required")
	}
	dog, ok := dogs[dogName]
	if !ok {
		return models.Dog{}, errors.New("dog not found")
	}
	return dog, nil
}

func (d Dog) Delete(nickname string) error {
	_, ok := dogs[dogName]
	if !ok {
		return errors.New("dog not found")
	}
	delete(dogs, dogName)
	return nil
}

func (d Dog) Create(dog models.Dog) error {
	if dog.Nickname == "" {
		fmt.Println(dog.Nickname)
		return errors.New("Empty nickname")
	}
	_, ok := dogs[dog.Nickname]
	if ok {
		return errors.New("Nickname occuped")
	}
	dogs[dog.Nickname] = dog
	return nil
}

func (d Dog) Update(dog models.Dog) error {
	if dogUpdateFields.Nickname == "" {
		return errors.New("Nickname is empty")
	}
	_, ok := dogs[dogUpdateFields.Nickname]
	if !ok {
		return errors.New("dog not found")
	}
	dogs[dogUpdateFields.Nickname] = dogUpdateFields
	// TODO: this PUT metod fix to PATCH.
	return nil
}

func (d Dog) Replace(dog models.Dog) error {
	if dogReplaceFields.Nickname == "" {
		return errors.New("Nickname is empty")
	}
	_, ok := dogs[dogReplaceFields.Nickname]
	if !ok {
		return errors.New("dog not found")
	}
	dogs[dogReplaceFields.Nickname] = dogReplaceFields
	return nil
}
