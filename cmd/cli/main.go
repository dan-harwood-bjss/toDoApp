package main

import (
	"bufio"
	"fmt"
	"os"

	cli "github.com/dan-harwood-bjss/toDoApp/pkg/cli"
	jsonStore "github.com/dan-harwood-bjss/toDoApp/pkg/stores/json"
)

func main() {
	var choice int
	stdin := bufio.NewReader(os.Stdin)
	file, err := os.OpenFile("db.json", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Received an error when opening file: %v\n", err)
		return
	}
	store, _ := jsonStore.NewJsonStore(file)
	file.Close()
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
