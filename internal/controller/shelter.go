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
