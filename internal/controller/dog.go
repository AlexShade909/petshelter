package controller

import (
	"net/http"
	"petshelter/internal/models"
	"petshelter/internal/service"
)

func Nicknames(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nicknames := service.DogsListNicknames(dogs)
		WriteJSON(w, http.StatusOK, map[string][]string{
			"nicknames": nicknames,
		})
	}
}

func DogInfo(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dogName := r.PathValue("dogName")
		dogInfo, err := service.DogInfo(dogs, dogName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, http.StatusOK, dogInfo)
		return
	}
}

func DeleteDog(dogs map[string]models.Dog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dogName := r.PathValue("dogName")
		err := service.DogDelete(dogs, dogName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, http.StatusOK, "Dog Deleted. His name: "+dogName)
	}
}
