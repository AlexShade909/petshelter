package repository

import (
	"context"
	"errors"
	"petshelter/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DogRepository struct {
	pool *pgxpool.Pool
}

func NewDogRepository(pool *pgxpool.Pool) *DogRepository {
	return &DogRepository{pool: pool}
}

func mapErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" { // нарушение внешнего ключа
		return models.ErrInvalidReference
	}
	return err
}
func (r *DogRepository) GetByID(id int) (models.Dog, error) {
	var d models.Dog
	ctx := context.Background() // заглушка
	err := r.pool.QueryRow(ctx, `
		SELECT d.id, d.nickname, d.age, d.weight_kg, d.check_in_date,
		       s.id, s.address, s.phone_number, s.working_time,
		       c.id, c.address, c.phone_number, c.working_time
		FROM dogs d
		JOIN shelters s ON s.id = d.shelter_id
		JOIN clinics  c ON c.id = d.clinic_id
		WHERE d.id = $1`, id,
	).Scan(
		&d.ID, &d.Nickname, &d.Age, &d.WeightKg, &d.CheckInDate,
		&d.Shelter.ID, &d.Shelter.Address, &d.Shelter.PhoneNumber, &d.Shelter.WorkingTime,
		&d.Сlinic.ID, &d.Сlinic.Address, &d.Сlinic.PhoneNumber, &d.Сlinic.WorkingTime,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Dog{}, errors.New("not found")
		}
		return models.Dog{}, err
	}
	return d, nil
}

func (r *DogRepository) ListNicknames() ([]string, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, `SELECT nickname FROM dogs ORDER BY nickname`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		names = append(names, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

func (r *DogRepository) Delete(id int) error {
	ctx := context.Background() // заглушка

	tag, err := r.pool.Exec(ctx, `DELETE FROM dogs WHERE id = $1`, id)
	if err != nil {
		return mapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *DogRepository) Create(dog models.Dog) (models.Dog, error) {
	ctx := context.Background() // заглушка
	var id int
	err := r.pool.QueryRow(ctx, `
		INSERT INTO dogs (nickname, age, weight_kg, check_in_date, shelter_id, clinic_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		dog.Nickname, dog.Age, dog.WeightKg, dog.CheckInDate,
		dog.Shelter.ID, dog.Сlinic.ID,
	).Scan(&id)
	if err != nil {
		return models.Dog{}, mapErr(err)
	}
	return r.GetByID(id) // вернуть собаку с заполненными Shelter и Clinic
}

// Replace (PUT): полная замена, обновляются все поля.
func (r *DogRepository) Replace(id int, dog models.Dog) (models.Dog, error) {
	ctx := context.Background() // заглушка
	tag, err := r.pool.Exec(ctx, `
		UPDATE dogs
		SET nickname = $2, age = $3, weight_kg = $4, check_in_date = $5,
		    shelter_id = $6, clinic_id = $7
		WHERE id = $1`,
		id, dog.Nickname, dog.Age, dog.WeightKg, dog.CheckInDate,
		dog.Shelter.ID, dog.Сlinic.ID,
	)
	if err != nil {
		return models.Dog{}, mapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return models.Dog{}, models.ErrNotFound
	}
	return r.GetByID(id)
}

// Update (PATCH): меняются только заполненные поля, нулевые значения игнорируются.
func (r *DogRepository) Update(id int, patch models.Dog) (models.Dog, error) {
	ctx := context.Background() // заглушка
	tag, err := r.pool.Exec(ctx, `
		UPDATE dogs
		SET nickname      = COALESCE(NULLIF($2::text, ''), nickname),
		    age           = COALESCE(NULLIF($3::int, 0),   age),
		    weight_kg     = COALESCE(NULLIF($4::int, 0),   weight_kg),
		    check_in_date = COALESCE(NULLIF($5::text, ''), check_in_date),
		    shelter_id    = COALESCE(NULLIF($6::int, 0),   shelter_id),
		    clinic_id     = COALESCE(NULLIF($7::int, 0),   clinic_id)
		WHERE id = $1`,
		id, patch.Nickname, patch.Age, patch.WeightKg, patch.CheckInDate,
		patch.Shelter.ID, patch.Сlinic.ID,
	)
	if err != nil {
		return models.Dog{}, mapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return models.Dog{}, models.ErrNotFound
	}
	return r.GetByID(id)
}
