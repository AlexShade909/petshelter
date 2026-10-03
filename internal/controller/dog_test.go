package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"petshelter/internal/models"
)

// ---------- mock ----------

type mockDogService struct {
	listNicknamesFn func() []string
	infoFn          func(nickname string) (models.Dog, error)
	deleteFn        func(nickname string) error
	createFn        func(dog models.Dog) error
	updateFn        func(dog models.Dog) (models.Dog, error)
	replaceFn       func(dog models.Dog) (models.Dog, error)

	calls int // сколько раз вызывали любой метод сервиса
}

func (m *mockDogService) ListNicknames() []string {
	m.calls++
	return m.listNicknamesFn()
}

func (m *mockDogService) Info(nickname string) (models.Dog, error) {
	m.calls++
	return m.infoFn(nickname)
}

func (m *mockDogService) Delete(nickname string) error {
	m.calls++
	return m.deleteFn(nickname)
}

func (m *mockDogService) Create(dog models.Dog) error {
	m.calls++
	return m.createFn(dog)
}

func (m *mockDogService) Update(dog models.Dog) (models.Dog, error) {
	m.calls++
	return m.updateFn(dog)
}

func (m *mockDogService) Replace(dog models.Dog) (models.Dog, error) {
	m.calls++
	return m.replaceFn(dog)
}

// ---------- helpers ----------

func mustJSON(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewReader(b)
}

func newReqWithName(method, name string, body *bytes.Reader) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, "/dogs/"+name, body)
	} else {
		req = httptest.NewRequest(method, "/dogs/"+name, nil)
	}
	req.SetPathValue("dogName", name)
	return req
}

func decodeBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(rec.Body).Decode(&v); err != nil {
		t.Fatalf("decode response body %q: %v", rec.Body.String(), err)
	}
	return v
}

// ---------- NicknamesHandler ----------

func TestDog_NicknamesHandler(t *testing.T) {
	tests := []struct {
		name  string
		names []string
	}{
		{"several nicknames", []string{"Rex", "Bobik", "Sharik"}},
		{"single nickname", []string{"Rex"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockDogService{
				listNicknamesFn: func() []string { return tt.names },
			}
			h := NewDog(svc)

			rec := httptest.NewRecorder()
			h.NicknamesHandler(rec, httptest.NewRequest(http.MethodGet, "/dogs", nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			got := decodeBody[[]string](t, rec)
			if !reflect.DeepEqual(got, tt.names) {
				t.Errorf("body = %v, want %v", got, tt.names)
			}
		})
	}
}

// ---------- InfoHandler ----------

func TestDog_InfoHandler(t *testing.T) {
	wantDog := models.Dog{Nickname: "Rex"}

	tests := []struct {
		name       string
		dogName    string
		svcDog     models.Dog
		svcErr     error
		wantStatus int
		wantErrMsg string
	}{
		{
			name:       "success",
			dogName:    "Rex",
			svcDog:     wantDog,
			wantStatus: http.StatusOK,
		},
		{
			name:       "service error",
			dogName:    "Ghost",
			svcErr:     errors.New("dog not found"),
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "dog not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotName string
			svc := &mockDogService{
				infoFn: func(n string) (models.Dog, error) {
					gotName = n
					return tt.svcDog, tt.svcErr
				},
			}
			h := NewDog(svc)

			rec := httptest.NewRecorder()
			h.InfoHandler(rec, newReqWithName(http.MethodGet, tt.dogName, nil))

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if gotName != tt.dogName {
				t.Errorf("service got name %q, want %q", gotName, tt.dogName)
			}
			if tt.wantErrMsg != "" {
				if got := strings.TrimSpace(rec.Body.String()); got != tt.wantErrMsg {
					t.Errorf("body = %q, want %q", got, tt.wantErrMsg)
				}
				return
			}
			if got := decodeBody[models.Dog](t, rec); !reflect.DeepEqual(got, tt.svcDog) {
				t.Errorf("body = %+v, want %+v", got, tt.svcDog)
			}
		})
	}
}

// ---------- DeleteHandler ----------

func TestDog_DeleteHandler(t *testing.T) {
	tests := []struct {
		name       string
		dogName    string
		svcErr     error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			dogName:    "Rex",
			wantStatus: http.StatusOK,
			wantBody:   "dog Deleted. His name: Rex",
		},
		{
			name:       "service error",
			dogName:    "Ghost",
			svcErr:     errors.New("dog not found"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "dog not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotName string
			svc := &mockDogService{
				deleteFn: func(n string) error {
					gotName = n
					return tt.svcErr
				},
			}
			h := NewDog(svc)

			rec := httptest.NewRecorder()
			h.DeleteHandler(rec, newReqWithName(http.MethodDelete, tt.dogName, nil))

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if gotName != tt.dogName {
				t.Errorf("service got name %q, want %q", gotName, tt.dogName)
			}
			if tt.svcErr != nil {
				if got := strings.TrimSpace(rec.Body.String()); got != tt.wantBody {
					t.Errorf("body = %q, want %q", got, tt.wantBody)
				}
				return
			}
			if got := decodeBody[string](t, rec); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}

// ---------- CreateHandler ----------

func TestDog_CreateHandler(t *testing.T) {
	validDog := models.Dog{Nickname: "Rex"}

	tests := []struct {
		name        string
		body        func(t *testing.T) *bytes.Reader
		svcErr      error
		wantStatus  int
		wantErrMsg  string
		wantSvcCall bool
	}{
		{
			name:       "invalid json",
			body:       func(t *testing.T) *bytes.Reader { return bytes.NewReader([]byte("{bad json")) },
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid request body",
		},
		{
			name:       "empty body",
			body:       func(t *testing.T) *bytes.Reader { return bytes.NewReader(nil) },
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid request body",
		},
		{
			name:       "empty nickname",
			body:       func(t *testing.T) *bytes.Reader { return mustJSON(t, models.Dog{}) },
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "empty nickname",
		},
		{
			name:        "service error",
			body:        func(t *testing.T) *bytes.Reader { return mustJSON(t, validDog) },
			svcErr:      errors.New("dog already exists"),
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "dog already exists",
			wantSvcCall: true,
		},
		{
			name:        "success",
			body:        func(t *testing.T) *bytes.Reader { return mustJSON(t, validDog) },
			wantStatus:  http.StatusOK,
			wantSvcCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotDog models.Dog
			svc := &mockDogService{
				createFn: func(d models.Dog) error {
					gotDog = d
					return tt.svcErr
				},
			}
			h := NewDog(svc)

			req := httptest.NewRequest(http.MethodPost, "/dogs", tt.body(t))
			rec := httptest.NewRecorder()
			h.CreateHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if called := svc.calls > 0; called != tt.wantSvcCall {
				t.Errorf("service called = %v, want %v", called, tt.wantSvcCall)
			}
			if tt.wantErrMsg != "" {
				if got := strings.TrimSpace(rec.Body.String()); got != tt.wantErrMsg {
					t.Errorf("body = %q, want %q", got, tt.wantErrMsg)
				}
				return
			}
			if !reflect.DeepEqual(gotDog, validDog) {
				t.Errorf("service got %+v, want %+v", gotDog, validDog)
			}
			if got := decodeBody[models.Dog](t, rec); !reflect.DeepEqual(got, validDog) {
				t.Errorf("body = %+v, want %+v", got, validDog)
			}
		})
	}
}

// ---------- UpdateHandler ----------

func TestDog_UpdateHandler(t *testing.T) {
	input := models.Dog{Nickname: "Rex"}
	updated := models.Dog{Nickname: "Rex"} // то, что вернёт сервис

	tests := []struct {
		name        string
		body        func(t *testing.T) *bytes.Reader
		svcErr      error
		wantStatus  int
		wantErrMsg  string
		wantSvcCall bool
	}{
		{
			name:       "invalid json",
			body:       func(t *testing.T) *bytes.Reader { return bytes.NewReader([]byte("{bad json")) },
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid request body",
		},
		{
			name:       "empty nickname",
			body:       func(t *testing.T) *bytes.Reader { return mustJSON(t, models.Dog{}) },
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "nickname is empty",
		},
		{
			name:        "service error",
			body:        func(t *testing.T) *bytes.Reader { return mustJSON(t, input) },
			svcErr:      errors.New("dog not found"),
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "dog not found",
			wantSvcCall: true,
		},
		{
			// NB: сейчас хендлер отвечает 201 — см. замечание про 200 ниже
			name:        "success",
			body:        func(t *testing.T) *bytes.Reader { return mustJSON(t, input) },
			wantStatus:  http.StatusOK,
			wantSvcCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotDog models.Dog
			svc := &mockDogService{
				updateFn: func(d models.Dog) (models.Dog, error) {
					gotDog = d
					if tt.svcErr != nil {
						return models.Dog{}, tt.svcErr
					}
					return updated, nil
				},
			}
			h := NewDog(svc)

			req := httptest.NewRequest(http.MethodPatch, "/dogs", tt.body(t))
			rec := httptest.NewRecorder()
			h.UpdateHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if called := svc.calls > 0; called != tt.wantSvcCall {
				t.Errorf("service called = %v, want %v", called, tt.wantSvcCall)
			}
			if tt.wantErrMsg != "" {
				if got := strings.TrimSpace(rec.Body.String()); got != tt.wantErrMsg {
					t.Errorf("body = %q, want %q", got, tt.wantErrMsg)
				}
				return
			}
			if !reflect.DeepEqual(gotDog, input) {
				t.Errorf("service got %+v, want %+v", gotDog, input)
			}
			if got := decodeBody[models.Dog](t, rec); !reflect.DeepEqual(got, updated) {
				t.Errorf("body = %+v, want %+v", got, updated)
			}
		})
	}
}

// ---------- ReplaceHandler ----------

func TestDog_ReplaceHandler(t *testing.T) {
	input := models.Dog{Nickname: "Rex"}
	replaced := models.Dog{Nickname: "Rex"}

	tests := []struct {
		name        string
		body        func(t *testing.T) *bytes.Reader
		svcErr      error
		wantStatus  int
		wantErrMsg  string
		wantSvcCall bool
	}{
		{
			name:       "invalid json",
			body:       func(t *testing.T) *bytes.Reader { return bytes.NewReader([]byte("{bad json")) },
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid request body",
		},
		{
			name:       "empty nickname",
			body:       func(t *testing.T) *bytes.Reader { return mustJSON(t, models.Dog{}) },
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "nickname is empty",
		},
		{
			name:        "service error",
			body:        func(t *testing.T) *bytes.Reader { return mustJSON(t, input) },
			svcErr:      errors.New("dog not found"),
			wantStatus:  http.StatusBadRequest,
			wantErrMsg:  "dog not found",
			wantSvcCall: true,
		},
		{
			name:        "success",
			body:        func(t *testing.T) *bytes.Reader { return mustJSON(t, input) },
			wantStatus:  http.StatusCreated,
			wantSvcCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotDog models.Dog
			svc := &mockDogService{
				replaceFn: func(d models.Dog) (models.Dog, error) {
					gotDog = d
					if tt.svcErr != nil {
						return models.Dog{}, tt.svcErr
					}
					return replaced, nil
				},
			}
			h := NewDog(svc)

			req := httptest.NewRequest(http.MethodPut, "/dogs", tt.body(t))
			rec := httptest.NewRecorder()
			h.ReplaceHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if called := svc.calls > 0; called != tt.wantSvcCall {
				t.Errorf("service called = %v, want %v", called, tt.wantSvcCall)
			}
			if tt.wantErrMsg != "" {
				if got := strings.TrimSpace(rec.Body.String()); got != tt.wantErrMsg {
					t.Errorf("body = %q, want %q", got, tt.wantErrMsg)
				}
				return
			}
			if !reflect.DeepEqual(gotDog, input) {
				t.Errorf("service got %+v, want %+v", gotDog, input)
			}
			if got := decodeBody[models.Dog](t, rec); !reflect.DeepEqual(got, replaced) {
				t.Errorf("body = %+v, want %+v", got, replaced)
			}
		})
	}
}
