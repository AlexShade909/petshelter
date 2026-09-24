package service

import "petshelter/internal"

func AddDog(dogs map[string]internal.Dog, nickname string, age string, weightKg string, checkInDate string, shelter *internal.Shelter, policlinic *internal.Policlinic) *internal.Dog {
	d := &internal.Dog{
		Nickname:    nickname,
		Age:         age,
		WeightKg:    weightKg,
		CheckInDate: checkInDate,
		Shelter:     shelter,
		Policlinic:  policlinic,
	}
	dogs[nickname] = *d
	return d
}

func RemoveDog(dogs map[string]internal.Dog, nickname string) bool {
	_, ok := dogs[nickname]
	if !ok {
		return false
	}
	delete(dogs, nickname)
	return true
}

func FindDog(dogs map[string]internal.Dog, nickname string) (internal.Dog, bool) {
	d, ok := dogs[nickname]
	return d, ok
}
