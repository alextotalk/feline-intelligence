package model

import "time"

type Cat struct {
	ID                int       `json:"id" example:"1"`
	Name              string    `json:"name" example:"Whiskers"`
	YearsOfExperience int       `json:"years_of_experience" example:"5"`
	Breed             string    `json:"breed" example:"Siamese"`
	Salary            float64   `json:"salary" example:"1000.0"`
	CreatedAt         time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}
