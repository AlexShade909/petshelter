package controller

import (
	"net/http"
	"petshelter/internal/models"
	"petshelter/internal/service"
)

func NicknamesHandler(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nicknames := service.Nicknames(dogs)
		WriteJSON(w, http.StatusOK, map[string][]string{
			"nicknames": nicknames,
		})
	}
}
