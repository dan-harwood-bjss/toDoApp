package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dan-harwood-bjss/toDoApp/pkg/server"
	jsonStore "github.com/dan-harwood-bjss/toDoApp/pkg/stores/json"
)

func saveToFile(store *jsonStore.JsonStore, timer *time.Ticker) {
	for range timer.C {
		file, _ := os.Create("db.json")
		jsonStore.WriteToFile(store, file)
		file.Close()
		fmt.Println("Stored data to file.")
	}
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
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", server.Read(store))
	http.HandleFunc("/create", server.Create(store))
	http.HandleFunc("/create-form", server.GetCreateForm(store))
	http.HandleFunc("/delete", server.Delete(store))
	http.HandleFunc("/update", server.Update(store))
	http.HandleFunc("/update-form", server.GetUpdateForm(store))
	timer := time.NewTicker(5 * time.Second)
	go saveToFile(store, timer)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
