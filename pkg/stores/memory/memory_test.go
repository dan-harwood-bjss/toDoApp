package memory

import (
	"reflect"
	"strconv"
	"sync"
	"testing"

	todo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
)

func TestCreate(t *testing.T) {
	t.Run("Can create an item", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewMemoryStore(nil)
		want := item

		got, ok := store.Create(item)

		if !ok {
			t.Fatal("Expected create to return ok but instead got not ok.")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but wanted %v", got, want)
		}
		if !reflect.DeepEqual(store.data[item.Name], item) {
			t.Errorf("Store does not contain %v instead contains %v", item, store.data[item.Name])
		}
		if len(store.data) != 1 {
			t.Errorf("Expected store to contain one item but got %d", len(store.data))
		}
	})
	t.Run("Returns not ok and does not create when item exists with key.", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewMemoryStore(map[string]todo.ToDo{item.Name: item})
		clashingItem := todo.NewToDo("Work", "Some other work", true)
		want := item

		got, ok := store.Create(clashingItem)

		if ok {
			t.Fatal("Expected create to return not ok but got ok.")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected create to return %v but got %v", want, got)
		}
		if len(store.data) != 1 {
			t.Errorf("Expected store to contain one item but got %d", len(store.data))
		}
		if !reflect.DeepEqual(store.data[item.Name], item) {
			t.Errorf("Expected store to contain %v but got %v", item, store.data[item.Name])
		}
	})
	t.Run("Returns not ok when not store not initialised.", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := &memoryStore{}

		_, ok := store.Create(item)

		if ok {
			t.Error("Read returned ok when Store not initialised.")
		}
	})
}
func TestRead(t *testing.T) {
	t.Run("Can read data from store.", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewMemoryStore(map[string]todo.ToDo{item.Name: item})
		want := item
		got, ok := store.Read(item.Name)

		if !ok {
			t.Fatal("Read returned not ok.")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but wanted %v", got, want)
		}
	})
	t.Run("Returns not ok when data is not present", func(t *testing.T) {
		store := NewMemoryStore(nil)

		_, ok := store.Read("Invalid")

		if ok {
			t.Error("Read returned ok as true, expected read to return ok as false.")
		}
	})
	t.Run("Returns not ok when not store not initialised.", func(t *testing.T) {
		store := &memoryStore{}

		_, ok := store.Read("Not initialised")

		if ok {
			t.Error("Read returned ok when Store not initialised.")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Can update store", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewMemoryStore(map[string]todo.ToDo{item.Name: item})
		updatedItem := todo.NewToDo("Work", "I have Changed", true)
		ok := store.Update(updatedItem)

		if !ok {
			t.Fatal("Expected update to return ok but returned not ok.")
		}
		if !reflect.DeepEqual(store.data[item.Name], updatedItem) {
			t.Errorf("Expected store to contain %v but instead it contains %v", updatedItem, store.data[item.Name])
		}
		if len(store.data) != 1 {
			t.Errorf("Expected store to contain 1 item but instead got %d", len(store.data))
		}
	})
	t.Run("Returns not ok when data not present.", func(t *testing.T) {
		store := NewMemoryStore(nil)
		item := todo.NewToDo("Work", "Do some work!", false)
		ok := store.Update(item)

		if ok {
			t.Error("Expected Update to return not ok but got ok.")
		}
	})
	t.Run("Returns not ok when store not initialised.", func(t *testing.T) {
		store := &memoryStore{}
		item := todo.NewToDo("Work", "Do some work!", false)
		ok := store.Update(item)

		if ok {
			t.Error("Expected store to return not ok when store not initialised but instead got ok.")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Deletes item successfully", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewMemoryStore(map[string]todo.ToDo{item.Name: item})
		ok := store.Delete(item.Name)

		if !ok {
			t.Fatal("Expected delete to return ok but got not ok.")
		}
	})
	t.Run("Returns not ok when item doesn't exist.", func(t *testing.T) {
		store := NewMemoryStore(nil)
		ok := store.Delete("Not exist")

		if ok {
			t.Error("Expected not ok but got ok.")
		}
	})
	t.Run("Returns not ok when store not initialised.", func(t *testing.T) {
		store := &memoryStore{}
		ok := store.Delete("Some key")

		if ok {
			t.Error("Expected store to return not ok when store not initialised but instead got ok.")
		}
	})
}

func TestReadAll(t *testing.T) {
	t.Run("Test happy paths", func(t *testing.T) {
		tests := []struct {
			Name string
			Want map[string]todo.ToDo
		}{
			{
				Name: "Returns empty map when there is no data.",
				Want: make(map[string]todo.ToDo),
			},
			{
				Name: "Returns all data when one item in map",
				Want: map[string]todo.ToDo{"Work": todo.NewToDo("Work", "Some Work", false)},
			},
			{
				Name: "Returns all data when more than one item.",
				Want: map[string]todo.ToDo{
					"Work":       todo.NewToDo("Work", "Some Work", false),
					"Other Work": todo.NewToDo("Other Work", "Some other Work", true),
				},
			},
		}
		for _, tc := range tests {
			store := NewMemoryStore(tc.Want)
			got, ok := store.ReadAll()
			if !ok {
				t.Fatal("Expected read all to return ok but got not ok.")
			}
			if !reflect.DeepEqual(got, tc.Want) {
				t.Errorf("got %v but wanted %v", got, tc.Want)
			}
		}
	})
	t.Run("Returns not ok when store is not initialised", func(t *testing.T) {
		store := &memoryStore{}
		_, ok := store.ReadAll()

		if ok {
			t.Error("Expected store to return not ok when store not initialised but instead got ok.")
		}
	})
}

func BenchmarkStore(b *testing.B) {
	store := NewMemoryStore(nil)
	items := []todo.ToDo{}
	for i := 0; i < b.N; i++ {
		items = append(items, todo.NewToDo(strconv.Itoa(i), "Some description", false))
	}
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			store.Create(items[i])
			wg.Done()
		}()
	}
	wg.Wait()
}
