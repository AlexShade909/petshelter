package service

import (
	"errors"
	"petshelter/internal/models"
)

type DogRepository interface {
	GetByID(id int) (models.Dog, error)
	ListNicknames() ([]string, error)
	Delete(ID int) error
	Create(dog models.Dog) (models.Dog, error)
	Update(ID int, patch models.Dog) (models.Dog, error)
	Replace(ID int, patch models.Dog) (models.Dog, error)
}

type DogService struct {
	dogRepository DogRepository
}

func NewDogService(dogRepository DogRepository) *DogService {
	return &DogService{dogRepository: dogRepository}
}

func (s *DogService) ListNicknames() ([]string, error) {
	nicknames, err := s.dogRepository.ListNicknames()
	if err != nil {
		return nil, errors.New("dog not found")
	}
	return nicknames, nil
}

func (s *DogService) GetByID(ID int) (models.Dog, error) {
	dog, err := s.dogRepository.GetByID(ID)
	if err != nil {
		return models.Dog{}, errors.New("dog not found")
	}
	return dog, nil
}

func (s *DogService) Delete(ID int) error {
	err := s.dogRepository.Delete(ID)
	if err != nil {
		return errors.New("dog not found")
	}
	return err
}

//var nexDogID = len(dogsData)

//func (d Dog) Create(dog models.Dog) (models.Dog, error) {
//	dog.ID = nexDogID
//	nexDogID++
//	dogsData[dog.ID] = dog
//	return dog, nil
//}

//func (d Dog) Update(ID int, patch models.Dog) (models.Dog, error) {
//	dog, ok := dogsData[ID]
//	if !ok {
//		return models.Dog{}, errors.New("dog not found")
//	}
//
//	if patch.Age != 0 {
//		dog.Age = patch.Age
//	}
//	if patch.Nickname != "" {
//		dog.Nickname = patch.Nickname
//	}
//	if patch.CheckInDate != "" {
//		dog.CheckInDate = patch.CheckInDate
//	}
//	if patch.WeightKg != 0 {
//		dog.WeightKg = patch.WeightKg
//	}
//	if patch.Сlinic.Address != "" {
//		dog.Сlinic = patch.Сlinic
//	}
//
//	if patch.Shelter.Address != "" {
//		dog.Shelter = patch.Shelter
//	}
//	dogsData[ID] = dog
//	return dogsData[ID], nil
//}

//func (d Dog) Replace(ID int, patch models.Dog) (models.Dog, error) {
//	_, ok := dogsData[ID]
//	patch.ID = ID
//	if !ok {
//		return models.Dog{}, errors.New("dog not found")
//	}
//	dogsData[ID] = patch
//	return dogsData[ID], nil
//}
