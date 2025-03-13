package repository

import (
	"database/sql"

	"github.com/alextotalk/feline-intelligence/internal/domain"
	"github.com/alextotalk/feline-intelligence/internal/domain/model"
)

type TargetPgRepository struct {
	db *sql.DB
}

func NewTargetPgRepository(db *sql.DB) domain.TargetRepository {
	return &TargetPgRepository{db: db}
}

func (r *TargetPgRepository) AddToMission(t *model.Target) error {
	query := `
        INSERT INTO targets (mission_id, name, country, notes, complete)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
	return r.db.QueryRow(query, t.MissionID, t.Name, t.Country, t.Notes, t.Complete).
		Scan(&t.ID, &t.CreatedAt)
}

func (r *TargetPgRepository) Update(t *model.Target) error {
	query := `
        UPDATE targets
        SET name = $1, country = $2, notes = $3, complete = $4
        WHERE id = $5
    `
	_, err := r.db.Exec(query, t.Name, t.Country, t.Notes, t.Complete, t.ID)
	return err
}

func (r *TargetPgRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM targets WHERE id=$1`, id)
	return err
}

func (r *TargetPgRepository) GetByID(id int) (*model.Target, error) {
	query := `
        SELECT id, mission_id, name, country, notes, complete, created_at
        FROM targets
        WHERE id=$1
    `
	var t model.Target
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.MissionID, &t.Name, &t.Country, &t.Notes, &t.Complete, &t.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &t, nil
}
