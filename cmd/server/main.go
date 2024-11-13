package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dan-harwood-bjss/toDoApp/pkg/server"
	jsonStore "github.com/dan-harwood-bjss/toDoApp/pkg/stores/json"
)

func PrintHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func main() {
	file, err := os.OpenFile("db.json", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Received an error when opening file: %v\n", err)
	}
	store, err := jsonStore.NewJsonStore(file)
	if err != nil {
		log.Fatalf("Received an error when opening file: %v\n", err)
	}
	http.HandleFunc("/list", server.Read(store))
	http.HandleFunc("/create", server.Create(store))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
