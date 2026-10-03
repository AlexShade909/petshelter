Стуктура Собака
имя
возраст
вес
дата регистрации
шелтер пребывания
поликлиника регистрации

- М. Забрать собаку

Стуктура Шелтер
название 
адресс
номер
время работы
список питомцев

- М. Удалить собаку из базы

Структура Поликлиника
название
адресс
номер
время работы
список пациентов


- М. Удалить пациента


type Dog struct {
nickname    string
age         int
weightKg    float64
checkInDate string
shelter     *Shelter
policlinic  *Policlinic

type Shelter struct {
name        string
address     string
number      string
workingTime string
listPets    []Dog
}

type Policlinic struct {
name         string
address      string
number       string
workingTime  string
listPatients []
}


project/
├── main.go        — точка входа, только запуск программы
├── models.go       — структуры Dog, Shelter, Policlinic
├── data.go         — CreatePoliclinics, CreateShelters, CreateDogs (начальные данные)
├── service.go      — AddDog, RemoveDog, FindDog (бизнес-логика над данными)
├── input.go        — примитивы консольного ввода (readMenuChoice, readNonEmptyString и т.д.)
├── output.go       — функции вывода (printDogInfo, printShelterInfo, printDogList)
└── scenarios.go     — сценарии (scenarioAddDog, scenarioTakeDog) + runMenu