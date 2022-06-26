package uow

import (
	"github.com/transactions/db"
	"github.com/transactions/entity"
	"gorm.io/gorm"
)

type todoUnitOfWork struct {
	todoDB db.TrxTodoDB
	subDB  db.TrxSubDB
	db     *gorm.DB
}

type TodoUnitOfWork interface {
	UpdateTodo(todo entity.Todo) error
	CreateTodo(todo entity.Todo) error
}

func NewTodo(todoDB db.TrxTodoDB, subDB db.TrxSubDB, db *gorm.DB) TodoUnitOfWork {
	return &todoUnitOfWork{
		todoDB: todoDB,
		subDB:  subDB,
		db:     db,
	}
}

func (t *todoUnitOfWork) UpdateTodo(todo entity.Todo) error {
	subs := todo.Subs
	todo.Subs = nil

	tx := t.db.Begin()
	err := t.todoDB.TrxUpdate(tx, todo)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = t.subDB.TrxRemoveByTodoID(tx, todo.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(subs) > 0 {
		err = t.subDB.TrxCreateList(tx, subs)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *todoUnitOfWork) CreateTodo(todo entity.Todo) error {
	var err error
	subs := todo.Subs
	todo.Subs = nil

	tx := t.db.Begin()
	err = t.todoDB.TrxCreate(tx, todo)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(subs) > 0 {
		err = t.subDB.TrxCreateList(tx, subs)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
