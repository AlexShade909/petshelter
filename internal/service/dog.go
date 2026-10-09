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

func (s *DogService) Create(dog models.Dog) (models.Dog, error) {
	dogCreated, err := s.dogRepository.Create(dog)
	if err != nil {
		return models.Dog{}, errors.New("not created")
	}
	return dogCreated, nil
}

func (s *DogService) Update(ID int, patch models.Dog) (models.Dog, error) {
	dogUpdate, err := s.dogRepository.Update(ID, patch)
	if err != nil {
		return models.Dog{}, errors.New("not updated")
	}
	return dogUpdate, err
}

//func (d Dog) Replace(ID int, patch models.Dog) (models.Dog, error) {
//	_, ok := dogsData[ID]
//	patch.ID = ID
//	if !ok {
//		return models.Dog{}, errors.New("dog not found")
//	}
//	dogsData[ID] = patch
//	return dogsData[ID], nil
//}
