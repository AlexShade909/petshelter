package controller

import (
	"encoding/json"
	"net/http"
	"petshelter/internal/models"
)

type DogService interface {
	ListNicknames() []string
	Info(nickname string) (models.Dog, error)
	Delete(nickname string) error
	Create(dog models.Dog) error
	Update(dog models.Dog) error
	Replace(dog models.Dog) error
}

type Dog struct {
	dogService DogService
}

func NewDog(dogService DogService) Dog {
	return Dog{dogService: dogService}
}

// TODO: убрать  dogs map[string]models.Dog
// TODO: убрать колбек функцию

func (d *Dog) NicknamesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//nicknames := service.ListNicknames(dogs)
		nicknames := d.dogService.ListNicknames()
		WriteJSON(w, http.StatusOK, map[string][]string{
			"nicknames": nicknames,
		})
	}
}

func (d *Dog) InfoHandler(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dogName := r.PathValue("dogName")
		dogInfo, err := d.dogService.Info(dogName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, http.StatusOK, dogInfo)
		return
	}
}

func (d *Dog) DeleteHandler(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dogName := r.PathValue("dogName")
		err := d.dogService.Delete(dogName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, http.StatusOK, "Dog Deleted. His name: "+dogName)
	}
}

func (d *Dog) CreateHandler(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dog models.Dog
		if err := json.NewDecoder(r.Body).Decode(&dog); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		err := d.dogService.Create(dog)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, http.StatusCreated, dog)
	}
}

func (d *Dog) UpdateHandler(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dogUpdateFields models.Dog
		if err := json.NewDecoder(r.Body).Decode(&dogUpdateFields); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		err := d.dogService.Update(dogUpdateFields)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, http.StatusCreated, dogs[dogUpdateFields.Nickname])
	}
}

func (d *Dog) ReplaceHandler(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dogReplaceFields models.Dog
		if err := json.NewDecoder(r.Body).Decode(&dogReplaceFields); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		err := d.dogService.Replace(dogReplaceFields)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, http.StatusCreated, dogs[dogReplaceFields.Nickname])
	}
}
