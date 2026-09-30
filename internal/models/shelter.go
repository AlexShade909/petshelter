package models

type Shelter struct {
	Address     string
	Number      string
	WorkingTime string
}

type ShelterPatch struct {
	Address     string `json:"Address"`
	Number      string `json:"Number"`
	WorkingTime string `json:"WorkingTime"`
}
