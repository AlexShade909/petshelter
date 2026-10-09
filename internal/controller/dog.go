package controller

import (
	"encoding/json"
	"net/http"
	"petshelter/internal/models"
	"strconv"
)

type DogService interface {
	ListNicknames() ([]string, error)
	GetByID(ID int) (models.Dog, error)
	Delete(ID int) error
	Create(dog models.Dog) (models.Dog, error)
	Update(ID int, patch models.Dog) (models.Dog, error)
	//Replace(ID int, dogReplaceFields models.Dog) (models.Dog, error)
}

type dog struct {
	dogService DogService
}

func NewDog(dogService DogService) dog {
	return dog{dogService: dogService}
}

func (d dog) NicknamesHandler(w http.ResponseWriter, r *http.Request) {
	data, err := d.dogService.ListNicknames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, data)
}

func (d dog) InfoHandler(w http.ResponseWriter, r *http.Request) {
	dogID, err := strconv.Atoi(r.PathValue("dogID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dogInfo, err := d.dogService.GetByID(dogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, dogInfo)
}

func (d dog) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	dogID, err := strconv.Atoi(r.PathValue("dogID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = d.dogService.Delete(dogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, "success")
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
	dogCreated, err := d.dogService.Create(dogCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, dogCreated)
}

func (d dog) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var patch models.Dog

	dogID, err := strconv.Atoi(r.PathValue("dogID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if patch.Nickname == "" {
		http.Error(w, "nickname is empty", http.StatusBadRequest)
		return
	}
	dogUpdated, err := d.dogService.Update(dogID, patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, dogUpdated)
}

//func (d dog) ReplaceHandler(w http.ResponseWriter, r *http.Request) {
//	var dogReplaceFields models.Dog
//	dogID, err := strconv.Atoi(r.PathValue("dogID"))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	if err := json.NewDecoder(r.Body).Decode(&dogReplaceFields); err != nil {
//		http.Error(w, "invalid request body", http.StatusBadRequest)
//		return
//	}
//	if dogReplaceFields.Nickname == "" {
//		http.Error(w, "nickname is empty", http.StatusBadRequest)
//		return
//	}
//	dogResp, err := d.dogService.Replace(dogID, dogReplaceFields)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	WriteJSON(w, http.StatusCreated, dogResp)
//}
