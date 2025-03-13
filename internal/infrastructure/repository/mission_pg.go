package repository

import (
	"database/sql"

	"github.com/alextotalk/feline-intelligence/internal/domain"
	"github.com/alextotalk/feline-intelligence/internal/domain/model"
)

type MissionPgRepository struct {
	db *sql.DB
}

func NewMissionPgRepository(db *sql.DB) domain.MissionRepository {
	return &MissionPgRepository{db: db}
}

func (r *MissionPgRepository) Create(m *model.Mission) error {
	query := `
        INSERT INTO missions (cat_id, completed)
        VALUES ($1, $2)
        RETURNING id, created_at
    `
	return r.db.QueryRow(query, m.CatID, m.Completed).
		Scan(&m.ID, &m.CreatedAt)
}

func (r *MissionPgRepository) GetByID(id int) (*model.Mission, error) {
	var ms model.Mission
	query := `
        SELECT id, cat_id, completed, created_at
        FROM missions
        WHERE id = $1
    `
	err := r.db.QueryRow(query, id).
		Scan(&ms.ID, &ms.CatID, &ms.Completed, &ms.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Витягуємо Targets для цієї місії
	tQuery := `
        SELECT id, mission_id, name, country, notes, complete, created_at
        FROM targets
        WHERE mission_id = $1
    `
	rows, err := r.db.Query(tQuery, ms.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []model.Target
	for rows.Next() {
		var t model.Target
		if err := rows.Scan(&t.ID, &t.MissionID, &t.Name, &t.Country, &t.Notes, &t.Complete, &t.CreatedAt); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	ms.Targets = targets

	return &ms, nil
}

func (r *MissionPgRepository) GetAll() ([]model.Mission, error) {
	rows, err := r.db.Query(`
        SELECT id, cat_id, completed, created_at
        FROM missions
        ORDER BY id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []model.Mission
	for rows.Next() {
		var ms model.Mission
		if err := rows.Scan(&ms.ID, &ms.CatID, &ms.Completed, &ms.CreatedAt); err != nil {
			return nil, err
		}

		// Витягуємо Targets
		tRows, err := r.db.Query(`
            SELECT id, mission_id, name, country, notes, complete, created_at
            FROM targets WHERE mission_id = $1
        `, ms.ID)
		if err != nil {
			return nil, err
		}
		var targets []model.Target
		for tRows.Next() {
			var t model.Target
			if err := tRows.Scan(&t.ID, &t.MissionID, &t.Name, &t.Country, &t.Notes, &t.Complete, &t.CreatedAt); err != nil {
				tRows.Close()
				return nil, err
			}
			targets = append(targets, t)
		}
		tRows.Close()

		ms.Targets = targets
		missions = append(missions, ms)
	}

	return missions, nil
}

func (r *MissionPgRepository) Update(m *model.Mission) error {
	query := `
        UPDATE missions
        SET cat_id = $1, completed = $2
        WHERE id = $3
    `
	_, err := r.db.Exec(query, m.CatID, m.Completed, m.ID)
	return err
}

func (r *MissionPgRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM missions WHERE id=$1`, id)
	return err
}

func (r *MissionPgRepository) AssignCat(missionID, catID int) error {
	// Записуємо в поле cat_id
	_, err := r.db.Exec(`UPDATE missions SET cat_id=$1 WHERE id=$2`, catID, missionID)
	return err
}
