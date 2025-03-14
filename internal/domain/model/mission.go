package model

import "time"

// Mission описує місію для кота
type Mission struct {
	ID        int       `json:"id" example:"1"`
	CatID     *int      `json:"cat_id" example:"1"`
	Completed bool      `json:"completed" example:"false"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	Targets   []Target  `json:"targets"`
}
