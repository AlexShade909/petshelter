package models

type Dog struct {
	ID          int
	Nickname    string
	Age         int
	WeightKg    int
	CheckInDate string
	Shelter     Shelter
	Clinic      Clinic
}
