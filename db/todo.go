package db

import (
	"github.com/google/uuid"
	"github.com/transactions/entity"
	"gorm.io/gorm"
)

type todoDB struct {
	db *gorm.DB
}

type TodoDB interface {
	GetByID(id uuid.UUID) (entity.Todo, error)
}
type TrxTodoDB interface {
	TrxUpdate(tx *gorm.DB, todo entity.Todo) error
	TrxCreate(tx *gorm.DB, todo entity.Todo) error
}

func NewTodoDB(db *gorm.DB) *todoDB {
	return &todoDB{db: db}
}

func (t *todoDB) GetByID(id uuid.UUID) (entity.Todo, error) {
	var todo entity.Todo
	err := t.db.Table("todos").Preload("Subs").First(&todo, "id = ?", id.String()).Error
	return todo, err
}

func (t *todoDB) TrxUpdate(tx *gorm.DB, todo entity.Todo) error {
	return tx.Model(&todo).
		Updates(map[string]interface{}{
			"name":     todo.Name,
			"deadline": todo.Deadline,
		}).Find(&todo, "id = ?", todo.ID.String()).Error
}

func (t *todoDB) TrxCreate(tx *gorm.DB, todo entity.Todo) error {
	return tx.Create(&todo).Error
}
