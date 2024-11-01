package store

import (
	"strconv"
	"testing"
)

func TestStore(t *testing.T) {
	store := &Store{
		data:          map[string]Person{"1": {"Dan", 30}},
		CreateChannel: make(chan Person),
		ReadChannel:   make(chan string),
		UpdateChannel: make(chan Person),
		DeleteChannel: make(chan Person),
	}
	go store.Loop()
	for i := 0; i < 500; i++ {
		go store.Read(strconv.Itoa(i))
	}
}
