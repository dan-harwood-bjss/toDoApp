package main

import (
	"bufio"
	"fmt"
	"os"

	cli "github.com/dan-harwood-bjss/toDoApp/pkg/cli"
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
	store := memory.NewMemoryStore(nil)
	for {
		choice = 0
		fmt.Println("To do list console app. Please select an option by typing a number for one of the below.")
		printOptions()
		fmt.Fscanln(stdin, &choice)
		discardBuffer(stdin)
		switch choice {
		case 1:
			cli.HandleRead(store)
		case 2:
			cli.HandleCreate(stdin, store)
		case 3:
			cli.HandleUpdate(stdin, store)
		case 4:
			cli.HandleDelete(stdin, store)
		case 5:
			return
		}
	}
}
