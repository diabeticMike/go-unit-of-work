package main

import (
	"github.com/transactions/db"
	"github.com/transactions/entity"
	"github.com/transactions/uow"
	"github.com/transactions/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	gormDB, err := NewConn()
	if err != nil {
		panic(err)
	}
	gormDB.AutoMigrate(&entity.Todo{}, &entity.Sub{})

	todoDB := db.NewTodoDB(gormDB)
	subDB := db.NewSubDB(gormDB)

	r := web.New(subDB, todoDB, uow.NewTodo(todoDB, subDB, gormDB))
	if err = http.ListenAndServe(":80", r); err != nil {
		panic(err)
	}
}

func NewConn() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("postgres://postgres:secret@localhost:5432/test?sslmode=disable"), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	return db, err
}
