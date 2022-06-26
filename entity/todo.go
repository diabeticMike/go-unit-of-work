package entity

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	ID       uuid.UUID `gorm:"column:id"`
	Name     string    `gorm:"column:name"`
	Deadline time.Time `gorm:"column:deadline"`
	Subs     []Sub
}
