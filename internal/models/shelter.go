package models

type Shelter struct {
	NumberShelter int
	Address       string
	Number        string
	WorkingTime   string
}

type ShelterPatch struct {
	Address     *string `json:"Address"`
	Number      *string `json:"Number"`
	WorkingTime *string `json:"WorkingTime"`
}
