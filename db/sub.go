package db

import (
	"github.com/google/uuid"
	"github.com/transactions/entity"
	"gorm.io/gorm"
)

type subDB struct {
	db *gorm.DB
}

type SubDB interface {
	Update(sub entity.Sub) error
	GetByID(id uuid.UUID) (entity.Sub, error)
}

type TrxSubDB interface {
	TrxCreateList(tx *gorm.DB, subs []entity.Sub) error
	TrxRemoveByTodoID(tx *gorm.DB, id uuid.UUID) error
}

func NewSubDB(db *gorm.DB) *subDB {
	return &subDB{db: db}
}

func (s *subDB) Update(sub entity.Sub) error {
	return s.db.Model(&sub).
		Updates(map[string]interface{}{
			"name": sub.Name,
		}).Find(&sub, "id = ?", sub.ID).Error
}

func (s *subDB) GetByID(id uuid.UUID) (entity.Sub, error) {
	var sub entity.Sub
	err := s.db.First(&sub, "id = ?", id.String()).Error
	return sub, err
}

func (s *subDB) TrxCreateList(tx *gorm.DB, subs []entity.Sub) error {
	return tx.Create(&subs).Error
}

func (s *subDB) TrxRemoveByTodoID(tx *gorm.DB, id uuid.UUID) error {
	return tx.Where("todo_id = ?", id.String()).Delete(&entity.Sub{}).Error
}
