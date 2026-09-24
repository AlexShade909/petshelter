package models

type Dog struct {
	Nickname    string
	Age         string
	WeightKg    string
	CheckInDate string
	Shelter     *Shelter
	Policlinic  *Policlinic
}
