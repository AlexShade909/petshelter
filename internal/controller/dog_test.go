package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"petshelter/internal/models"
)

type mockDogService struct {
	listNicknamesFn func() []string
	infoFn          func(int) (models.Dog, error)
	deleteFn        func(int) (error, string)
	createFn        func(models.Dog) (error, models.Dog)
	updateFn        func(int, models.Dog) (models.Dog, error)
	replaceFn       func(int, models.Dog) (models.Dog, error)
}

func (m mockDogService) ListNicknames() []string                 { return m.listNicknamesFn() }
func (m mockDogService) Info(id int) (models.Dog, error)         { return m.infoFn(id) }
func (m mockDogService) Delete(id int) (error, string)           { return m.deleteFn(id) }
func (m mockDogService) Create(d models.Dog) (error, models.Dog) { return m.createFn(d) }
func (m mockDogService) Update(id int, d models.Dog) (models.Dog, error) {
	return m.updateFn(id, d)
}
func (m mockDogService) Replace(id int, d models.Dog) (models.Dog, error) {
	return m.replaceFn(id, d)
}

// path-параметр здесь называется dogID (в тестах клиник был NumberClinic)
func newDogReq(method, body, id string) *http.Request {
	req := httptest.NewRequest(method, "/dogs/"+id, strings.NewReader(body))
	req.SetPathValue("dogID", id)
	return req
}

func TestDogNicknamesHandler(t *testing.T) {
	h := NewDog(mockDogService{
		listNicknamesFn: func() []string { return []string{"Rex", "Bim"} },
	})
	rec := httptest.NewRecorder()

	h.NicknamesHandler(rec, newDogReq(http.MethodGet, "", ""))

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want %d", rec.Code, http.StatusOK)
	}
	var got []string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got) != 2 || got[0] != "Rex" || got[1] != "Bim" {
		t.Errorf("got %v, want [Rex Bim]", got)
	}
}

func TestDogInfoHandler(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		infoErr    error
		wantCode   int
		wantCalled bool
	}{
		{"ok", "5", nil, http.StatusOK, true},
		{"invalid id", "abc", nil, http.StatusBadRequest, false},
		{"service error", "5", errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called, gotID := false, 0
			h := NewDog(mockDogService{
				infoFn: func(id int) (models.Dog, error) {
					called, gotID = true, id
					return models.Dog{ID: id, Nickname: "Rex"}, tt.infoErr
				},
			})
			rec := httptest.NewRecorder()

			h.InfoHandler(rec, newDogReq(http.MethodGet, "", tt.id))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
			if tt.wantCode == http.StatusOK {
				if gotID != 5 {
					t.Errorf("service got id %d, want 5", gotID)
				}
				var resp models.Dog
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.ID != 5 || resp.Nickname != "Rex" {
					t.Errorf("response = %+v, want ID=5 Nickname=Rex", resp)
				}
			}
		})
	}
}

func TestDogDeleteHandler(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		deleteErr  error
		wantCode   int
		wantCalled bool
	}{
		{"ok", "5", nil, http.StatusOK, true},
		{"invalid id", "abc", nil, http.StatusBadRequest, false},
		{"service error", "5", errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			h := NewDog(mockDogService{
				deleteFn: func(int) (error, string) {
					called = true
					return tt.deleteErr, "Rex"
				},
			})
			rec := httptest.NewRecorder()

			h.DeleteHandler(rec, newDogReq(http.MethodDelete, "", tt.id))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
			if tt.wantCode == http.StatusOK {
				var msg string
				if err := json.NewDecoder(rec.Body).Decode(&msg); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if msg != "dog Deleted: Rex" {
					t.Errorf("message = %q, want %q", msg, "dog Deleted: Rex")
				}
			}
		})
	}
}

func TestDogCreateHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		createErr  error
		wantCode   int
		wantCalled bool
	}{
		{"ok", `{"Nickname":"Rex","Age":"3"}`, nil, http.StatusOK, true},
		{"invalid json", `{broken`, nil, http.StatusBadRequest, false},
		{"empty nickname", `{"Age":"3"}`, nil, http.StatusBadRequest, false},
		{"service error", `{"Nickname":"Rex"}`, errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			var got models.Dog
			h := NewDog(mockDogService{
				createFn: func(d models.Dog) (error, models.Dog) {
					called, got = true, d
					d.ID = 1 // сервис присваивает ID
					return tt.createErr, d
				},
			})
			rec := httptest.NewRecorder()

			h.CreateHandler(rec, newDogReq(http.MethodPost, tt.body, ""))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
			if tt.wantCode == http.StatusOK {
				if got.Nickname != "Rex" || got.Age != "3" {
					t.Errorf("service got %+v, want Nickname=Rex Age=3", got)
				}
				var resp models.Dog
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.ID != 1 || resp.Nickname != "Rex" {
					t.Errorf("response = %+v, want ID=1 Nickname=Rex", resp)
				}
			}
		})
	}
}

func TestDogUpdateHandler(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		body       string
		updateErr  error
		wantCode   int
		wantCalled bool
	}{
		{"ok", "5", `{"Nickname":"Max"}`, nil, http.StatusOK, true},
		{"invalid id", "abc", `{"Nickname":"Max"}`, nil, http.StatusBadRequest, false},
		{"invalid json", "5", `{broken`, nil, http.StatusBadRequest, false},
		{"empty nickname", "5", `{"Age":"4"}`, nil, http.StatusBadRequest, false},
		{"service error", "5", `{"Nickname":"Max"}`, errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called, gotID := false, 0
			var got models.Dog
			h := NewDog(mockDogService{
				updateFn: func(id int, d models.Dog) (models.Dog, error) {
					called, gotID, got = true, id, d
					return d, tt.updateErr
				},
			})
			rec := httptest.NewRecorder()

			h.UpdateHandler(rec, newDogReq(http.MethodPatch, tt.body, tt.id))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
			if tt.wantCode == http.StatusOK && (gotID != 5 || got.Nickname != "Max") {
				t.Errorf("service got id=%d dog=%+v, want id=5 Nickname=Max", gotID, got)
			}
		})
	}
}

func TestDogReplaceHandler(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		body       string
		replaceErr error
		wantCode   int
		wantCalled bool
	}{
		{"ok", "5", `{"Nickname":"Max"}`, nil, http.StatusCreated, true},
		{"invalid id", "abc", `{"Nickname":"Max"}`, nil, http.StatusBadRequest, false},
		{"invalid json", "5", `{broken`, nil, http.StatusBadRequest, false},
		{"empty nickname", "5", `{"Age":"4"}`, nil, http.StatusBadRequest, false},
		{"service error", "5", `{"Nickname":"Max"}`, errService, http.StatusBadRequest, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called, gotID := false, 0
			var got models.Dog
			h := NewDog(mockDogService{
				replaceFn: func(id int, d models.Dog) (models.Dog, error) {
					called, gotID, got = true, id, d
					return d, tt.replaceErr
				},
			})
			rec := httptest.NewRecorder()

			h.ReplaceHandler(rec, newDogReq(http.MethodPut, tt.body, tt.id))

			if rec.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", rec.Code, tt.wantCode)
			}
			if called != tt.wantCalled {
				t.Errorf("service called = %v, want %v", called, tt.wantCalled)
			}
			if tt.wantCode == http.StatusCreated && (gotID != 5 || got.Nickname != "Max") {
				t.Errorf("service got id=%d dog=%+v, want id=5 Nickname=Max", gotID, got)
			}
		})
	}
}
