package server

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	todo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	jsonStore "github.com/dan-harwood-bjss/toDoApp/pkg/stores/json"
	"github.com/google/uuid"
)

type TodoPageData struct {
	PageTitle string
	Todos     map[string]todo.ToDo
}

func GetCreateForm(store *jsonStore.JsonStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/createForm.html"))
		tmpl.Execute(w, nil)
	}
}

func Create(store *jsonStore.JsonStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Default()
		if r.Method != http.MethodPost {
			logger.Printf("Received %s method to a POST endpoint.\n", r.Method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		logger.Printf("Received request of method: %v\n", r.Method)
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
		if filepath.Clean(r.URL.Path) != "/" {
			http.NotFound(w, r)
			return
		}
		logger := log.Default()
		if r.Method != http.MethodGet {
			logger.Printf("Received %s method to a GET endpoint.\n", r.Method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		logger.Printf("Received request method: %v\n", r.Method)
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

func GetUpdateForm(store *jsonStore.JsonStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		item, _ := jsonStore.GetItem(store, id)
		tmpl := template.Must(template.ParseFiles("./templates/updateForm.html"))
		tmpl.Execute(w, item)
	}
}

func Update(store *jsonStore.JsonStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		r.ParseForm()
		name := r.Form.Get("name")
		description := r.Form.Get("description")
		completed := r.Form.Get("completed")
		isCompleted := false
		if completed != "" {
			isCompleted = true
		}
		item := todo.NewToDo(id, name, description, isCompleted)
		jsonStore.Update(store, item)
		http.Redirect(w, r, "", http.StatusSeeOther)
	}
}

func Delete(store *jsonStore.JsonStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		jsonStore.Delete(store, id)
		http.Redirect(w, r, "", http.StatusSeeOther)
	}
}
