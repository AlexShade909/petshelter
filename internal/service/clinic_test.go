package service

import (
	"errors"
	"testing"

	"petshelter/internal/models"
)

// =========================
// Восстановление данных
// =========================

func backupClinicData(t *testing.T) {
	t.Helper()

	oldClinicsData := clinicsData
	oldDogsData := dogsData
	oldNextClinicID := nextClinicID

	t.Cleanup(func() {
		clinicsData = oldClinicsData
		dogsData = oldDogsData
		nextClinicID = oldNextClinicID
	})
}

// =========================
// FullInfo
// =========================

func TestClinic_FullInfo(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{
		0: {
			Address:     "Мира 1",
			PhoneNumber: "+375291111111",
			WorkingTime: "10:00-23:00",
		},
		1: {
			Address:     "Ленина 10",
			PhoneNumber: "+375292222222",
			WorkingTime: "09:00-18:00",
		},
	}

	service := NewClinic()

	data, err := service.FullInfo()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 2 {
		t.Fatalf("expected 2 clinics, got %d", len(data))
	}

	if data[0].Address != "Мира 1" {
		t.Errorf("unexpected address: %s", data[0].Address)
	}

	if data[1].PhoneNumber != "+375292222222" {
		t.Errorf("unexpected phone: %s", data[1].PhoneNumber)
	}
}

// =========================
// Info
// =========================

func TestClinic_Info(t *testing.T) {
	backupClinicData(t)

	expected := models.Clinic{
		Address:     "Мира 1",
		PhoneNumber: "+375291111111",
		WorkingTime: "10:00-23:00",
	}

	clinicsData = map[int]models.Clinic{
		1: expected,
	}

	service := NewClinic()

	data, err := service.Info(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != expected {
		t.Errorf("expected %+v, got %+v", expected, data)
	}
}

// =========================
// Info - not found
// =========================

func TestClinic_Info_NotFound(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{}

	service := NewClinic()

	_, err := service.Info(1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, errors.New("clinic not found")) &&
		err.Error() != "clinic not found" {
		t.Errorf("unexpected error: %v", err)
	}
}

// =========================
// ListClinics
// =========================

func TestClinic_ListClinics(t *testing.T) {
	backupClinicData(t)

	clinic := models.Clinic{
		Address:     "Мира 1",
		PhoneNumber: "+375291111111",
		WorkingTime: "10:00-23:00",
	}

	clinicsData = map[int]models.Clinic{
		0: clinic,
	}

	// ВАЖНО:
	// Здесь используется твоя структура Dog.
	// Поле должно называться Сlinic так же,
	// как в production-коде.
	dogsData = map[string]models.Dog{
		"Барсик": {
			Сlinic: clinic,
		},
		"Шарик": {
			Сlinic: clinic,
		},
		"Рекс": {
			Сlinic: models.Clinic{
				Address:     "Ленина 10",
				PhoneNumber: "+375292222222",
				WorkingTime: "09:00-18:00",
			},
		},
	}

	service := NewClinic()

	data, err := service.ListDogs(0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 2 {
		t.Fatalf("expected 2 dogs, got %d", len(data))
	}

	found := map[string]bool{}

	for _, name := range data {
		found[name] = true
	}

	if !found["Барсик"] {
		t.Error("Барсик should be in clinic list")
	}

	if !found["Шарик"] {
		t.Error("Шарик should be in clinic list")
	}

	if found["Рекс"] {
		t.Error("Рекс should not be in clinic list")
	}
}

// =========================
// ListClinics - invalid number
// =========================

func TestClinic_ListClinics_InvalidNumber(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{
		0: {
			Address:     "Мира 1",
			PhoneNumber: "+375291111111",
			WorkingTime: "10:00-23:00",
		},
	}

	service := NewClinic()

	_, err := service.ListDogs(-1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "number not correct" {
		t.Errorf(
			"expected error %q, got %q",
			"number not correct",
			err.Error(),
		)
	}
}

func TestClinic_ListClinics_NumberTooLarge(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{
		0: {
			Address:     "Мира 1",
			PhoneNumber: "+375291111111",
			WorkingTime: "10:00-23:00",
		},
	}

	service := NewClinic()

	_, err := service.ListDogs(1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "number not correct" {
		t.Errorf(
			"expected error %q, got %q",
			"number not correct",
			err.Error(),
		)
	}
}

// =========================
// Create
// =========================

func TestClinic_Create(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{}
	nextClinicID = 5

	newClinic := models.Clinic{
		Address:     "Новая 10",
		PhoneNumber: "+375293333333",
		WorkingTime: "08:00-20:00",
	}

	service := NewClinic()

	err := service.Create(newClinic)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	created, ok := clinicsData[5]

	if !ok {
		t.Fatal("clinic was not created")
	}

	if created != newClinic {
		t.Errorf(
			"expected %+v, got %+v",
			newClinic,
			created,
		)
	}

	if nextClinicID != 6 {
		t.Errorf(
			"expected nextClinicID 6, got %d",
			nextClinicID,
		)
	}
}

// =========================
// Delete
// =========================

func TestClinic_Delete(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{
		1: {
			Address:     "Мира 1",
			PhoneNumber: "+375291111111",
			WorkingTime: "10:00-23:00",
		},
		2: {
			Address:     "Ленина 10",
			PhoneNumber: "+375292222222",
			WorkingTime: "09:00-18:00",
		},
	}

	service := NewClinic()

	err := service.Delete(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := clinicsData[1]; ok {
		t.Error("clinic should have been deleted")
	}

	if _, ok := clinicsData[2]; !ok {
		t.Error("clinic 2 should still exist")
	}
}

// =========================
// Delete - not found
// =========================

func TestClinic_Delete_NotFound(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{}

	service := NewClinic()

	err := service.Delete(1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "clinic not found" {
		t.Errorf(
			"expected error %q, got %q",
			"clinic not found",
			err.Error(),
		)
	}
}

// =========================
// Update
// =========================

func TestClinic_Update(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{
		1: {
			Address:     "Мира 1",
			PhoneNumber: "+375291111111",
			WorkingTime: "10:00-23:00",
		},
	}

	patch := models.ClinicPatch{
		Address:     "Новая 20",
		PhoneNumber: "+375294444444",
		WorkingTime: "08:00-20:00",
	}

	service := NewClinic()

	data, err := service.Update(1, patch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := models.Clinic{
		Address:     "Новая 20",
		PhoneNumber: "+375294444444",
		WorkingTime: "08:00-20:00",
	}

	if data != expected {
		t.Errorf(
			"expected %+v, got %+v",
			expected,
			data,
		)
	}

	if clinicsData[1] != expected {
		t.Errorf(
			"data in map was not updated: %+v",
			clinicsData[1],
		)
	}
}

// =========================
// Update - partial patch
// =========================

func TestClinic_Update_PartialPatch(t *testing.T) {
	backupClinicData(t)

	original := models.Clinic{
		Address:     "Мира 1",
		PhoneNumber: "+375291111111",
		WorkingTime: "10:00-23:00",
	}

	clinicsData = map[int]models.Clinic{
		1: original,
	}

	patch := models.ClinicPatch{
		Address: "Новая 20",
	}

	service := NewClinic()

	data, err := service.Update(1, patch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Address != "Новая 20" {
		t.Errorf(
			"expected address %q, got %q",
			"Новая 20",
			data.Address,
		)
	}

	if data.PhoneNumber != original.PhoneNumber {
		t.Errorf(
			"phone number should not change, got %q",
			data.PhoneNumber,
		)
	}

	if data.WorkingTime != original.WorkingTime {
		t.Errorf(
			"working time should not change, got %q",
			data.WorkingTime,
		)
	}
}

// =========================
// Update - not found
// =========================

func TestClinic_Update_NotFound(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{}

	service := NewClinic()

	patch := models.ClinicPatch{
		Address: "Новая 20",
	}

	_, err := service.Update(1, patch)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "clinic not found" {
		t.Errorf(
			"expected error %q, got %q",
			"clinic not found",
			err.Error(),
		)
	}
}

// =========================
// Replace
// =========================

func TestClinic_Replace(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{
		1: {
			Address:     "Мира 1",
			PhoneNumber: "+375291111111",
			WorkingTime: "10:00-23:00",
		},
	}

	replacement := models.Clinic{
		Address:     "Ленина 100",
		PhoneNumber: "+375295555555",
		WorkingTime: "09:00-21:00",
	}

	service := NewClinic()

	data, err := service.Replace(1, replacement)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != replacement {
		t.Errorf(
			"expected %+v, got %+v",
			replacement,
			data,
		)
	}

	if clinicsData[1] != replacement {
		t.Errorf(
			"clinic was not replaced: %+v",
			clinicsData[1],
		)
	}
}

// =========================
// Replace - new clinic
// =========================

func TestClinic_Replace_NewClinic(t *testing.T) {
	backupClinicData(t)

	clinicsData = map[int]models.Clinic{}

	newClinic := models.Clinic{
		Address:     "Ленина 100",
		PhoneNumber: "+375295555555",
		WorkingTime: "09:00-21:00",
	}

	service := NewClinic()

	data, err := service.Replace(10, newClinic)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != newClinic {
		t.Errorf(
			"expected %+v, got %+v",
			newClinic,
			data,
		)
	}

	if clinicsData[10] != newClinic {
		t.Error("new clinic was not added")
	}
}
