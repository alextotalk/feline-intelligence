package model

import "time"

type Target struct {
	ID        int       `json:"id" example:"1"`
	MissionID int       `json:"mission_id" example:"1"`
	Name      string    `json:"name" example:"Target Alpha"`
	Country   string    `json:"country" example:"Meowland"`
	Notes     string    `json:"notes" example:"Highly guarded"`
	Complete  bool      `json:"complete" example:"false"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}
