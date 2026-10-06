package controller

import (
	"encoding/json"
	"net/http"
	"petshelter/internal/models"
)

type DogService interface {
	ListNicknames() []string
	Info(ID int) (models.Dog, error)
	Delete(nickname string) error
	Create(dog models.Dog) error
	Update(dogUpdateFields models.Dog) (models.Dog, error)
	Replace(dogReplaceFields models.Dog) (models.Dog, error)
}

type dog struct {
	dogService DogService
}

func NewDog(dogService DogService) dog {
	return dog{dogService: dogService}
}

func (d dog) NicknamesHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, d.dogService.ListNicknames())
}

func (d dog) InfoHandler(w http.ResponseWriter, r *http.Request) {
	dogID := r.PathValue("dogID")
	dogInfo, err := d.dogService.Info(dogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, dogInfo)
}

func (d dog) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	dogName := r.PathValue("dogName")
	err := d.dogService.Delete(dogName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, "dog Deleted. His name: "+dogName)
}

func (d dog) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var dogCreate models.Dog
	if err := json.NewDecoder(r.Body).Decode(&dogCreate); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if dogCreate.Nickname == "" {
		http.Error(w, "empty nickname", http.StatusBadRequest)
		return
	}
	err := d.dogService.Create(dogCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, dogCreate)
}

func (d dog) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var dogUpdateFields models.Dog
	if err := json.NewDecoder(r.Body).Decode(&dogUpdateFields); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if dogUpdateFields.Nickname == "" {
		http.Error(w, "nickname is empty", http.StatusBadRequest)
		return
	}
	dogUpdated, err := d.dogService.Update(dogUpdateFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, dogUpdated)
}

func (d dog) ReplaceHandler(w http.ResponseWriter, r *http.Request) {
	var dogReplaceFields models.Dog
	if err := json.NewDecoder(r.Body).Decode(&dogReplaceFields); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if dogReplaceFields.Nickname == "" {
		http.Error(w, "nickname is empty", http.StatusBadRequest)
		return
	}
	dogResp, err := d.dogService.Replace(dogReplaceFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusCreated, dogResp)
}
