package models

type Shelter struct {
	Address     string
	PhoneNumber string
	WorkingTime string
}

type ShelterPatch struct {
	Address     string `json:"Address"`
	PhoneNumber string `json:"PhoneNumber"`
	WorkingTime string `json:"WorkingTime"`
}
