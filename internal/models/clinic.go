package models

type Clinic struct {
	ID          int
	Address     string
	PhoneNumber string
	WorkingTime string
}
type ClinicPatch struct {
	Address     string `json:"Address"`
	PhoneNumber string `json:"PhoneNumber"`
	WorkingTime string `json:"WorkingTime"`
}
