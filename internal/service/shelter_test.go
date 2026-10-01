package service

import (
	"testing"

	"petshelter/internal/models"
)

// =====================================================
// Helpers
// =====================================================

// Сохраняем исходное состояние глобальных данных,
// чтобы один тест не влиял на другой.
func backupShelterData(t *testing.T) func() {
	t.Helper()

	oldSheltersData := sheltersData
	oldDogsData := dogsData
	oldNextShelterID := nextShelterID

	t.Cleanup(func() {
		sheltersData = oldSheltersData
		dogsData = oldDogsData
		nextShelterID = oldNextShelterID
	})

	return func() {}
}

// =====================================================
// FullInfo
// =====================================================

func TestShelter_FullInfo(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{
		1: {
			Address:     "Mira 1",
			PhoneNumber: "1",
			WorkingTime: "10:00-18:00",
		},
		2: {
			Address:     "Lenina 10",
			PhoneNumber: "2",
			WorkingTime: "09:00-17:00",
		},
	}

	service := NewShelter()

	result, err := service.FullInfo()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 shelters, got %d", len(result))
	}

	if result[1].Address != "Mira 1" {
		t.Fatalf(
			"expected address %q, got %q",
			"Mira 1",
			result[1].Address,
		)
	}

	if result[2].Address != "Lenina 10" {
		t.Fatalf(
			"expected address %q, got %q",
			"Lenina 10",
			result[2].Address,
		)
	}
}

// =====================================================
// Info
// =====================================================

func TestShelter_Info(t *testing.T) {
	backupShelterData(t)

	expected := models.Shelter{
		Address:     "Mira 1",
		PhoneNumber: "1",
		WorkingTime: "10:00-18:00",
	}

	sheltersData = map[int]models.Shelter{
		1: expected,
	}

	service := NewShelter()

	result, err := service.Info(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Fatalf(
			"expected %+v, got %+v",
			expected,
			result,
		)
	}
}

func TestShelter_Info_NotFound(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{}

	service := NewShelter()

	result, err := service.Info(999)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	// В текущей реализации Info не возвращает ошибку,
	// если такого ключа нет. Возвращается zero-value Shelter.
	if result != (models.Shelter{}) {
		t.Fatalf(
			"expected empty shelter, got %+v",
			result,
		)
	}
}

// =====================================================
// ListDogs
// =====================================================

func TestShelter_ListDogs(t *testing.T) {
	backupShelterData(t)

	shelter1 := models.Shelter{
		Address:     "Mira 1",
		PhoneNumber: "1",
		WorkingTime: "10:00-18:00",
	}

	shelter2 := models.Shelter{
		Address:     "Lenina 10",
		PhoneNumber: "2",
		WorkingTime: "09:00-17:00",
	}

	sheltersData = map[int]models.Shelter{
		1: shelter1,
		2: shelter2,
	}

	dogsData = map[string]models.Dog{
		"Bobik": {
			Shelter: shelter1,
		},
		"Sharik": {
			Shelter: shelter1,
		},
		"Rex": {
			Shelter: shelter2,
		},
	}

	service := NewShelter()

	result, err := service.ListDogs(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 dogs, got %d",
			len(result),
		)
	}

	// Проверяем, что именно нужные собаки попали
	// в результат.
	found := map[string]bool{}

	for _, name := range result {
		found[name] = true
	}

	if !found["Bobik"] {
		t.Error("expected Bobik in result")
	}

	if !found["Sharik"] {
		t.Error("expected Sharik in result")
	}

	if found["Rex"] {
		t.Error("Rex should not be in result")
	}
}

func TestShelter_ListDogs_InvalidNumber(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{
		0: {},
		1: {},
	}

	service := NewShelter()

	tests := []struct {
		name   string
		number int
	}{
		{
			name:   "negative number",
			number: -1,
		},
		{
			name:   "number greater than length",
			number: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ListDogs(tt.number)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != "number not correct" {
				t.Fatalf(
					"expected %q, got %q",
					"number not correct",
					err.Error(),
				)
			}

			if len(result) != 0 {
				t.Fatalf(
					"expected empty result, got %v",
					result,
				)
			}
		})
	}
}

// =====================================================
// Create
// =====================================================

func TestShelter_Create(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{}

	nextShelterID = 1

	service := NewShelter()

	newShelter := models.Shelter{
		Address:     "Mira 1",
		PhoneNumber: "1",
		WorkingTime: "10:00-18:00",
	}

	err := service.Create(newShelter)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, ok := sheltersData[1]

	if !ok {
		t.Fatal("shelter was not created")
	}

	if result != newShelter {
		t.Fatalf(
			"expected %+v, got %+v",
			newShelter,
			result,
		)
	}

	if nextShelterID != 2 {
		t.Fatalf(
			"expected nextShelterID 2, got %d",
			nextShelterID,
		)
	}
}

// =====================================================
// Delete
// =====================================================

func TestShelter_Delete(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{
		1: {
			Address: "Mira 1",
		},
		2: {
			Address: "Lenina 10",
		},
	}

	service := NewShelter()

	err := service.Delete(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := sheltersData[1]; ok {
		t.Fatal("shelter 1 was not deleted")
	}

	if _, ok := sheltersData[2]; !ok {
		t.Fatal("shelter 2 should still exist")
	}
}

func TestShelter_Delete_NotFound(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{
		1: {},
	}

	service := NewShelter()

	// Текущая реализация Delete не проверяет,
	// существует ли shelter.
	err := service.Delete(999)

	if err != nil {
		t.Fatalf(
			"expected nil error, got %v",
			err,
		)
	}
}

// =====================================================
// Update
// =====================================================

func TestShelter_Update(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{
		1: {
			Address:     "Old address",
			PhoneNumber: "1",
			WorkingTime: "10:00-18:00",
		},
	}

	service := NewShelter()

	patch := models.ShelterPatch{
		Address:     "New address",
		PhoneNumber: "10",
		WorkingTime: "09:00-20:00",
	}

	result, err := service.Update(1, patch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Address != "New address" {
		t.Fatalf(
			"expected address %q, got %q",
			"New address",
			result.Address,
		)
	}

	if result.PhoneNumber != "10" {
		t.Fatalf(
			"expected number %q, got %q",
			"10",
			result.PhoneNumber,
		)
	}

	if result.WorkingTime != "09:00-20:00" {
		t.Fatalf(
			"expected working time %q, got %q",
			"09:00-20:00",
			result.WorkingTime,
		)
	}
}

func TestShelter_Update_PartialPatch(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{
		1: {
			Address:     "Old address",
			PhoneNumber: "1",
			WorkingTime: "10:00-18:00",
		},
	}

	service := NewShelter()

	patch := models.ShelterPatch{
		Address: "New address",
	}

	result, err := service.Update(1, patch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Address != "New address" {
		t.Fatalf(
			"expected address %q, got %q",
			"New address",
			result.Address,
		)
	}

	// Эти поля должны остаться неизменными.
	if result.PhoneNumber != "1" {
		t.Fatalf(
			"expected number %q, got %q",
			"1",
			result.PhoneNumber,
		)
	}

	if result.WorkingTime != "10:00-18:00" {
		t.Fatalf(
			"expected working time %q, got %q",
			"10:00-18:00",
			result.WorkingTime,
		)
	}
}

func TestShelter_Update_NotFound(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{}

	service := NewShelter()

	patch := models.ShelterPatch{
		Address: "New address",
	}

	result, err := service.Update(999, patch)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "shelter not found" {
		t.Fatalf(
			"expected %q, got %q",
			"shelter not found",
			err.Error(),
		)
	}

	if result != (models.Shelter{}) {
		t.Fatalf(
			"expected empty shelter, got %+v",
			result,
		)
	}
}

// =====================================================
// Replace
// =====================================================

func TestShelter_Replace(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{
		1: {
			Address:     "Old address",
			PhoneNumber: "1",
			WorkingTime: "10:00-18:00",
		},
	}

	service := NewShelter()

	newShelter := models.Shelter{
		Address:     "New address",
		PhoneNumber: "2",
		WorkingTime: "09:00-20:00",
	}

	result, err := service.Replace(1, newShelter)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != newShelter {
		t.Fatalf(
			"expected %+v, got %+v",
			newShelter,
			result,
		)
	}

	if sheltersData[1] != newShelter {
		t.Fatalf(
			"data was not replaced: got %+v",
			sheltersData[1],
		)
	}
}

func TestShelter_Replace_NewShelter(t *testing.T) {
	backupShelterData(t)

	sheltersData = map[int]models.Shelter{}

	service := NewShelter()

	newShelter := models.Shelter{
		Address:     "New address",
		PhoneNumber: "1",
		WorkingTime: "10:00-20:00",
	}

	result, err := service.Replace(5, newShelter)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != newShelter {
		t.Fatalf(
			"expected %+v, got %+v",
			newShelter,
			result,
		)
	}

	if sheltersData[5] != newShelter {
		t.Fatalf(
			"expected shelter at key 5, got %+v",
			sheltersData[5],
		)
	}
}
