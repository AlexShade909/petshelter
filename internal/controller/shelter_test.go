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
// Mock
// =========================

type mockShelterService struct {
	fullInfoFunc func() (map[int]models.Shelter, error)
	infoFunc     func(int) (models.Shelter, error)
	listDogsFunc func(int) ([]string, error)
	createFunc   func(models.Shelter) error
	deleteFunc   func(int) error
	updateFunc   func(int, models.ShelterPatch) (models.Shelter, error)
	replaceFunc  func(int, models.Shelter) (models.Shelter, error)
}

func (m *mockShelterService) FullInfo() (map[int]models.Shelter, error) {
	return m.fullInfoFunc()
}

func (m *mockShelterService) Info(number int) (models.Shelter, error) {
	return m.infoFunc(number)
}

func (m *mockShelterService) ListDogs(number int) ([]string, error) {
	return m.listDogsFunc(number)
}

func (m *mockShelterService) Create(shelter models.Shelter) error {
	return m.createFunc(shelter)
}

func (m *mockShelterService) Delete(number int) error {
	return m.deleteFunc(number)
}

func (m *mockShelterService) Update(
	number int,
	patch models.ShelterPatch,
) (models.Shelter, error) {
	return m.updateFunc(number, patch)
}

func (m *mockShelterService) Replace(
	number int,
	shelter models.Shelter,
) (models.Shelter, error) {
	return m.replaceFunc(number, shelter)
}

// =========================
// FullInfoHandler
// =========================

func TestShelter_FullInfoHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mockShelterService{
			fullInfoFunc: func() (map[int]models.Shelter, error) {
				return map[int]models.Shelter{
					1: {},
					2: {},
				}, nil
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/shelters",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.FullInfoHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockShelterService{
			fullInfoFunc: func() (map[int]models.Shelter, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/shelters",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.FullInfoHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}

		if !strings.Contains(rec.Body.String(), "database error") {
			t.Fatalf(
				"expected error message, got %q",
				rec.Body.String(),
			)
		}
	})
}

// =========================
// InfoHandler
// =========================

func TestShelter_InfoHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var receivedNumber int

		service := &mockShelterService{
			infoFunc: func(number int) (models.Shelter, error) {
				receivedNumber = number
				return models.Shelter{}, nil
			},

			listDogsFunc: func(number int) ([]string, error) {
				return []string{
					"Bobik",
					"Sharik",
				}, nil
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/shelters/42",
			nil,
		)

		req.SetPathValue("NumberShelter", "42")

		rec := httptest.NewRecorder()

		handler.InfoHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}

		if receivedNumber != 42 {
			t.Fatalf(
				"expected number 42, got %d",
				receivedNumber,
			)
		}
	})

	t.Run("info error", func(t *testing.T) {
		service := &mockShelterService{
			infoFunc: func(number int) (models.Shelter, error) {
				return models.Shelter{}, errors.New("shelter not found")
			},

			listDogsFunc: func(number int) ([]string, error) {
				t.Fatal("ListDogs should not be called")
				return nil, nil
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/shelters/42",
			nil,
		)

		req.SetPathValue("NumberShelter", "42")

		rec := httptest.NewRecorder()

		handler.InfoHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("list dogs error", func(t *testing.T) {
		service := &mockShelterService{
			infoFunc: func(number int) (models.Shelter, error) {
				return models.Shelter{}, nil
			},

			listDogsFunc: func(number int) ([]string, error) {
				return nil, errors.New("cannot get dogs")
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/shelters/42",
			nil,
		)

		req.SetPathValue("NumberShelter", "42")

		rec := httptest.NewRecorder()

		handler.InfoHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})
}

// =========================
// CreateHandler
// =========================

func TestShelter_CreateHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var called bool

		service := &mockShelterService{
			createFunc: func(s models.Shelter) error {
				called = true
				return nil
			},
		}

		handler := NewShelter(service)

		body := `{
"NumberShelter": 1
}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/shelters",
			strings.NewReader(body),
		)

		rec := httptest.NewRecorder()

		handler.CreateHandler(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusCreated,
				rec.Code,
			)
		}

		if !called {
			t.Fatal("expected Create to be called")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		service := &mockShelterService{
			createFunc: func(s models.Shelter) error {
				t.Fatal("Create should not be called")
				return nil
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodPost,
			"/shelters",
			strings.NewReader(`invalid json`),
		)

		rec := httptest.NewRecorder()

		handler.CreateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockShelterService{
			createFunc: func(s models.Shelter) error {
				return errors.New("shelter already exists")
			},
		}

		handler := NewShelter(service)

		body := `{
"NumberShelter": 1
}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/shelters",
			strings.NewReader(body),
		)

		rec := httptest.NewRecorder()

		handler.CreateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})
}

// =========================
// DeleteHandler
// =========================

func TestShelter_DeleteHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var receivedNumber int

		service := &mockShelterService{
			deleteFunc: func(number int) error {
				receivedNumber = number
				return nil
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/shelters/15",
			nil,
		)

		req.SetPathValue("NumberShelter", "15")

		rec := httptest.NewRecorder()

		handler.DeleteHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}

		if receivedNumber != 15 {
			t.Fatalf(
				"expected number 15, got %d",
				receivedNumber,
			)
		}

		if !strings.Contains(rec.Body.String(), "15") {
			t.Fatal("expected response to contain shelter number")
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockShelterService{
			deleteFunc: func(number int) error {
				return errors.New("shelter not found")
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/shelters/15",
			nil,
		)

		req.SetPathValue("NumberShelter", "15")

		rec := httptest.NewRecorder()

		handler.DeleteHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})
}

// =========================
// UpdateHandler
// =========================

func TestShelter_UpdateHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var receivedNumber int

		service := &mockShelterService{
			updateFunc: func(
				number int,
				patch models.ShelterPatch,
			) (models.Shelter, error) {

				receivedNumber = number

				return models.Shelter{}, nil
			},
		}

		handler := NewShelter(service)

		body := `{
"Address": "New address"
}`

		req := httptest.NewRequest(
			http.MethodPatch,
			"/shelters/10",
			strings.NewReader(body),
		)

		req.SetPathValue("NumberShelter", "10")

		rec := httptest.NewRecorder()

		handler.UpdateHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}

		if receivedNumber != 10 {
			t.Fatalf(
				"expected number 10, got %d",
				receivedNumber,
			)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		service := &mockShelterService{
			updateFunc: func(
				number int,
				patch models.ShelterPatch,
			) (models.Shelter, error) {

				t.Fatal("Update should not be called")

				return models.Shelter{}, nil
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/shelters/10",
			strings.NewReader(`invalid`),
		)

		req.SetPathValue("NumberShelter", "10")

		rec := httptest.NewRecorder()

		handler.UpdateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		service := &mockShelterService{
			updateFunc: func(
				number int,
				patch models.ShelterPatch,
			) (models.Shelter, error) {

				t.Fatal("Update should not be called")

				return models.Shelter{}, nil
			},
		}

		handler := NewShelter(service)

		body := `{
"SomeUnknownField": "test"
}`

		req := httptest.NewRequest(
			http.MethodPatch,
			"/shelters/10",
			strings.NewReader(body),
		)

		req.SetPathValue("NumberShelter", "10")

		rec := httptest.NewRecorder()

		handler.UpdateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockShelterService{
			updateFunc: func(
				number int,
				patch models.ShelterPatch,
			) (models.Shelter, error) {

				return models.Shelter{}, errors.New("update failed")
			},
		}

		handler := NewShelter(service)

		body := `{
"Address": "New address"
}`

		req := httptest.NewRequest(
			http.MethodPatch,
			"/shelters/10",
			strings.NewReader(body),
		)

		req.SetPathValue("NumberShelter", "10")

		rec := httptest.NewRecorder()

		handler.UpdateHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})
}

// =========================
// ReplaceHandler
// =========================

func TestShelter_ReplaceHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var receivedNumber int

		service := &mockShelterService{
			replaceFunc: func(
				number int,
				shelter models.Shelter,
			) (models.Shelter, error) {

				receivedNumber = number

				return shelter, nil
			},
		}

		handler := NewShelter(service)

		body := `{
"NumberShelter": 20
}`

		req := httptest.NewRequest(
			http.MethodPut,
			"/shelters/20",
			strings.NewReader(body),
		)

		req.SetPathValue("NumberShelter", "20")

		rec := httptest.NewRecorder()

		handler.ReplaceHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}

		if receivedNumber != 20 {
			t.Fatalf(
				"expected number 20, got %d",
				receivedNumber,
			)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		service := &mockShelterService{
			replaceFunc: func(
				number int,
				shelter models.Shelter,
			) (models.Shelter, error) {

				t.Fatal("Replace should not be called")

				return models.Shelter{}, nil
			},
		}

		handler := NewShelter(service)

		req := httptest.NewRequest(
			http.MethodPut,
			"/shelters/20",
			strings.NewReader(`invalid json`),
		)

		req.SetPathValue("NumberShelter", "20")

		rec := httptest.NewRecorder()

		handler.ReplaceHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockShelterService{
			replaceFunc: func(
				number int,
				shelter models.Shelter,
			) (models.Shelter, error) {

				return models.Shelter{}, errors.New("replace failed")
			},
		}

		handler := NewShelter(service)

		body := `{
"NumberShelter": 20
}`

		req := httptest.NewRequest(
			http.MethodPut,
			"/shelters/20",
			strings.NewReader(body),
		)

		req.SetPathValue("NumberShelter", "20")

		rec := httptest.NewRecorder()

		handler.ReplaceHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})
}
