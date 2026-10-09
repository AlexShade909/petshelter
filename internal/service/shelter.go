package service

type shelter struct{}

func NewShelter() shelter {
	return shelter{}
}

//func (s shelter) FullInfo() (map[int]models.Shelter, error) {
//	return sheltersData, nil
//}
//
//func (s shelter) GetByID(shelterNumber int) (models.Shelter, error) {
//	shelterInfo, ok := sheltersData[shelterNumber]
//	if !ok {
//		return models.Shelter{}, errors.New("sheler not found")
//	}
//	return shelterInfo, nil
//}
//
//func (s shelter) ListDogs(shelterNumber int) ([]string, error) {
//	listDogs := make([]string, 0)
//	sh, ok := sheltersData[shelterNumber]
//	if !ok {
//		return []string{}, errors.New("number not correct")
//	}
//	for _, d := range dogsData {
//		if d.Shelter == sh {
//			listDogs = append(listDogs, d.Nickname)
//		}
//	}
//	return listDogs, nil
//}
//
//func (s shelter) Create(shelterCreate models.Shelter) error {
//	sheltersData[nextShelterID] = shelterCreate
//	nextShelterID++
//	return nil
//}
//
//func (s shelter) Delete(shelterNumber int) error {
//	if _, ok := sheltersData[shelterNumber]; !ok {
//		return errors.New("clinic not found")
//	}
//	delete(sheltersData, shelterNumber)
//	return nil
//}
//
//func (s shelter) Update(shelterNumber int, patch models.ShelterPatch) (models.Shelter, error) {
//	sh, ok := sheltersData[shelterNumber]
//	if !ok {
//		return models.Shelter{}, errors.New("shelter not found")
//	}
//
//	if patch.Address != "" {
//		sh.Address = patch.Address
//	}
//	if patch.PhoneNumber != "" {
//		sh.PhoneNumber = patch.PhoneNumber
//	}
//	if patch.WorkingTime != "" {
//		sh.WorkingTime = patch.WorkingTime
//	}
//
//	sheltersData[shelterNumber] = sh
//	return sheltersData[shelterNumber], nil
//}
//
//func (s shelter) Replace(shelterNumber int, patch models.Shelter) (models.Shelter, error) {
//	if _, ok := sheltersData[shelterNumber]; !ok {
//		return models.Shelter{}, errors.New("shelter not found")
//	}
//	sheltersData[shelterNumber] = patch
//	return sheltersData[shelterNumber], nil
//}
