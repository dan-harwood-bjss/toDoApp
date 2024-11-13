package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	jsonStore "github.com/dan-harwood-bjss/toDoApp/pkg/stores/json"
	"github.com/google/uuid"
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

func HandleChoice(choice int, store *jsonStore.JsonStore, stdin *bufio.Reader) bool {
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
		HandleExit(store)
		return true
	default:
		fmt.Println("Option not understood, please choose a valid option.")
	}
	discardBuffer(stdin)
	return false
}
func HandleRead(store *jsonStore.JsonStore) {
	items, err := jsonStore.Read(store)
	if err != nil {
		fmt.Println("Error reading store:", err)
	} else if len(items) == 0 {
		fmt.Println("Nothing in list.")
	} else {
		fmt.Printf("To Do Items:\n%v\n", items)
	}
}

func HandleCreate(stdin *bufio.Reader, store *jsonStore.JsonStore) {
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
	err := jsonStore.Create(store, toDo.NewToDo(uuid.NewString(), name, description, false))
	if err != nil {
		fmt.Println("Failed to add to do due to error:", err)
	}
}

func HandleUpdate(stdin *bufio.Reader, store *jsonStore.JsonStore) {
	fmt.Println("Enter the id of the To Do you wish to change status: ")
	fmt.Print("Id: ")
	id, _ := stdin.ReadString('\n')
	if ok := validateString(id); !ok {
		fmt.Println("Could not update item as invalid id given.")
		return
	}
	id = removeNewlineFromString(id)
	item, err := jsonStore.GetItem(store, id)
	if err != nil {
		fmt.Println("Error updating item:", err)
		return
	}
	item.Completed = !item.Completed
	err = jsonStore.Update(store, item)
	if err != nil {
		fmt.Println("Error updating item:", err)
		return
	}
}

func HandleDelete(stdin *bufio.Reader, store *jsonStore.JsonStore) {
	fmt.Println("Enter the id of the To Do you wish to delete: ")
	fmt.Print("id: ")
	id, _ := stdin.ReadString('\n')
	if ok := validateString(id); !ok {
		fmt.Println("Could not delete item as invalid id given.")
		return
	}
	id = removeNewlineFromString(id)
	if err := jsonStore.Delete(store, id); err != nil {
		fmt.Println("Failed to delete To Do due to error:", err)
	}
}

func HandleExit(store *jsonStore.JsonStore) {
	fmt.Println("Committing data to file.")
	file, err := os.Create("db.json")
	if err != nil {
		fmt.Println("Could not commit data to file, data will be wiped due to error:", err)
		return
	}
	defer file.Close()
	err = jsonStore.WriteToFile(store, file)
	if err != nil {
		fmt.Println("Could not commit data to file, data will be wiped due to error:", err)
		return
	}
}
