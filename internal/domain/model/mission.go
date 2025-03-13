package model

import "time"

// Mission описує місію для кота
type Mission struct {
	ID        int       `json:"id"`
	CatID     *int      `json:"cat_id"` // Якщо кіт ще не призначений, може бути nil
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`

	Targets []Target `json:"targets"`
}
