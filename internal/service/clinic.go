package service

import (
	"errors"
	"petshelter/internal/models"
)

type clinic struct{}

func NewClinic() clinic {
	return clinic{}
}

func (c clinic) FullInfo() (map[int]models.Clinic, error) {
	return clinicsData, nil
}

func (c clinic) Info(clinicNumber int) (models.Clinic, error) {
	clinic, ok := clinicsData[clinicNumber]
	if !ok {
		return models.Clinic{}, errors.New("clinic not found")
	}

	return clinic, nil
}

func (c clinic) ListDogs(clinicNumber int) ([]string, error) {
	var listDogs []string
	if clinicNumber >= len(clinicsData) || clinicNumber < 0 {
		return []string{}, errors.New("number not correct")
	}
	for i, v := range dogsData {
		if v.Сlinic == clinicsData[clinicNumber] {
			listDogs = append(listDogs, i)
		}
	}
	return listDogs, nil
}

func (c clinic) Create(clinicCreate models.Clinic) error {
	clinicsData[nextClinicID] = clinicCreate
	nextClinicID++

	return nil
}

func (c clinic) Delete(clinicNumber int) error {
	if _, ok := clinicsData[clinicNumber]; !ok {
		return errors.New("clinic not found")
	}
	delete(clinicsData, clinicNumber)
	return nil
}

func (c clinic) Update(clinicNumber int, patch models.ClinicPatch) (models.Clinic, error) {

	cl, ok := clinicsData[clinicNumber]
	if !ok {
		return models.Clinic{}, errors.New("clinic not found")
	}

	if patch.Address != "" {
		cl.Address = patch.Address
	}

	if patch.PhoneNumber != "" {
		cl.PhoneNumber = patch.PhoneNumber
	}

	if patch.WorkingTime != "" {
		cl.WorkingTime = patch.WorkingTime
	}

	clinicsData[clinicNumber] = cl

	return cl, nil
}

func (c clinic) Replace(clinicNumber int, clinic models.Clinic) (models.Clinic, error) {
	clinicsData[clinicNumber] = clinic
	return clinic, nil
}
