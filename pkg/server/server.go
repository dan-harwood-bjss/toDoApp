package server

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	todo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	jsonStore "github.com/dan-harwood-bjss/toDoApp/pkg/stores/json"
	"github.com/google/uuid"
)

type readResponse struct {
	items map[string]todo.ToDo
	err   error
}
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
func Read(store *jsonStore.JsonStore, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
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
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		responseChan := make(chan readResponse)
		go func() {
			items, err := jsonStore.Read(store)
			if err != nil {
				responseChan <- readResponse{nil, err}
				return
			}
			responseChan <- readResponse{items, nil}
		}()
		for {
			select {
			case <-ctx.Done():
				http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
				return
			case resp := <-responseChan:
				tmpl := template.Must(template.ParseFiles("./templates/todoList.html"))
				data := TodoPageData{
					PageTitle: "My TODO List",
					Todos:     resp.items,
				}
				tmpl.Execute(w, data)
				return
			}
		}
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
