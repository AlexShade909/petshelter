package service

import (
	"errors"
	"petshelter/internal/models"
	"petshelter/internal/repository"
	"sort"
	"strconv"
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
	for i, v := range dogsData {
		nicknames = append(nicknames, v.Nickname+" ID: "+strconv.Itoa(i))
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

func (d Dog) Delete(ID int) (error, string) {
	_, ok := dogsData[ID]
	if !ok {
		return errors.New("dog not found"), ""
	}
	nickname := dogsData[ID].Nickname
	delete(dogsData, ID)
	return nil, nickname
}

var nexDogID = len(dogsData)

func (d Dog) Create(dog models.Dog) (models.Dog, error) {
	dog.ID = nexDogID
	nexDogID++
	dogsData[dog.ID] = dog
	return dog, nil
}

func (d Dog) Update(ID int, patch models.Dog) (models.Dog, error) {
	dog, ok := dogsData[ID]
	if !ok {
		return models.Dog{}, errors.New("dog not found")
	}

	if patch.Age != "" {
		dog.Age = patch.Age
	}
	if patch.Nickname != "" {
		dog.Nickname = patch.Nickname
	}
	if patch.CheckInDate != "" {
		dog.CheckInDate = patch.CheckInDate
	}
	if patch.WeightKg != "" {
		dog.WeightKg = patch.WeightKg
	}
	if patch.Сlinic.Address != "" {
		dog.Сlinic = patch.Сlinic
	}

	if patch.Shelter.Address != "" {
		dog.Shelter = patch.Shelter
	}
	dogsData[ID] = dog
	return dogsData[ID], nil
}

func (d Dog) Replace(ID int, patch models.Dog) (models.Dog, error) {
	_, ok := dogsData[ID]
	patch.ID = ID
	if !ok {
		return models.Dog{}, errors.New("dog not found")
	}
	dogsData[ID] = patch
	return dogsData[ID], nil
}
