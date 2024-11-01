package store

import "errors"

type StoreResponse struct {
	Item  Person
	Error error
}
type Person struct {
	Name string
	Age  int
}

type Store struct {
	data          map[string]Person
	CreateChannel chan Person
	ReadChannel   chan string
	UpdateChannel chan Person
	DeleteChannel chan Person
	OutChannel    chan StoreResponse
}

func (s *Store) Loop() {
	for {
		select {
		case key := <-s.ReadChannel:
			s.read(key)
		case item := <-s.UpdateChannel:
			s.update(item)
		}
	}
}

func (s *Store) Read(key string) StoreResponse {
	s.ReadChannel <- key
	return <-s.OutChannel
}
func (s *Store) read(key string) {
	value, ok := s.data[key]
	if !ok {
		s.OutChannel <- StoreResponse{Error: errors.New("not found")}
	}
	s.OutChannel <- StoreResponse{Item: value}
}

func (s *Store) update(item Person) {

}

//what if have a server that handles requests of different types, the server has access to a store and the store is looping, when the server wants to read it sends to
// server should handle validation
// store is literally for returning data and nothing else.
// a get channel and when it writes it sends to a done channel, the done channel
// could maybe have a data and error channel and the store can respond via the channel
