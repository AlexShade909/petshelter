package controller

import (
	"encoding/json"
	"net/http"
	"petshelter/internal/models"
	"strconv"
)

type shelterService interface {
	SheltersList() ([]int, error)
	Create(shelter models.Shelter) error
	Info(shelterNumber int) (models.Shelter, error)
	ListDogs(shelterNumber int) ([]string, error)
	Delete(shelterNumber int) error
	Update(shelterNumber int, shelterUpdateData models.Shelter) (models.Shelter, error)
	Replace(shelterNumber int, shelterUpdateData models.Shelter) (models.Shelter, error)
}
type shelter struct {
	shelterService shelterService
}

func NewShelter(shelterService shelterService) shelter {
	return shelter{shelterService: shelterService}
}

func (s shelter) ListHandler(w http.ResponseWriter, r *http.Request) {
	list, err := s.shelterService.SheltersList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (s shelter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var shelterCreate models.Shelter
	if err := json.NewDecoder(r.Body).Decode(&shelterCreate); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := s.shelterService.Create(shelterCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusCreated, shelterCreate)
}
func (s shelter) InfoHandler(w http.ResponseWriter, r *http.Request) {
	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	shelterInfo, err := s.shelterService.Info(shelterNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, shelterInfo)
	return
}

func (s shelter) ListDogsHandler(w http.ResponseWriter, r *http.Request) {
	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	listDogsNicknames, err := s.shelterService.ListDogs(shelterNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, listDogsNicknames)
	return
}

func (s shelter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	err := s.shelterService.Delete(shelterNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, "shelter has been deleted")
}

func (s shelter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	var shelterUpdateData models.Shelter
	if err := json.NewDecoder(r.Body).Decode(&shelterUpdateData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	data, err := s.shelterService.Update(shelterNumber, shelterUpdateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, data)
}

func (s shelter) ReplaceHandler(w http.ResponseWriter, r *http.Request) {
	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	var shelterUpdateData models.Shelter
	if err := json.NewDecoder(r.Body).Decode(&shelterUpdateData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	data, err := s.shelterService.Replace(shelterNumber, shelterUpdateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, data)
}
