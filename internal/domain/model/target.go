package model

import "time"

type Target struct {
	ID        int       `json:"id"`
	MissionID int       `json:"mission_id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Notes     string    `json:"notes"`
	Complete  bool      `json:"complete"`
	CreatedAt time.Time `json:"created_at"`
}
