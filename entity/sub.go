package entity

import (
	"github.com/google/uuid"
)

type Sub struct {
	ID     uuid.UUID `gorm:"id"`
	Name   string    `gorm:"name"`
	TodoID uuid.UUID `gorm:"todo_id"`
}
