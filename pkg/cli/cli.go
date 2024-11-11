package cli

import (
	"bufio"
	"fmt"

	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	store "github.com/dan-harwood-bjss/toDoApp/pkg/stores/store"
)

func HandleRead(store store.Store) {
	items, ok := store.ReadAll()
	if !ok {
		fmt.Println("Not OK")
	} else if len(items) == 0 {
		fmt.Println("Nothing in list.")
	} else {
		fmt.Printf("To Do Items:\n%v\n", items)
	}
}

func HandleCreate(stdin *bufio.Reader, store store.Store) {
	fmt.Println("Please enter the information below:")
	fmt.Print("Name: ")
	name, _ := stdin.ReadString('\n')
	fmt.Print("Description: ")
	description, _ := stdin.ReadString('\n')
	_, ok := store.Create(toDo.NewToDo(name, description, false))
	if !ok {
		fmt.Println("Failed to add to do.")
	}
}

func HandleUpdate(stdin *bufio.Reader, store store.Store) {
	fmt.Println("Enter the name of the To Do you wish to change status: (Case sensitive.)")
	fmt.Print("Name: ")
	name, _ := stdin.ReadString('\n')
	item, ok := store.Read(name)
	if !ok {
		fmt.Println("Could not update item as it was not found.")
	} else {
		item.Completed = !item.Completed
		ok = store.Update(item)
		if !ok {
			fmt.Println("Failed to update to do.")
		}
	}
}

func HandleDelete(stdin *bufio.Reader, store store.Store) {
	fmt.Println("Enter the name of the To Do you wish to delete: (Case sensitive.)")
	fmt.Print("Name: ")
	name, _ := stdin.ReadString('\n')
	if ok := store.Delete(name); !ok {
		fmt.Println("Failed to delete To Do. Check that it exists.")
	}
}
