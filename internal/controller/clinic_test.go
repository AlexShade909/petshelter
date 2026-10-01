package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"petshelter/internal/models"
)

// =========================
// Mock service
// =========================

type mockClinicService struct {
	fullInfoFunc    func() (map[int]models.Clinic, error)
	infoFunc        func(int) (models.Clinic, error)
	listClinicsFunc func(int) ([]string, error)
	createFunc      func(models.Clinic) error
	deleteFunc      func(int) error
	updateFunc      func(int, models.ClinicPatch) (models.Clinic, error)
	replaceFunc     func(int, models.Clinic) (models.Clinic, error)
}

func (m *mockClinicService) FullInfo() (map[int]models.Clinic, error) {
	return m.fullInfoFunc()
}

func (m *mockClinicService) Info(clinicNumber int) (models.Clinic, error) {
	return m.infoFunc(clinicNumber)
}

func (m *mockClinicService) ListClinics(clinicNumber int) ([]string, error) {
	return m.listClinicsFunc(clinicNumber)
}

func (m *mockClinicService) Create(clinic models.Clinic) error {
	return m.createFunc(clinic)
}

func (m *mockClinicService) Delete(clinicNumber int) error {
	return m.deleteFunc(clinicNumber)
}

func (m *mockClinicService) Update(
	clinicNumber int,
	patch models.ClinicPatch,
) (models.Clinic, error) {
	return m.updateFunc(clinicNumber, patch)
}

func (m *mockClinicService) Replace(
	clinicNumber int,
	clinic models.Clinic,
) (models.Clinic, error) {
	return m.replaceFunc(clinicNumber, clinic)
}

// =========================
// FullInfoHandler
// =========================

func TestClinic_FullInfoHandler(t *testing.T) {
	expected := map[int]models.Clinic{
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

	service := &mockClinicService{
		fullInfoFunc: func() (map[int]models.Clinic, error) {
			return expected, nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(http.MethodGet, "/clinics", nil)
	rec := httptest.NewRecorder()

	controller.FullInfoHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rec.Code)
	}
}

// =========================
// FullInfoHandler - error
// =========================

func TestClinic_FullInfoHandler_Error(t *testing.T) {
	service := &mockClinicService{
		fullInfoFunc: func() (map[int]models.Clinic, error) {
			return nil, errors.New("database error")
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(http.MethodGet, "/clinics", nil)
	rec := httptest.NewRecorder()

	controller.FullInfoHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// InfoHandler
// =========================

func TestClinic_InfoHandler(t *testing.T) {
	expectedClinic := models.Clinic{
		Address:     "Мира 1",
		PhoneNumber: "+375291111111",
		WorkingTime: "10:00-23:00",
	}

	service := &mockClinicService{
		infoFunc: func(clinicNumber int) (models.Clinic, error) {
			if clinicNumber != 1 {
				t.Fatalf("expected clinic number 1, got %d", clinicNumber)
			}

			return expectedClinic, nil
		},

		listClinicsFunc: func(clinicNumber int) ([]string, error) {
			if clinicNumber != 1 {
				t.Fatalf("expected clinic number 1, got %d", clinicNumber)
			}

			return []string{"clinic-1", "clinic-2"}, nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/clinics/1",
		nil,
	)

	// PathValue() работает только если значение
	// было установлено через SetPathValue.
	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.InfoHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rec.Code)
	}
}

// =========================
// InfoHandler - Info error
// =========================

func TestClinic_InfoHandler_InfoError(t *testing.T) {
	service := &mockClinicService{
		infoFunc: func(clinicNumber int) (models.Clinic, error) {
			return models.Clinic{}, errors.New("clinic not found")
		},

		listClinicsFunc: func(clinicNumber int) ([]string, error) {
			t.Fatal("ListClinics should not be called")
			return nil, nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/clinics/1",
		nil,
	)

	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.InfoHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// InfoHandler - ListClinics error
// =========================

func TestClinic_InfoHandler_ListClinicsError(t *testing.T) {
	service := &mockClinicService{
		infoFunc: func(clinicNumber int) (models.Clinic, error) {
			return models.Clinic{
				Address:     "Мира 1",
				PhoneNumber: "+375291111111",
				WorkingTime: "10:00-23:00",
			}, nil
		},

		listClinicsFunc: func(clinicNumber int) ([]string, error) {
			return nil, errors.New("cannot get clinic list")
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/clinics/1",
		nil,
	)

	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.InfoHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// CreateHandler
// =========================

func TestClinic_CreateHandler(t *testing.T) {
	body := `{
"Address": "Мира 1",
"PhoneNumber": "+375291111111",
"WorkingTime": "10:00-23:00"
}`

	service := &mockClinicService{
		createFunc: func(clinic models.Clinic) error {
			if clinic.Address != "Мира 1" {
				t.Errorf("unexpected address: %s", clinic.Address)
			}

			if clinic.PhoneNumber != "+375291111111" {
				t.Errorf("unexpected phone: %s", clinic.PhoneNumber)
			}

			if clinic.WorkingTime != "10:00-23:00" {
				t.Errorf("unexpected working time: %s",
					clinic.WorkingTime)
			}

			return nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/clinics",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	controller.CreateHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d",
			http.StatusCreated, rec.Code)
	}
}

// =========================
// CreateHandler - invalid JSON
// =========================

func TestClinic_CreateHandler_InvalidJSON(t *testing.T) {
	service := &mockClinicService{
		createFunc: func(clinic models.Clinic) error {
			t.Fatal("Create should not be called")
			return nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/clinics",
		strings.NewReader(`invalid json`),
	)

	rec := httptest.NewRecorder()

	controller.CreateHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// CreateHandler - service error
// =========================

func TestClinic_CreateHandler_ServiceError(t *testing.T) {
	service := &mockClinicService{
		createFunc: func(clinic models.Clinic) error {
			return errors.New("create error")
		},
	}

	controller := NewClinic(service)

	body := `{
"Address": "Мира 1",
"PhoneNumber": "+375291111111",
"WorkingTime": "10:00-23:00"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/clinics",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	controller.CreateHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// DeleteHandler
// =========================

func TestClinic_DeleteHandler(t *testing.T) {
	service := &mockClinicService{
		deleteFunc: func(clinicNumber int) error {
			if clinicNumber != 5 {
				t.Fatalf("expected clinic number 5, got %d",
					clinicNumber)
			}

			return nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/clinics/5",
		nil,
	)

	req.SetPathValue("NumberClinic", "5")

	rec := httptest.NewRecorder()

	controller.DeleteHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rec.Code)
	}

	expected := "data in clinic has been deleted. Number deleted clinic: 5"

	if !strings.Contains(rec.Body.String(), expected) {
		t.Fatalf("expected response to contain %q, got %q",
			expected, rec.Body.String())
	}
}

// =========================
// DeleteHandler - error
// =========================

func TestClinic_DeleteHandler_Error(t *testing.T) {
	service := &mockClinicService{
		deleteFunc: func(clinicNumber int) error {
			return errors.New("clinic not found")
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/clinics/5",
		nil,
	)

	req.SetPathValue("NumberClinic", "5")

	rec := httptest.NewRecorder()

	controller.DeleteHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// UpdateHandler
// =========================

func TestClinic_UpdateHandler(t *testing.T) {
	body := `{
"Address": "Новая 10",
"PhoneNumber": "+375293333333",
"WorkingTime": "08:00-20:00"
}`

	service := &mockClinicService{
		updateFunc: func(
			clinicNumber int,
			patch models.ClinicPatch,
		) (models.Clinic, error) {

			if clinicNumber != 2 {
				t.Fatalf("expected clinic number 2, got %d",
					clinicNumber)
			}

			if patch.Address != "Новая 10" {
				t.Errorf("unexpected address: %s",
					patch.Address)
			}

			if patch.PhoneNumber != "+375293333333" {
				t.Errorf("unexpected phone: %s",
					patch.PhoneNumber)
			}

			if patch.WorkingTime != "08:00-20:00" {
				t.Errorf("unexpected working time: %s",
					patch.WorkingTime)
			}

			return models.Clinic{
				Address:     "Новая 10",
				PhoneNumber: "+375293333333",
				WorkingTime: "08:00-20:00",
			}, nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/clinics/2",
		strings.NewReader(body),
	)

	req.SetPathValue("NumberClinic", "2")

	rec := httptest.NewRecorder()

	controller.UpdateHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rec.Code)
	}
}

// =========================
// UpdateHandler - invalid JSON
// =========================

func TestClinic_UpdateHandler_InvalidJSON(t *testing.T) {
	service := &mockClinicService{
		updateFunc: func(
			clinicNumber int,
			patch models.ClinicPatch,
		) (models.Clinic, error) {

			t.Fatal("Update should not be called")

			return models.Clinic{}, nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/clinics/1",
		strings.NewReader(`invalid json`),
	)

	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.UpdateHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// UpdateHandler - unknown field
// =========================

func TestClinic_UpdateHandler_UnknownField(t *testing.T) {
	service := &mockClinicService{
		updateFunc: func(
			clinicNumber int,
			patch models.ClinicPatch,
		) (models.Clinic, error) {

			t.Fatal("Update should not be called")

			return models.Clinic{}, nil
		},
	}

	controller := NewClinic(service)

	body := `{
"UnknownField": "test"
}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/clinics/1",
		strings.NewReader(body),
	)

	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.UpdateHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// UpdateHandler - service error
// =========================

func TestClinic_UpdateHandler_ServiceError(t *testing.T) {
	service := &mockClinicService{
		updateFunc: func(
			clinicNumber int,
			patch models.ClinicPatch,
		) (models.Clinic, error) {

			return models.Clinic{}, errors.New("clinic not found")
		},
	}

	controller := NewClinic(service)

	body := `{
"Address": "Новая 10"
}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/clinics/1",
		strings.NewReader(body),
	)

	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.UpdateHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// ReplaceHandler
// =========================

func TestClinic_ReplaceHandler(t *testing.T) {
	body := `{
"Address": "Ленина 100",
"PhoneNumber": "+375294444444",
"WorkingTime": "09:00-21:00"
}`

	service := &mockClinicService{
		replaceFunc: func(
			clinicNumber int,
			clinic models.Clinic,
		) (models.Clinic, error) {

			if clinicNumber != 3 {
				t.Fatalf("expected clinic number 3, got %d",
					clinicNumber)
			}

			if clinic.Address != "Ленина 100" {
				t.Errorf("unexpected address: %s",
					clinic.Address)
			}

			return clinic, nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodPut,
		"/clinics/3",
		strings.NewReader(body),
	)

	req.SetPathValue("NumberClinic", "3")

	rec := httptest.NewRecorder()

	controller.ReplaceHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rec.Code)
	}
}

// =========================
// ReplaceHandler - invalid JSON
// =========================

func TestClinic_ReplaceHandler_InvalidJSON(t *testing.T) {
	service := &mockClinicService{
		replaceFunc: func(
			clinicNumber int,
			clinic models.Clinic,
		) (models.Clinic, error) {

			t.Fatal("Replace should not be called")

			return models.Clinic{}, nil
		},
	}

	controller := NewClinic(service)

	req := httptest.NewRequest(
		http.MethodPut,
		"/clinics/1",
		strings.NewReader(`invalid json`),
	)

	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.ReplaceHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

// =========================
// ReplaceHandler - service error
// =========================

func TestClinic_ReplaceHandler_ServiceError(t *testing.T) {
	service := &mockClinicService{
		replaceFunc: func(
			clinicNumber int,
			clinic models.Clinic,
		) (models.Clinic, error) {

			return models.Clinic{}, errors.New("clinic not found")
		},
	}

	controller := NewClinic(service)

	body := `{
"Address": "Ленина 100",
"PhoneNumber": "+375294444444",
"WorkingTime": "09:00-21:00"
}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/clinics/1",
		strings.NewReader(body),
	)

	req.SetPathValue("NumberClinic", "1")

	rec := httptest.NewRecorder()

	controller.ReplaceHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}
