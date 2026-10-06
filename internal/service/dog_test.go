package service

import (
	"reflect"
	"testing"

	"petshelter/internal/models"
)

// resetDogs подменяет глобальную мапу и счётчик ID известными данными,
// а после теста восстанавливает исходные значения.
func resetDogs(t *testing.T, data map[int]models.Dog) {
	t.Helper()
	origData, origNext := dogsData, nexDogID
	dogsData = data
	nexDogID = len(data)
	t.Cleanup(func() {
		dogsData, nexDogID = origData, origNext
	})
}

func TestDog_ListNicknames(t *testing.T) {
	t.Run("sorted with ids", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{
			1: {ID: 1, Nickname: "Rex"},
			2: {ID: 2, Nickname: "Bim"},
		})

		got := NewDog().ListNicknames()
		want := []string{"Bim ID: 2", "Rex ID: 1"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{})

		got := NewDog().ListNicknames()

		if got == nil || len(got) != 0 {
			t.Errorf("got %#v, want empty non-nil slice", got)
		}
	})
}

func TestDog_Info(t *testing.T) {
	resetDogs(t, map[int]models.Dog{
		1: {ID: 1, Nickname: "Rex", Age: "3"},
	})

	t.Run("found", func(t *testing.T) {
		got, err := NewDog().Info(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 1 || got.Nickname != "Rex" || got.Age != "3" {
			t.Errorf("got %+v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := NewDog().Info(99)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDog_Delete(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{
			1: {ID: 1, Nickname: "Rex"},
			2: {ID: 2, Nickname: "Bim"},
		})

		err, nickname := NewDog().Delete(1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nickname != "Rex" {
			t.Errorf("nickname = %q, want Rex", nickname)
		}
		if _, ok := dogsData[1]; ok {
			t.Error("dog 1 still in storage")
		}
		if _, ok := dogsData[2]; !ok {
			t.Error("dog 2 was removed by mistake")
		}
	})

	t.Run("not found", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{1: {ID: 1, Nickname: "Rex"}})

		err, nickname := NewDog().Delete(99)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if nickname != "" {
			t.Errorf("nickname = %q, want empty", nickname)
		}
		if len(dogsData) != 1 {
			t.Errorf("storage size = %d, want 1", len(dogsData))
		}
	})
}

func TestDog_Create(t *testing.T) {
	t.Run("assigns id and stores", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{
			0: {ID: 0, Nickname: "Old"},
		})

		got, err := NewDog().Create(models.Dog{Nickname: "Rex", Age: "3"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Nickname != "Rex" || got.Age != "3" {
			t.Errorf("got %+v", got)
		}
		stored, ok := dogsData[got.ID]
		if !ok {
			t.Fatalf("dog with id %d not in storage", got.ID)
		}
		if stored.ID != got.ID {
			t.Errorf("stored.ID = %d, key = %d, must match", stored.ID, got.ID)
		}
		if dogsData[0].Nickname != "Old" {
			t.Errorf("existing dog overwritten: %+v", dogsData[0])
		}
	})

	t.Run("does not overwrite existing after delete", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{
			0: {ID: 0, Nickname: "A"},
			1: {ID: 1, Nickname: "B"},
			2: {ID: 2, Nickname: "C"},
		})
		svc := NewDog()
		svc.Delete(0)

		got, err := svc.Create(models.Dog{Nickname: "New"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dogsData[2].Nickname != "C" {
			t.Errorf("dog 2 overwritten: %+v", dogsData[2])
		}
		if got.ID == 2 {
			t.Errorf("new dog got occupied id %d", got.ID)
		}
		if len(dogsData) != 3 {
			t.Errorf("storage size = %d, want 3", len(dogsData))
		}
	})

	t.Run("ids are unique", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{})
		svc := NewDog()

		a, _ := svc.Create(models.Dog{Nickname: "A"})
		b, _ := svc.Create(models.Dog{Nickname: "B"})

		if a.ID == b.ID {
			t.Errorf("both dogs got id %d", a.ID)
		}
		if len(dogsData) != 2 {
			t.Errorf("storage size = %d, want 2", len(dogsData))
		}
	})

	t.Run("id from body is ignored", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{})

		got, _ := NewDog().Create(models.Dog{ID: 777, Nickname: "Rex"})

		if got.ID == 777 {
			t.Error("client-provided ID must be overwritten by the service")
		}
	})
}

func TestDog_Update(t *testing.T) {
	orig := models.Dog{
		ID:          1,
		Nickname:    "Rex",
		Age:         "3",
		WeightKg:    "20",
		CheckInDate: "2026-01-01",
		Shelter:     models.Shelter{Address: "old shelter"},
	}

	t.Run("partial update keeps other fields", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{1: orig})

		got, err := NewDog().Update(1, models.Dog{Age: "4"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Age != "4" {
			t.Errorf("Age = %q, want 4", got.Age)
		}
		if got.Nickname != "Rex" {
			t.Errorf("Nickname = %q, want Rex (must not be wiped)", got.Nickname)
		}
		if got.WeightKg != "20" || got.CheckInDate != "2026-01-01" {
			t.Errorf("other fields wiped: %+v", got)
		}
		if got.Shelter.Address != "old shelter" {
			t.Errorf("Shelter = %+v, want old shelter", got.Shelter)
		}
		if got.ID != 1 {
			t.Errorf("ID = %d, want 1", got.ID)
		}
		if !reflect.DeepEqual(dogsData[1], got) {
			t.Errorf("storage %+v differs from returned %+v", dogsData[1], got)
		}
	})

	t.Run("updates all provided fields", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{1: orig})

		got, err := NewDog().Update(1, models.Dog{
			Nickname:    "Max",
			Age:         "5",
			WeightKg:    "25",
			CheckInDate: "2026-02-02",
			Shelter:     models.Shelter{Address: "new shelter"},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Nickname != "Max" || got.Age != "5" || got.WeightKg != "25" ||
			got.CheckInDate != "2026-02-02" || got.Shelter.Address != "new shelter" {
			t.Errorf("got %+v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{1: orig})

		_, err := NewDog().Update(99, models.Dog{Age: "4"})

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDog_Replace(t *testing.T) {
	t.Run("replaces fully and forces id", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{
			1: {ID: 1, Nickname: "Rex", Age: "3", WeightKg: "20"},
		})

		got, err := NewDog().Replace(1, models.Dog{ID: 777, Nickname: "Max"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 1 {
			t.Errorf("ID = %d, want 1 (taken from path, not body)", got.ID)
		}
		if got.Nickname != "Max" || got.Age != "" || got.WeightKg != "" {
			t.Errorf("old fields must be cleared, got %+v", got)
		}
		if len(dogsData) != 1 {
			t.Errorf("storage size = %d, want 1", len(dogsData))
		}
		if dogsData[1].Nickname != "Max" {
			t.Errorf("storage not updated: %+v", dogsData[1])
		}
	})

	t.Run("not found", func(t *testing.T) {
		resetDogs(t, map[int]models.Dog{1: {ID: 1, Nickname: "Rex"}})

		_, err := NewDog().Replace(99, models.Dog{Nickname: "Max"})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if _, ok := dogsData[99]; ok {
			t.Error("dog 99 must not be created by Replace")
		}
	})
}
