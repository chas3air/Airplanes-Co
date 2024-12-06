package models

import "github.com/google/uuid"

type Plane struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Capacity int       `json:"capacity"`
}
