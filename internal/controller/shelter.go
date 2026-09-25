package controller

import (
	"net/http"
	"petshelter/internal/models"
	"petshelter/internal/service"
)

func SheltersListHandler(Shelters []models.Shelter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := service.SheltersList(Shelters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		WriteJSON(w, http.StatusOK, list)
	}
}
