package cli

import (
	"bufio"
	"fmt"
	"strings"

	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	store "github.com/dan-harwood-bjss/toDoApp/pkg/stores/store"
)

func validateString(str string) bool {
	if str == "\n" || str == "" {
		return false
	}
	if str[0] == ' ' {
		return false
	}
	if str[len(str)-2] == ' ' {
		return false
	}
	return true
}

func removeNewlineFromString(str string) string {
	return strings.Split(str, "\n")[0]
}

func discardBuffer(r *bufio.Reader) {
	r.Discard(r.Buffered())
}

func PrintOptions() {
	fmt.Println("1 = list To Dos")
	fmt.Println("2 = create To Do")
	fmt.Println("3 = update To Do")
	fmt.Println("4 = delete To Do")
	fmt.Println("5 = exit")

}

func HandleChoice(choice int, store store.Store, stdin *bufio.Reader) bool {
	switch choice {
	case 1:
		HandleRead(store)
	case 2:
		HandleCreate(stdin, store)
	case 3:
		HandleUpdate(stdin, store)
	case 4:
		HandleDelete(stdin, store)
	case 5:
		return true
	default:
		fmt.Println("Option not understood, please choose a valid option.")
	}
	discardBuffer(stdin)
	return false
}
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
	if ok := validateString(name); !ok {
		fmt.Println("Could not Create Item as invalid name given.")
		return
	}
	fmt.Print("Description: ")
	description, _ := stdin.ReadString('\n')
	if ok := validateString(description); !ok {
		fmt.Println("Could not Create Item as invalid description given.")
		return
	}
	name = removeNewlineFromString(name)
	description = removeNewlineFromString(description)
	_, ok := store.Create(toDo.NewToDo(name, description, false))
	if !ok {
		fmt.Println("Failed to add to do.")
	}
}

func HandleUpdate(stdin *bufio.Reader, store store.Store) {
	fmt.Println("Enter the name of the To Do you wish to change status: (Case sensitive.)")
	fmt.Print("Name: ")
	name, _ := stdin.ReadString('\n')
	if ok := validateString(name); !ok {
		fmt.Println("Could not update item as invalid name given.")
		return
	}
	name = removeNewlineFromString(name)
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
	if ok := validateString(name); !ok {
		fmt.Println("Could not delete item as invalid name given.")
		return
	}
	name = removeNewlineFromString(name)
	if ok := store.Delete(name); !ok {
		fmt.Println("Failed to delete To Do. Check that it exists.")
	}
}
