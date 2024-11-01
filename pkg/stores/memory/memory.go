package memory

import (
	toDo "example.com/toDoApp/pkg/models/toDo"
)

type storeOutcome struct {
	item toDo.ToDo
	ok   bool
}

type Store struct {
	data          map[string]toDo.ToDo
	createChannel chan toDo.ToDo
	readChannel   chan string
	updateChannel chan toDo.ToDo
	deleteChannel chan string
	outChannel    chan storeOutcome
	initialised   bool
}

func (s *Store) loop() {
	for {
		select {
		case item := <-s.createChannel:
			create(s, item)
		case key := <-s.readChannel:
			read(s, key)
		case item := <-s.updateChannel:
			s.update(item)
		case key := <-s.deleteChannel:
			s.delete(key)
		}
	}
}

func create(s *Store, item toDo.ToDo) {
	value, ok := s.data[item.Name]
	if ok {
		s.outChannel <- storeOutcome{value, !ok}
		return
	}
	s.data[item.Name] = item
	s.outChannel <- storeOutcome{item, !ok}
}

func read(s *Store, key string) {
	value, ok := s.data[key]
	s.outChannel <- storeOutcome{value, ok}
}

func (s *Store) update(item toDo.ToDo) {
	_, ok := s.data[item.Name]
	if !ok {
		s.outChannel <- storeOutcome{ok: ok}
		return
	}
	s.data[item.Name] = item
	s.outChannel <- storeOutcome{item, ok}
}

func (s *Store) delete(key string) {
	_, ok := s.data[key]
	if !ok {
		s.outChannel <- storeOutcome{ok: ok}
		return
	}
	delete(s.data, key)
	s.outChannel <- storeOutcome{ok: ok}
}

func NewStore(data map[string]toDo.ToDo) *Store {
	if data == nil {
		data = make(map[string]toDo.ToDo)
	}
	store := &Store{
		data:          data,
		createChannel: make(chan toDo.ToDo),
		readChannel:   make(chan string),
		updateChannel: make(chan toDo.ToDo),
		deleteChannel: make(chan string),
		outChannel:    make(chan storeOutcome),
		initialised:   true,
	}
	go store.loop()
	return store
}

func Create(s *Store, item toDo.ToDo) (value toDo.ToDo, ok bool) {
	if !s.initialised {
		return
	}
	s.createChannel <- item
	output := <-s.outChannel
	return output.item, output.ok
}

func Read(s *Store, key string) (value toDo.ToDo, ok bool) {
	if !s.initialised {
		return
	}
	s.readChannel <- key
	item := <-s.outChannel
	return item.item, item.ok
}

func Update(s *Store, item toDo.ToDo) (ok bool) {
	if !s.initialised {
		return
	}
	s.updateChannel <- item
	output := <-s.outChannel
	return output.ok
}

func Delete(s *Store, key string) (ok bool) {
	if !s.initialised {
		return
	}
	s.deleteChannel <- key
	output := <-s.outChannel
	return output.ok
}
