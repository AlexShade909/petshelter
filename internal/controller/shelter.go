package controller

import (
	"encoding/json"
	"net/http"
	"petshelter/internal/models"
	"strconv"
)

type shelterService interface {
	FullInfo() (map[int]models.Shelter, error)
	Info(shelterNumber int) (models.Shelter, error)
	ListDogs(shelterNumber int) ([]string, error)
	Create(shelter models.Shelter) error
	Delete(shelterNumber int) error
	Update(shelterNumber int, patch models.ShelterPatch) (models.Shelter, error)
	Replace(shelterNumber int, patch models.Shelter) (models.Shelter, error)
}
type shelter struct {
	shelterService shelterService
}

func NewShelter(shelterService shelterService) shelter {
	return shelter{shelterService: shelterService}
}

func (s shelter) FullInfoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := s.shelterService.FullInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, data)
}

func (s shelter) InfoHandler(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		Info     models.Shelter
		ListDogs []string
	}
	var resp responce
	var err error

	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	resp.Info, err = s.shelterService.Info(shelterNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.ListDogs, err = s.shelterService.ListDogs(shelterNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, resp)
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

func (s shelter) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	err := s.shelterService.Delete(shelterNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := "data in shelter has been deleted. PhoneNumber deleted shelter: " + strconv.Itoa(shelterNumber)
	WriteJSON(w, http.StatusOK, resp)
}

func (s shelter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))

	var patch models.ShelterPatch
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&patch); err != nil {
		http.Error(w, "invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}
	data, err := s.shelterService.Update(shelterNumber, patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, data)
}

func (s shelter) ReplaceHandler(w http.ResponseWriter, r *http.Request) {
	shelterNumber, _ := strconv.Atoi(r.PathValue("NumberShelter"))
	var patch models.Shelter
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	data, err := s.shelterService.Replace(shelterNumber, patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, data)
}
