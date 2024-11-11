package main

import (
	"bufio"
	"fmt"
	"os"

	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
	memory "github.com/dan-harwood-bjss/toDoApp/pkg/stores/memory"
)

func discardBuffer(r *bufio.Reader) {
	r.Discard(r.Buffered())
}

func printOptions() {
	fmt.Println("1 = list To Dos")
	fmt.Println("2 = create To Do")
	fmt.Println("3 = update To Do")
	fmt.Println("4 = delete To Do")
	fmt.Println("5 = exit")

}

func main() {
	choice := 0
	stdin := bufio.NewReader(os.Stdin)
	store := memory.NewStore(nil)
	for {
		choice = 0
		fmt.Println("To do list console app. Please select an option by typing a number for one of the below.")
		printOptions()
		fmt.Fscanln(stdin, &choice)
		discardBuffer(stdin)
		switch choice {
		case 1:
			items, ok := memory.ReadAll(store)
			if !ok {
				fmt.Println("Not OK")
			} else if len(items) == 0 {
				fmt.Println("Nothing in list.")
			} else {
				fmt.Printf("To Do Items:\n%v\n", items)
			}
		case 2:
			name := ""
			description := ""
			fmt.Println("Please enter the information below:")
			fmt.Print("Name: ")
			name, _ = stdin.ReadString('\n')
			fmt.Print("Description: ")
			description, _ = stdin.ReadString('\n')
			_, ok := memory.Create(store, toDo.NewToDo(name, description, false))
			if !ok {
				fmt.Println("Failed to add to do.")
			}
		case 3:
			fmt.Println("Enter the name of the To Do you wish to change status: (Case sensitive.)")
			fmt.Print("Name: ")
			name, _ := stdin.ReadString('\n')
			item, ok := memory.Read(store, name)
			if !ok {
				fmt.Println("Could not update item as it was not found.")
			} else {
				item.Completed = !item.Completed
				ok = memory.Update(store, item)
				if !ok {
					fmt.Println("Failed to update to do.")
				}
			}
		case 4:
			fmt.Println("Enter the name of the To Do you wish to delete: (Case sensitive.)")
			fmt.Print("Name: ")
			name, _ := stdin.ReadString('\n')
			if ok := memory.Delete(store, name); !ok {
				fmt.Println("Failed to delete To Do. Check that it exists.")
			}
		case 5:
			return
		}
	}
}
