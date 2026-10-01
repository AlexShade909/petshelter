package controller

import (
	"encoding/json"
	"net/http"
	"petshelter/internal/models"
	"strconv"
)

type clinicService interface {
	FullInfo() (map[int]models.Clinic, error)
	Info(clinicNumber int) (models.Clinic, error)
	ListClinics(clinicNumber int) ([]string, error)
	Create(clinic models.Clinic) error
	Delete(clinicNumber int) error
	Update(clinicNumber int, patch models.ClinicPatch) (models.Clinic, error)
	Replace(clinicNumber int, clinic models.Clinic) (models.Clinic, error)
}

type clinic struct {
	clinicService clinicService
}

func NewClinic(clinicService clinicService) clinic {
	return clinic{clinicService: clinicService}
}

func (c clinic) FullInfoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := c.clinicService.FullInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusOK, data)
}

func (c clinic) InfoHandler(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		Info     models.Clinic
		ListDogs []string
	}
	var resp responce
	var err error

	clinicNumber, _ := strconv.Atoi(r.PathValue("NumberClinic"))
	resp.Info, err = c.clinicService.Info(clinicNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.ListDogs, err = c.clinicService.ListClinics(clinicNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, http.StatusOK, resp)
}

func (c clinic) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var clinicCreate models.Clinic

	if err := json.NewDecoder(r.Body).Decode(&clinicCreate); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := c.clinicService.Create(clinicCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusCreated, clinicCreate)
}

func (c clinic) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	clinicNumber, _ := strconv.Atoi(r.PathValue("NumberClinic"))

	err := c.clinicService.Delete(clinicNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := "data in clinic has been deleted. Number deleted clinic: " +
		strconv.Itoa(clinicNumber)

	WriteJSON(w, http.StatusOK, resp)
}

func (c clinic) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	clinicNumber, _ := strconv.Atoi(r.PathValue("NumberClinic"))

	var patch models.ClinicPatch

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&patch); err != nil {
		http.Error(w, "invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.clinicService.Update(clinicNumber, patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusOK, data)
}

func (c clinic) ReplaceHandler(w http.ResponseWriter, r *http.Request) {
	clinicNumber, _ := strconv.Atoi(r.PathValue("NumberClinic"))

	var clinic models.Clinic

	if err := json.NewDecoder(r.Body).Decode(&clinic); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	data, err := c.clinicService.Replace(clinicNumber, clinic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusOK, data)
}
