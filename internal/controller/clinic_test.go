package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"petshelter/internal/models"
)

// --- мок сервиса: у каждого метода своя функция, задаём только нужные ---

type mockClinicService struct {
	fullInfoFn func() (map[int]models.Clinic, error)
	infoFn     func(int) (models.Clinic, error)
	listDogsFn func(int) ([]string, error)
	createFn   func(models.Clinic) error
	deleteFn   func(int) error
	updateFn   func(int, models.ClinicPatch) (models.Clinic, error)
	replaceFn  func(int, models.Clinic) (models.Clinic, error)
}

func (m mockClinicService) FullInfo() (map[int]models.Clinic, error) { return m.fullInfoFn() }
func (m mockClinicService) Info(n int) (models.Clinic, error)        { return m.infoFn(n) }
func (m mockClinicService) ListDogs(n int) ([]string, error)         { return m.listDogsFn(n) }
func (m mockClinicService) Create(c models.Clinic) error             { return m.createFn(c) }
func (m mockClinicService) Delete(n int) error                       { return m.deleteFn(n) }
func (m mockClinicService) Update(n int, p models.ClinicPatch) (models.Clinic, error) {
	return m.updateFn(n, p)
}
func (m mockClinicService) Replace(n int, c models.Clinic) (models.Clinic, error) {
	return m.replaceFn(n, c)
}

// --- хелпер: запрос с path-параметром (нужен Go 1.22+) ---

func newReq(method, body, number string) *http.Request {
	req := httptest.NewRequest(method, "/clinics/"+number, strings.NewReader(body))
	req.SetPathValue("NumberClinic", number)
	return req
}

var errService = errors.New("service error")

func TestFullInfoHandler(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (map[int]models.Clinic, error)
		wantCode int
	}{
		{"ok", func() (map[int]models.Clinic, error) { return map[int]models.Clinic{1: {}}, nil }, http.StatusOK},
		{"service error", func() (map[int]models.Clinic, error) { return nil, errService }, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewClinic(mockClinicService{fullInfoFn: tt.fn})
			rec := httptest.NewRecorder()

			h.FullInfoHandler(rec, newReq(http.MethodGet, "", "1"))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func TestInfoHandler(t *testing.T) {
	tests := []struct {
		name     string
		svc      mockClinicService
		wantCode int
	}{
		{
			name: "ok",
			svc: mockClinicService{
				infoFn:     func(int) (models.Clinic, error) { return models.Clinic{}, nil },
				listDogsFn: func(int) ([]string, error) { return []string{"Rex", "Bim"}, nil },
			},
			wantCode: http.StatusOK,
		},
		{
			name: "info error",
			svc: mockClinicService{
				infoFn: func(int) (models.Clinic, error) { return models.Clinic{}, errService },
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "list dogs error",
			svc: mockClinicService{
				infoFn:     func(int) (models.Clinic, error) { return models.Clinic{}, nil },
				listDogsFn: func(int) ([]string, error) { return nil, errService },
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewClinic(tt.svc)
			rec := httptest.NewRecorder()

			h.InfoHandler(rec, newReq(http.MethodGet, "", "1"))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
		})
	}
}

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		createErr  error
		wantCode   int
		wantCalled bool
	}{
		{"ok", `{}`, nil, http.StatusCreated, true},
		{"invalid json", `{broken`, nil, http.StatusBadRequest, false},
		{"service error", `{}`, errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			h := NewClinic(mockClinicService{
				createFn: func(models.Clinic) error {
					called = true
					return tt.createErr
				},
			})
			rec := httptest.NewRecorder()

			h.CreateHandler(rec, newReq(http.MethodPost, tt.body, ""))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name      string
		deleteErr error
		wantCode  int
	}{
		{"ok", nil, http.StatusOK},
		{"service error", errService, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotNumber int
			h := NewClinic(mockClinicService{
				deleteFn: func(n int) error {
					gotNumber = n
					return tt.deleteErr
				},
			})
			rec := httptest.NewRecorder()

			h.DeleteHandler(rec, newReq(http.MethodDelete, "", "7"))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if gotNumber != 7 {
				t.Errorf("clinic number = %d, want 7", gotNumber)
			}
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		updateErr  error
		wantCode   int
		wantCalled bool
	}{
		{"ok", `{}`, nil, http.StatusOK, true},
		{"unknown field", `{"no_such_field": 1}`, nil, http.StatusBadRequest, false},
		{"invalid json", `{broken`, nil, http.StatusBadRequest, false},
		{"service error", `{}`, errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			h := NewClinic(mockClinicService{
				updateFn: func(int, models.ClinicPatch) (models.Clinic, error) {
					called = true
					return models.Clinic{}, tt.updateErr
				},
			})
			rec := httptest.NewRecorder()

			h.UpdateHandler(rec, newReq(http.MethodPatch, tt.body, "1"))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
		})
	}
}

func TestReplaceHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		replaceErr error
		wantCode   int
		wantCalled bool
	}{
		{"ok", `{}`, nil, http.StatusOK, true},
		{"invalid json", `{broken`, nil, http.StatusBadRequest, false},
		{"service error", `{}`, errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			h := NewClinic(mockClinicService{
				replaceFn: func(int, models.Clinic) (models.Clinic, error) {
					called = true
					return models.Clinic{}, tt.replaceErr
				},
			})
			rec := httptest.NewRecorder()

			h.ReplaceHandler(rec, newReq(http.MethodPut, tt.body, "1"))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
		})
	}
}
