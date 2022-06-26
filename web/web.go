package web

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/transactions/db"
	"github.com/transactions/entity"
	"github.com/transactions/uow"
	"gorm.io/gorm"
	"net/http"
)

func New(subDB db.SubDB, todoDB db.TodoDB, todoUOW uow.TodoUnitOfWork) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	c := newController(todoUOW, todoDB, subDB)
	router.HandleFunc("/todos", c.CreateTodo).Methods(http.MethodPost)
	router.HandleFunc("/todos/{todo_id}", c.UpdateTodo).Methods(http.MethodPut)
	router.HandleFunc("/todos/{todo_id}", c.GetTodo).Methods(http.MethodGet)
	router.HandleFunc("/subs/{sub_id}", c.UpdateSub).Methods(http.MethodPut)

	return router
}

type ctrl struct {
	todoUOW uow.TodoUnitOfWork
	todoDB  db.TodoDB
	subDB   db.SubDB
}

func newController(todoUOW uow.TodoUnitOfWork, todoDB db.TodoDB, subDB db.SubDB) *ctrl {
	return &ctrl{todoUOW: todoUOW, todoDB: todoDB, subDB: subDB}
}

func (c *ctrl) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo entity.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.ID = uuid.New()
	if len(todo.Subs) > 0 {
		for i := range todo.Subs {
			todo.Subs[i].ID = uuid.New()
			todo.Subs[i].TodoID = todo.ID
		}
	}

	err = c.todoUOW.CreateTodo(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newTodo, err := c.todoDB.GetByID(todo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (c *ctrl) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoUUID, err := uuid.Parse(mux.Vars(r)["todo_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var todo entity.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.ID = todoUUID
	if len(todo.Subs) > 0 {
		for i := range todo.Subs {
			todo.Subs[i].ID = uuid.New()
			todo.Subs[i].TodoID = todo.ID
		}
	}
	err = c.todoUOW.UpdateTodo(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedTodo, err := c.todoDB.GetByID(todo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&updatedTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (c *ctrl) GetTodo(w http.ResponseWriter, r *http.Request) {
	todoUUID, err := uuid.Parse(mux.Vars(r)["todo_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := c.todoDB.GetByID(todoUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (c *ctrl) UpdateSub(w http.ResponseWriter, r *http.Request) {
	subUUID, err := uuid.Parse(mux.Vars(r)["sub_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sub entity.Sub
	err = json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub.ID = subUUID
	err = c.subDB.Update(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedSub, err := c.subDB.GetByID(sub.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&updatedSub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
