package memory

import (
	"reflect"
	"strconv"
	"testing"

	todo "example.com/toDoApp/pkg/models/toDo"
)

func TestCreate(t *testing.T) {
	t.Run("Can create an item", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewStore(nil)
		want := item

		got, ok := Create(store, item)

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
		store := NewStore(map[string]todo.ToDo{item.Name: item})
		clashingItem := todo.NewToDo("Work", "Some other work", true)
		want := item

		got, ok := Create(store, clashingItem)

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
		store := &Store{}

		_, ok := Create(store, item)

		if ok {
			t.Error("Read returned ok when Store not initialised.")
		}
	})
}
func TestRead(t *testing.T) {
	t.Run("Can read data from store.", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewStore(map[string]todo.ToDo{item.Name: item})
		want := item
		got, ok := Read(store, item.Name)

		if !ok {
			t.Fatal("Read returned not ok.")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but wanted %v", got, want)
		}
	})
	t.Run("Returns not ok when data is not present", func(t *testing.T) {
		store := NewStore(nil)

		_, ok := Read(store, "Invalid")

		if ok {
			t.Error("Read returned ok as true, expected read to return ok as false.")
		}
	})
	t.Run("Returns not ok when not store not initialised.", func(t *testing.T) {
		store := &Store{}

		_, ok := Read(store, "Not initialised")

		if ok {
			t.Error("Read returned ok when Store not initialised.")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Can update store", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewStore(map[string]todo.ToDo{item.Name: item})
		updatedItem := todo.NewToDo("Work", "I have Changed", true)
		ok := Update(store, updatedItem)

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
		store := NewStore(nil)
		item := todo.NewToDo("Work", "Do some work!", false)
		ok := Update(store, item)

		if ok {
			t.Error("Expected Update to return not ok but got ok.")
		}
	})
	t.Run("Returns not ok when store not initialised.", func(t *testing.T) {
		store := &Store{}
		item := todo.NewToDo("Work", "Do some work!", false)
		ok := Update(store, item)

		if ok {
			t.Error("Expected store to return not ok when store not initialised but instead got ok.")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Deletes item successfully", func(t *testing.T) {
		item := todo.NewToDo("Work", "Do some work!", false)
		store := NewStore(map[string]todo.ToDo{item.Name: item})
		ok := Delete(store, item.Name)

		if !ok {
			t.Fatal("Expected delete to return ok but got not ok.")
		}
	})
	t.Run("Returns not ok when item doesn't exist.", func(t *testing.T) {
		store := NewStore(nil)
		ok := Delete(store, "Not exist")

		if ok {
			t.Error("Expected not ok but got ok.")
		}
	})
	t.Run("Returns not ok when store not initialised.", func(t *testing.T) {
		store := &Store{}
		ok := Delete(store, "Some key")

		if ok {
			t.Error("Expected store to return not ok when store not initialised but instead got ok.")
		}
	})
}
func BenchmarkStore(b *testing.B) {
	store := NewStore(nil)
	items := []todo.ToDo{}
	for i := 0; i < b.N; i++ {
		items = append(items, todo.NewToDo(strconv.Itoa(i), "Some description", false))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go Create(store, items[i])
		go Read(store, items[i].Name)
		go Update(store, items[i])
		go Delete(store, items[i].Name)
	}
}
