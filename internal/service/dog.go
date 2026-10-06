package service

import (
	"errors"
	"petshelter/internal/models"
	"petshelter/internal/repository"
	"sort"
)

type Dog struct {
}

func NewDog() Dog {
	return Dog{}
}

var (
	sheltersData  = repository.CreateShelters()
	nextShelterID = len(sheltersData)
	clinicsData   = repository.CreateClinics()
	nextClinicID  = len(clinicsData)
	dogsData      = repository.CreateDogs(sheltersData, clinicsData)
)

func (d Dog) ListNicknames() []string {
	nicknames := make([]string, 0, len(dogsData))
	for _, v := range dogsData {
		nicknames = append(nicknames, v.Nickname)
	}
	sort.Strings(nicknames)
	return nicknames
}

func (d Dog) Info(ID int) (models.Dog, error) {
	dog, ok := dogsData[ID]
	if !ok {
		return models.Dog{}, errors.New("dog not found")
	}
	return dog, nil
}

func (d Dog) Delete(nickname string) error {
	_, ok := dogsData[nickname]
	if !ok {
		return errors.New("dog not found")
	}
	delete(dogsData, nickname)
	return nil
}

func (d Dog) Create(dog models.Dog) error {
	_, ok := dogsData[dog.Nickname]
	if ok {
		return errors.New("nickname occupied")
	}
	dogsData[dog.Nickname] = dog
	return nil
}

func (d Dog) Update(dogUpdateFields models.Dog) (models.Dog, error) {
	_, ok := dogsData[dogUpdateFields.Nickname]
	if !ok {
		return models.Dog{}, errors.New("dog not found")
	}
	dogsData[dogUpdateFields.Nickname] = dogUpdateFields
	// TODO: this PUT metod fix to PATCH.
	return dogsData[dogUpdateFields.Nickname], nil
}

func (d Dog) Replace(dogReplaceFields models.Dog) (models.Dog, error) {
	_, ok := dogsData[dogReplaceFields.Nickname]
	if !ok {
		return models.Dog{}, errors.New("dog not found")
	}
	dogsData[dogReplaceFields.Nickname] = dogReplaceFields
	return dogsData[dogReplaceFields.Nickname], nil
}
