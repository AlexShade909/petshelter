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
	clinicInfo, ok := clinicsData[clinicNumber]
	if !ok {
		return models.Clinic{}, errors.New("clinic not found")
	}

	return clinicInfo, nil
}

func (c clinic) ListDogs(clinicNumber int) ([]string, error) {
	listDogs := make([]string, 0)
	cl, ok := clinicsData[clinicNumber]
	if !ok {
		return []string{}, errors.New("number not correct")
	}
	for _, d := range dogsData {
		if d.Сlinic == cl {
			listDogs = append(listDogs, d.Nickname)
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

func (c clinic) Replace(clinicNumber int, patch models.Clinic) (models.Clinic, error) {
	if _, ok := clinicsData[clinicNumber]; !ok {
		return models.Clinic{}, errors.New("patch not found")
	}
	clinicsData[clinicNumber] = patch
	return patch, nil
}
