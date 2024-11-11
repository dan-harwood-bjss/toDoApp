package store

import (
	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
)

type Store interface {
	Create(toDo.ToDo) (toDo.ToDo, bool)
	Read(string) (toDo.ToDo, bool)
	Update(toDo.ToDo) bool
	Delete(string) bool
	ReadAll() (map[string]toDo.ToDo, bool)
}
