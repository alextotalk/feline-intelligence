package repository

import (
	"database/sql"
	"fmt"

	"github.com/alextotalk/feline-intelligence/internal/domain"
	"github.com/alextotalk/feline-intelligence/internal/domain/model"
)

type CatPgRepository struct {
	db *sql.DB
}

func NewCatPgRepository(db *sql.DB) domain.CatRepository {
	return &CatPgRepository{db: db}
}

func (r *CatPgRepository) Create(cat *model.Cat) error {
	query := `
        INSERT INTO spy_cats (name, years_of_experience, breed, salary)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
	return r.db.QueryRow(query, cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary).
		Scan(&cat.ID, &cat.CreatedAt)
}

func (r *CatPgRepository) GetByID(id int) (*model.Cat, error) {
	cat := model.Cat{}
	query := `
        SELECT id, name, years_of_experience, breed, salary, created_at
        FROM spy_cats
        WHERE id = $1
    `
	err := r.db.QueryRow(query, id).
		Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary, &cat.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CatPgRepository) GetAll() ([]model.Cat, error) {
	rows, err := r.db.Query(`
        SELECT id, name, years_of_experience, breed, salary, created_at
        FROM spy_cats
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.Cat
	for rows.Next() {
		var c model.Cat
		if err := rows.Scan(&c.ID, &c.Name, &c.YearsOfExperience, &c.Breed, &c.Salary, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *CatPgRepository) Update(cat *model.Cat) error {
	query := `
        UPDATE spy_cats
        SET name = $1, years_of_experience = $2, breed = $3, salary = $4
        WHERE id = $5
    `
	_, err := r.db.Exec(query, cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary, cat.ID)
	return err
}

func (r *CatPgRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM spy_cats WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no cat found with id %d", id)
	}
	return nil
}
