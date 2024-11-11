package main

import (
	"bufio"
	"fmt"
	"os"

	cli "github.com/dan-harwood-bjss/toDoApp/pkg/cli"
	memory "github.com/dan-harwood-bjss/toDoApp/pkg/stores/memory"
)

func main() {
	var choice int
	stdin := bufio.NewReader(os.Stdin)
	store := memory.NewMemoryStore(nil)
	for {
		choice = -1
		fmt.Println("To do list console app. Please select an option by typing a number for one of the below.")
		cli.PrintOptions()
		fmt.Fscanln(stdin, &choice)
		if exit := cli.HandleChoice(choice, store, stdin); exit {
			return
		}
	}
}
