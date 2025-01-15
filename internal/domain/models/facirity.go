package models

import "github.com/google/uuid"

type Facility struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
