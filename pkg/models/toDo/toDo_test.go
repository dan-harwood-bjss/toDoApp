package todo

import "fmt"

func ExampleNewToDo() {
	fmt.Println(NewToDo("1", "Work", "Do Work", false))
	// Output: {1 Work Do Work false}
}
