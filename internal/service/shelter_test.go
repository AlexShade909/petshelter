package service

import (
	"reflect"
	"sort"
	"testing"

	"petshelter/internal/models"
)

// resetShelters подменяет глобальную мапу приютов и счётчик ID известными данными,
// а после теста восстанавливает исходные значения.
// resetDogs находится в dog_test.go (тот же пакет).
func resetShelters(t *testing.T, data map[int]models.Shelter) {
	t.Helper()
	origData, origNext := sheltersData, nextShelterID
	sheltersData = data
	nextShelterID = len(data)
	t.Cleanup(func() {
		sheltersData, nextShelterID = origData, origNext
	})
}

var (
	testShelterA = models.Shelter{Address: "Шелтер 0", PhoneNumber: "+111", WorkingTime: "10:00-22:00"}
	testShelterB = models.Shelter{Address: "Шелтер 1", PhoneNumber: "+222", WorkingTime: "11:00-21:00"}
	testShelterC = models.Shelter{Address: "Шелтер 2", PhoneNumber: "+333", WorkingTime: "08:00-16:00"}
)

func TestShelter_FullInfo(t *testing.T) {
	resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB})

	got, err := NewShelter().FullInfo()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[int]models.Shelter{0: testShelterA, 1: testShelterB}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestShelter_Info(t *testing.T) {
	resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB})

	t.Run("found", func(t *testing.T) {
		got, err := NewShelter().Info(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != testShelterB {
			t.Errorf("got %+v, want %+v", got, testShelterB)
		}
	})

	// Этот тест падает на текущем коде: Info не проверяет наличие ключа
	// и возвращает пустой Shelter{} без ошибки.
	t.Run("not found", func(t *testing.T) {
		got, err := NewShelter().Info(99)
		if err == nil {
			t.Errorf("expected error, got nil (returned %+v)", got)
		}
	})

	t.Run("negative number", func(t *testing.T) {
		if _, err := NewShelter().Info(-1); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestShelter_ListDogs(t *testing.T) {
	dogs := map[int]models.Dog{
		0: {ID: 0, Nickname: "Rex", Shelter: testShelterA},
		1: {ID: 1, Nickname: "Bim", Shelter: testShelterA},
		2: {ID: 2, Nickname: "Max", Shelter: testShelterB},
	}

	t.Run("returns only dogs of the shelter", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB})
		resetDogs(t, dogs)

		got, err := NewShelter().ListDogs(0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		sort.Strings(got) // порядок обхода мапы случайный
		want := []string{"Bim", "Rex"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("shelter without dogs", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB, 2: testShelterC})
		resetDogs(t, dogs)

		got, err := NewShelter().ListDogs(2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("got %v, want no dogs", got)
		}
	})

	t.Run("invalid number", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB})
		resetDogs(t, dogs)

		for _, n := range []int{-1, 2, 99} {
			if _, err := NewShelter().ListDogs(n); err == nil {
				t.Errorf("number %d: expected error, got nil", n)
			}
		}
	})

	// Падает на текущем коде: проверка shelterNumber >= len(sheltersData)
	// ломается после удаления приюта, ключи перестают быть 0..len-1.
	t.Run("existing shelter works after another was deleted", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB, 2: testShelterC})
		resetDogs(t, map[int]models.Dog{
			0: {ID: 0, Nickname: "Rex", Shelter: testShelterC},
		})
		svc := NewShelter()
		if err := svc.Delete(1); err != nil { // размер мапы 2, а ключ 2 существует
			t.Fatalf("delete: %v", err)
		}

		got, err := svc.ListDogs(2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, []string{"Rex"}) {
			t.Errorf("got %v, want [Rex]", got)
		}
	})

	// Падает на текущем коде: номер удалённого приюта меньше len(), проверка проходит,
	// и для пустого Shelter{} находятся «собаки без приюта».
	t.Run("deleted shelter number returns error", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA, 2: testShelterC}) // ключа 1 нет
		resetDogs(t, map[int]models.Dog{
			0: {ID: 0, Nickname: "Stray"}, // без приюта
		})

		got, err := NewShelter().ListDogs(1)

		if err == nil {
			t.Errorf("expected error for missing shelter, got %v", got)
		}
	})
}

func TestShelter_Create(t *testing.T) {
	t.Run("stores shelter", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA})

		err := NewShelter().Create(testShelterB)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sheltersData) != 2 {
			t.Fatalf("storage size = %d, want 2", len(sheltersData))
		}
		if sheltersData[0] != testShelterA {
			t.Errorf("existing shelter changed: %+v", sheltersData[0])
		}
		if sheltersData[1] != testShelterB {
			t.Errorf("new shelter = %+v, want %+v", sheltersData[1], testShelterB)
		}
	})

	t.Run("does not overwrite existing after delete", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB, 2: testShelterC})
		svc := NewShelter()
		if err := svc.Delete(0); err != nil {
			t.Fatalf("delete: %v", err)
		}

		if err := svc.Create(models.Shelter{Address: "New"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if sheltersData[2] != testShelterC {
			t.Errorf("shelter 2 overwritten: %+v", sheltersData[2])
		}
		if len(sheltersData) != 3 {
			t.Errorf("storage size = %d, want 3", len(sheltersData))
		}
	})

	t.Run("each create gets its own id", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{})
		svc := NewShelter()

		_ = svc.Create(testShelterA)
		_ = svc.Create(testShelterB)

		if len(sheltersData) != 2 {
			t.Errorf("storage size = %d, want 2", len(sheltersData))
		}
	})
}

func TestShelter_Delete(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA, 1: testShelterB})

		err := NewShelter().Delete(0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := sheltersData[0]; ok {
			t.Error("shelter 0 still in storage")
		}
		if _, ok := sheltersData[1]; !ok {
			t.Error("shelter 1 was removed by mistake")
		}
	})

	// Падает на текущем коде: Delete всегда возвращает nil,
	// даже если приюта нет (у клиник и собак в этом случае ошибка).
	t.Run("not found", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{0: testShelterA})

		err := NewShelter().Delete(99)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if len(sheltersData) != 1 {
			t.Errorf("storage size = %d, want 1", len(sheltersData))
		}
	})
}

func TestShelter_Update(t *testing.T) {
	t.Run("partial update keeps other fields", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{1: testShelterA})

		got, err := NewShelter().Update(1, models.ShelterPatch{PhoneNumber: "+999"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := models.Shelter{
			Address:     testShelterA.Address,
			PhoneNumber: "+999",
			WorkingTime: testShelterA.WorkingTime,
		}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
		if sheltersData[1] != want {
			t.Errorf("storage %+v differs from returned %+v", sheltersData[1], want)
		}
	})

	t.Run("updates all fields", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{1: testShelterA})
		patch := models.ShelterPatch{Address: "A2", PhoneNumber: "P2", WorkingTime: "W2"}

		got, err := NewShelter().Update(1, patch)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := models.Shelter{Address: "A2", PhoneNumber: "P2", WorkingTime: "W2"}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("empty patch changes nothing", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{1: testShelterA})

		got, err := NewShelter().Update(1, models.ShelterPatch{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != testShelterA {
			t.Errorf("got %+v, want %+v", got, testShelterA)
		}
	})

	t.Run("not found", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{1: testShelterA})

		_, err := NewShelter().Update(99, models.ShelterPatch{Address: "X"})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if _, ok := sheltersData[99]; ok {
			t.Error("shelter 99 must not be created by Update")
		}
	})
}

func TestShelter_Replace(t *testing.T) {
	t.Run("replaces fully", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{1: testShelterA})
		replacement := models.Shelter{Address: "Only address"}

		got, err := NewShelter().Replace(1, replacement)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != replacement {
			t.Errorf("got %+v, want %+v", got, replacement)
		}
		if sheltersData[1] != replacement {
			t.Errorf("storage = %+v, want %+v (old fields must be cleared)", sheltersData[1], replacement)
		}
		if len(sheltersData) != 1 {
			t.Errorf("storage size = %d, want 1", len(sheltersData))
		}
	})

	// Падает на текущем коде: Replace не проверяет, что приют существует,
	// и молча создаёт новую запись с произвольным ключом.
	t.Run("not found", func(t *testing.T) {
		resetShelters(t, map[int]models.Shelter{1: testShelterA})

		_, err := NewShelter().Replace(99, testShelterB)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if _, ok := sheltersData[99]; ok {
			t.Error("shelter 99 must not be created by Replace")
		}
	})
}
