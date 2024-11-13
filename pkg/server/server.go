package server

import (
	"log"
	"net/http"
	"text/template"

	todo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	jsonStore "github.com/dan-harwood-bjss/toDoApp/pkg/stores/json"
	"github.com/google/uuid"
)

type TodoPageData struct {
	PageTitle string
	Todos     map[string]todo.ToDo
}

func Create(store *jsonStore.JsonStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Default()
		if r.Method != http.MethodPost {
			logger.Printf("Received %s method to a POST endpoint.\n", r.Method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		logger.Printf("Received request: %v", r)
		r.ParseForm()
		name := r.Form.Get("name")
		if name == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		description := r.Form.Get("description")
		if description == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		completed := r.Form.Get("completed")
		isCompleted := false
		if completed != "" {
			isCompleted = true
		}
		item := todo.NewToDo(uuid.NewString(), name, description, isCompleted)
		if err := jsonStore.Create(store, item); err != nil {
			logger.Println("Got an error creating item:", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "", http.StatusSeeOther)
	}
}
func Read(store *jsonStore.JsonStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Default()
		if r.Method != http.MethodGet {
			logger.Printf("Received %s method to a GET endpoint.\n", r.Method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		logger.Printf("Received request: %v", r)
		items, err := jsonStore.Read(store)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles("./templates/todoList.html"))
		data := TodoPageData{
			PageTitle: "My TODO List",
			Todos:     items,
		}
		tmpl.Execute(w, data)
	}
}
