package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"sync"
	"testing"

	todo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	"github.com/google/uuid"
)

func TestNewJsonStore(t *testing.T) {
	t.Run("Create new store with data.", func(t *testing.T) {
		id := uuid.NewString()
		data := map[string]todo.ToDo{
			id: todo.NewToDo(id, "Work", "Do Work", false),
		}
		jsonData, _ := json.Marshal(data)
		buffer := &bytes.Buffer{}
		buffer.Write(jsonData)
		store, err := NewJsonStore(buffer)
		if err != nil {
			t.Fatalf("Expected no errors but got: %v", err)
		}
		if !reflect.DeepEqual(store.data, data) {
			t.Errorf("Wanted %v but got %v", data, store.data)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("Can create a ToDo item in the store", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		item := todo.NewToDo(uuid.NewString(), "Work", "Do work.", false)
		store, err := NewJsonStore(buffer)
		if err != nil {
			t.Fatalf("Got an error when creating store. Error: %v", err)
		}
		err = Create(store, item)
		if err != nil {
			t.Fatalf("Got an error when adding an item to the store. Error: %v", err)
		}
		got := store.data[item.Id]
		if !reflect.DeepEqual(got, item) {
			t.Errorf("Expected item to be in store but got %v", got)
		}
		if len(store.data) != 1 {
			t.Errorf("Expected store to contain one item but instead it contains %d", len(store.data))
		}
	})
}

func TestRead(t *testing.T) {
	id := uuid.NewString()
	data := map[string]todo.ToDo{
		id: todo.NewToDo(id, "Work", "Do Work", false),
	}
	jsonData, _ := json.Marshal(data)
	buffer := &bytes.Buffer{}
	buffer.Write(jsonData)
	store, err := NewJsonStore(buffer)
	if err != nil {
		t.Fatalf("Expected no errors but got: %v", err)
	}
	got, err := Read(store)
	if err != nil {
		t.Fatalf("Got an error when reading from store: %v", err)
	}
	if !reflect.DeepEqual(got, data) {
		t.Errorf("wanted %v but got %v", got, data)
	}

}
func BenchmarkJsonStore(b *testing.B) {
	buffer := &bytes.Buffer{}
	store, err := NewJsonStore(buffer)
	if err != nil {
		b.Fatalf("Received an error when creating store. Error: %v", err)
	}
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			Create(store, todo.NewToDo(uuid.NewString(), strconv.Itoa(i), "Something", false))
			wg.Done()
		}()
	}
	wg.Wait()
	if len(store.data) != b.N {
		b.Errorf("Expected store to have %d items but got %d items instead.", b.N, len(store.data))
	}
}

func FuzzCreate(f *testing.F) {
	testcases := []struct {
		Name        string
		Description string
		Completed   bool
	}{
		{
			Name:        "Work",
			Description: "Do work",
			Completed:   false,
		},
		{
			Name:        "Clean",
			Description: "Tidy up",
			Completed:   true,
		},
	}
	for _, tc := range testcases {
		f.Add(tc.Name, tc.Description, tc.Completed)
	}
	buffer := &bytes.Buffer{}
	store, err := NewJsonStore(buffer)
	if err != nil {
		f.Fatalf("Got an error when creating store. Error: %v", err)
	}
	f.Fuzz(func(t *testing.T, name string, description string, completed bool) {
		item := todo.NewToDo(uuid.NewString(), name, description, completed)
		err := Create(store, item)
		if err != nil {
			t.Errorf("Received an error: %v for inputs: %s, %s, %v", err, name, description, completed)
		}
	})
}
