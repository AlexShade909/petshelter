package service

import (
	"reflect"
	"sort"
	"testing"

	"petshelter/internal/models"
)

// resetClinics подменяет глобальную мапу клиник и счётчик ID известными данными,
// а после теста восстанавливает исходные значения.
// resetDogs находится в dog_test.go (тот же пакет).
func resetClinics(t *testing.T, data map[int]models.Clinic) {
	t.Helper()
	origData, origNext := clinicsData, nextClinicID
	clinicsData = data
	nextClinicID = len(data)
	t.Cleanup(func() {
		clinicsData, nextClinicID = origData, origNext
	})
}

var (
	testClinicA = models.Clinic{Address: "Мира 1", PhoneNumber: "+111", WorkingTime: "09:00-18:00"}
	testClinicB = models.Clinic{Address: "Ленина 133", PhoneNumber: "+222", WorkingTime: "10:00-20:00"}
	testClinicC = models.Clinic{Address: "Победы 5", PhoneNumber: "+333", WorkingTime: "08:00-16:00"}
)

func TestClinic_FullInfo(t *testing.T) {
	resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB})

	got, err := NewClinic().FullInfo()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[int]models.Clinic{0: testClinicA, 1: testClinicB}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestClinic_Info(t *testing.T) {
	resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB})

	t.Run("found", func(t *testing.T) {
		got, err := NewClinic().Info(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != testClinicB {
			t.Errorf("got %+v, want %+v", got, testClinicB)
		}
	})

	t.Run("not found", func(t *testing.T) {
		if _, err := NewClinic().Info(99); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestClinic_ListDogs(t *testing.T) {
	// ВНИМАНИЕ: поле в models.Dog называется с кириллической «С» (Сlinic).
	dogs := map[int]models.Dog{
		0: {ID: 0, Nickname: "Rex", Сlinic: testClinicA},
		1: {ID: 1, Nickname: "Bim", Сlinic: testClinicA},
		2: {ID: 2, Nickname: "Max", Сlinic: testClinicB},
	}

	t.Run("returns only dogs of the clinic", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB})
		resetDogs(t, dogs)

		got, err := NewClinic().ListDogs(0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		sort.Strings(got) // порядок обхода мапы случайный
		want := []string{"Bim", "Rex"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("clinic without dogs", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB, 2: testClinicC})
		resetDogs(t, dogs)

		got, err := NewClinic().ListDogs(2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("got %v, want no dogs", got)
		}
	})

	t.Run("invalid number", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB})
		resetDogs(t, dogs)

		for _, n := range []int{-1, 2, 99} {
			if _, err := NewClinic().ListDogs(n); err == nil {
				t.Errorf("number %d: expected error, got nil", n)
			}
		}
	})

	// Этот тест падает на текущем коде: проверка clinicNumber >= len(clinicsData)
	// ломается после удаления клиники, потому что ключи перестают быть 0..len-1.
	t.Run("existing clinic works after another was deleted", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB, 2: testClinicC})
		resetDogs(t, map[int]models.Dog{
			0: {ID: 0, Nickname: "Rex", Сlinic: testClinicC},
		})
		svc := NewClinic()
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

	// Этот тест тоже падает на текущем коде: номер удалённой клиники меньше len(),
	// проверка проходит, и для пустой Clinic{} находятся «собаки без клиники».
	t.Run("deleted clinic number returns error", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA, 2: testClinicC}) // ключа 1 нет
		resetDogs(t, map[int]models.Dog{
			0: {ID: 0, Nickname: "Stray"}, // без клиники
		})

		got, err := NewClinic().ListDogs(1)

		if err == nil {
			t.Errorf("expected error for missing clinic, got %v", got)
		}
	})
}

func TestClinic_Create(t *testing.T) {
	t.Run("stores clinic", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA})

		err := NewClinic().Create(testClinicB)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(clinicsData) != 2 {
			t.Fatalf("storage size = %d, want 2", len(clinicsData))
		}
		if clinicsData[0] != testClinicA {
			t.Errorf("existing clinic changed: %+v", clinicsData[0])
		}
		if clinicsData[1] != testClinicB {
			t.Errorf("new clinic = %+v, want %+v", clinicsData[1], testClinicB)
		}
	})

	t.Run("does not overwrite existing after delete", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB, 2: testClinicC})
		svc := NewClinic()
		if err := svc.Delete(0); err != nil {
			t.Fatalf("delete: %v", err)
		}

		if err := svc.Create(models.Clinic{Address: "New"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if clinicsData[2] != testClinicC {
			t.Errorf("clinic 2 overwritten: %+v", clinicsData[2])
		}
		if len(clinicsData) != 3 {
			t.Errorf("storage size = %d, want 3", len(clinicsData))
		}
	})

	t.Run("each create gets its own id", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{})
		svc := NewClinic()

		_ = svc.Create(testClinicA)
		_ = svc.Create(testClinicB)

		if len(clinicsData) != 2 {
			t.Errorf("storage size = %d, want 2", len(clinicsData))
		}
	})
}

func TestClinic_Delete(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA, 1: testClinicB})

		err := NewClinic().Delete(0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := clinicsData[0]; ok {
			t.Error("clinic 0 still in storage")
		}
		if _, ok := clinicsData[1]; !ok {
			t.Error("clinic 1 was removed by mistake")
		}
	})

	t.Run("not found", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{0: testClinicA})

		err := NewClinic().Delete(99)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if len(clinicsData) != 1 {
			t.Errorf("storage size = %d, want 1", len(clinicsData))
		}
	})
}

func TestClinic_Update(t *testing.T) {
	t.Run("partial update keeps other fields", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{1: testClinicA})

		got, err := NewClinic().Update(1, models.ClinicPatch{PhoneNumber: "+999"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := models.Clinic{
			Address:     testClinicA.Address,
			PhoneNumber: "+999",
			WorkingTime: testClinicA.WorkingTime,
		}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
		if clinicsData[1] != want {
			t.Errorf("storage %+v differs from returned %+v", clinicsData[1], want)
		}
	})

	t.Run("updates all fields", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{1: testClinicA})
		patch := models.ClinicPatch{Address: "A2", PhoneNumber: "P2", WorkingTime: "W2"}

		got, err := NewClinic().Update(1, patch)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := models.Clinic{Address: "A2", PhoneNumber: "P2", WorkingTime: "W2"}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("empty patch changes nothing", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{1: testClinicA})

		got, err := NewClinic().Update(1, models.ClinicPatch{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != testClinicA {
			t.Errorf("got %+v, want %+v", got, testClinicA)
		}
	})

	t.Run("not found", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{1: testClinicA})

		_, err := NewClinic().Update(99, models.ClinicPatch{Address: "X"})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if _, ok := clinicsData[99]; ok {
			t.Error("clinic 99 must not be created by Update")
		}
	})
}

func TestClinic_Replace(t *testing.T) {
	t.Run("replaces fully", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{1: testClinicA})
		replacement := models.Clinic{Address: "Only address"}

		got, err := NewClinic().Replace(1, replacement)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != replacement {
			t.Errorf("got %+v, want %+v", got, replacement)
		}
		if clinicsData[1] != replacement {
			t.Errorf("storage = %+v, want %+v (old fields must be cleared)", clinicsData[1], replacement)
		}
		if len(clinicsData) != 1 {
			t.Errorf("storage size = %d, want 1", len(clinicsData))
		}
	})

	// Этот тест падает на текущем коде: Replace не проверяет, что клиника существует,
	// и молча создаёт новую запись. У собак Replace возвращает ошибку.
	t.Run("not found", func(t *testing.T) {
		resetClinics(t, map[int]models.Clinic{1: testClinicA})

		_, err := NewClinic().Replace(99, testClinicB)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if _, ok := clinicsData[99]; ok {
			t.Error("clinic 99 must not be created by Replace")
		}
	})
}
