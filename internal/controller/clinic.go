package controller

//
//import (
//	"encoding/json"
//	"log"
//	"net/http"
//
//	"petshelter/internal/models"
//)
//
//// TODO: Написать контроллер поликлинники
//
//type clinicService interface {
//	GetAllClinics() []models.Policlinic
//}
//
//type Clinic struct {
//	clinicService clinicService
//}
//
//func NewClinic(clinicService clinicService) *Clinic {
//	return &Clinic{clinicService: clinicService}
//}
//
//func (c *Clinic) GetAll(w http.ResponseWriter, r *http.Request) {
//	clinics := c.clinicService.GetAllClinics()
//
//	rawJSON, err := json.Marshal(clinics)
//	if err != nil {
//		log.Panicln(err.Error())
//		return
//	}
//
//	if _, err := w.Write(rawJSON); err != nil {
//		log.Panicln(err.Error())
//		return
//	}
//}
