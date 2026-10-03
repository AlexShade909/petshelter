package service

import (
	"reflect"
	"testing"

	"petshelter/internal/models"
)

// setDogs подменяет глобальное хранилище dogsData на время теста
// и возвращает исходное состояние в t.Cleanup.
// Из-за глобального состояния тесты НЕ должны запускаться через t.Parallel().
func setDogs(t *testing.T, dogs ...models.Dog) {
	t.Helper()
	orig := dogsData
	dogsData = make(map[string]models.Dog, len(dogs))
	for _, d := range dogs {
		dogsData[d.Nickname] = d
	}
	t.Cleanup(func() { dogsData = orig })
}

func TestNewDog(t *testing.T) {
	_ = NewDog() // конструктор не должен паниковать
}

func TestDog_ListNicknames(t *testing.T) {
	tests := []struct {
		name string
		dogs []models.Dog
		want []string
	}{
		{
			name: "empty storage returns empty non-nil slice",
			dogs: nil,
			want: []string{},
		},
		{
			name: "single dog",
			dogs: []models.Dog{{Nickname: "Rex"}},
			want: []string{"Rex"},
		},
		{
			name: "result is sorted",
			dogs: []models.Dog{{Nickname: "Sharik"}, {Nickname: "Bobik"}, {Nickname: "Rex"}},
			want: []string{"Bobik", "Rex", "Sharik"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setDogs(t, tt.dogs...)

			got := NewDog().ListNicknames()

			if got == nil {
				t.Fatal("got nil slice, want non-nil (JSON would encode it as null)")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDog_Info(t *testing.T) {
	rex := models.Dog{Nickname: "Rex"}

	tests := []struct {
		name    string
		nick    string
		want    models.Dog
		wantErr string
	}{
		{name: "existing dog", nick: "Rex", want: rex},
		{name: "unknown dog", nick: "Ghost", wantErr: "dog not found"},
		{name: "empty nickname", nick: "", wantErr: "dog not found"},
		{name: "case sensitive", nick: "rex", wantErr: "dog not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setDogs(t, rex)

			got, err := NewDog().Info(tt.nick)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				if !reflect.DeepEqual(got, models.Dog{}) {
					t.Errorf("got %+v on error, want zero value", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDog_Delete(t *testing.T) {
	t.Run("existing dog is removed", func(t *testing.T) {
		setDogs(t, models.Dog{Nickname: "Rex"}, models.Dog{Nickname: "Bobik"})

		if err := NewDog().Delete("Rex"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := dogsData["Rex"]; ok {
			t.Error("Rex still in storage after Delete")
		}
		if _, ok := dogsData["Bobik"]; !ok {
			t.Error("Bobik was removed, but should stay")
		}
	})

	t.Run("unknown dog returns error and storage is untouched", func(t *testing.T) {
		setDogs(t, models.Dog{Nickname: "Rex"})

		err := NewDog().Delete("Ghost")

		if err == nil || err.Error() != "dog not found" {
			t.Fatalf("err = %v, want %q", err, "dog not found")
		}
		if len(dogsData) != 1 {
			t.Errorf("storage size = %d, want 1", len(dogsData))
		}
	})

	t.Run("second delete of the same dog fails", func(t *testing.T) {
		setDogs(t, models.Dog{Nickname: "Rex"})
		d := NewDog()

		if err := d.Delete("Rex"); err != nil {
			t.Fatalf("first delete: %v", err)
		}
		if err := d.Delete("Rex"); err == nil {
			t.Error("second delete: want error, got nil")
		}
	})
}

func TestDog_Create(t *testing.T) {
	t.Run("new dog is stored", func(t *testing.T) {
		setDogs(t)
		rex := models.Dog{Nickname: "Rex"}

		if err := NewDog().Create(rex); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, ok := dogsData["Rex"]
		if !ok {
			t.Fatal("dog not in storage after Create")
		}
		if !reflect.DeepEqual(got, rex) {
			t.Errorf("stored %+v, want %+v", got, rex)
		}
	})

	t.Run("duplicate nickname returns error and keeps original", func(t *testing.T) {
		original := models.Dog{Nickname: "Rex"}
		setDogs(t, original)

		err := NewDog().Create(models.Dog{Nickname: "Rex"})

		if err == nil || err.Error() != "nickname occupied" {
			t.Fatalf("err = %v, want %q", err, "nickname occupied")
		}
		if !reflect.DeepEqual(dogsData["Rex"], original) {
			t.Errorf("original dog was overwritten: %+v", dogsData["Rex"])
		}
		if len(dogsData) != 1 {
			t.Errorf("storage size = %d, want 1", len(dogsData))
		}
	})

	t.Run("created dog is visible via Info and ListNicknames", func(t *testing.T) {
		setDogs(t)
		d := NewDog()

		if err := d.Create(models.Dog{Nickname: "Rex"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := d.Info("Rex"); err != nil {
			t.Errorf("Info after Create: %v", err)
		}
		if got := d.ListNicknames(); !reflect.DeepEqual(got, []string{"Rex"}) {
			t.Errorf("ListNicknames = %v, want [Rex]", got)
		}
	})
}

// Update и Replace сейчас имеют одинаковую реализацию (см. TODO в коде),
// поэтому тестируем их одним набором. Когда Update превратится в PATCH,
// для него нужен будет отдельный тест на частичное обновление полей.
func TestDog_UpdateAndReplace(t *testing.T) {
	methods := []struct {
		name string
		call func(d Dog, dog models.Dog) (models.Dog, error)
	}{
		{"Update", func(d Dog, dog models.Dog) (models.Dog, error) { return d.Update(dog) }},
		{"Replace", func(d Dog, dog models.Dog) (models.Dog, error) { return d.Replace(dog) }},
	}

	for _, m := range methods {
		t.Run(m.name+"/existing dog", func(t *testing.T) {
			setDogs(t, models.Dog{Nickname: "Rex"})
			// TODO: добавь сюда изменённое поле models.Dog (возраст, порода и т.п.),
			// иначе тест не отличит «заменил» от «ничего не сделал».
			newDog := models.Dog{Nickname: "Rex"}

			got, err := m.call(NewDog(), newDog)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, newDog) {
				t.Errorf("returned %+v, want %+v", got, newDog)
			}
			if !reflect.DeepEqual(dogsData["Rex"], newDog) {
				t.Errorf("stored %+v, want %+v", dogsData["Rex"], newDog)
			}
		})

		t.Run(m.name+"/unknown dog", func(t *testing.T) {
			setDogs(t, models.Dog{Nickname: "Rex"})

			got, err := m.call(NewDog(), models.Dog{Nickname: "Ghost"})

			if err == nil || err.Error() != "dog not found" {
				t.Fatalf("err = %v, want %q", err, "dog not found")
			}
			if !reflect.DeepEqual(got, models.Dog{}) {
				t.Errorf("got %+v on error, want zero value", got)
			}
			if _, ok := dogsData["Ghost"]; ok {
				t.Error("unknown dog must not be created by " + m.name)
			}
			if len(dogsData) != 1 {
				t.Errorf("storage size = %d, want 1", len(dogsData))
			}
		})
	}
}
