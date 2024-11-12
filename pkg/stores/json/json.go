package json

import (
	"encoding/json"
	"io"

	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
)

type readAllOutcome struct {
	items map[string]toDo.ToDo
	ok    bool
}
type storeOutcome struct {
	item toDo.ToDo
	ok   bool
}

type jsonStore struct {
	data                 map[string]toDo.ToDo
	createChannel        chan toDo.ToDo
	readChannel          chan string
	updateChannel        chan toDo.ToDo
	deleteChannel        chan string
	readAllChannel       chan struct{}
	outChannel           chan storeOutcome
	readAllReturnChannel chan readAllOutcome
	initialised          bool
}

func readFile(file io.Reader) (map[string]toDo.ToDo, error) {
	data := make(map[string]toDo.ToDo)
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(fileData, &data)
	return data, nil
}

func NewJsonStore(file io.Reader) (*jsonStore, error) {
	data, err := readFile(file)
	if err != nil {
		return nil, err
	}
	store := &jsonStore{
		data:                 data,
		createChannel:        make(chan toDo.ToDo),
		readChannel:          make(chan string),
		updateChannel:        make(chan toDo.ToDo),
		deleteChannel:        make(chan string),
		outChannel:           make(chan storeOutcome),
		readAllChannel:       make(chan struct{}),
		readAllReturnChannel: make(chan readAllOutcome),
		initialised:          true,
	}
	go store.loop()
	return store, nil
}

func (s *jsonStore) loop() {
	for {
		select {
		case item := <-s.createChannel:
			s.create(item)
		case key := <-s.readChannel:
			s.read(key)
		case item := <-s.updateChannel:
			s.update(item)
		case key := <-s.deleteChannel:
			s.delete(key)
		case <-s.readAllChannel:
			s.readAll()
		}
	}
}

func (s *jsonStore) create(item toDo.ToDo) {
	value, ok := s.data[item.Name]
	if ok {
		s.outChannel <- storeOutcome{value, !ok}
		return
	}
	s.data[item.Name] = item
	s.outChannel <- storeOutcome{item, !ok}
}

func (s *jsonStore) read(key string) {
	value, ok := s.data[key]
	s.outChannel <- storeOutcome{value, ok}
}

func (s *jsonStore) update(item toDo.ToDo) {
	_, ok := s.data[item.Name]
	if !ok {
		s.outChannel <- storeOutcome{ok: ok}
		return
	}
	s.data[item.Name] = item
	s.outChannel <- storeOutcome{item, ok}
}

func (s *jsonStore) delete(key string) {
	_, ok := s.data[key]
	if !ok {
		s.outChannel <- storeOutcome{ok: ok}
		return
	}
	delete(s.data, key)
	s.outChannel <- storeOutcome{ok: ok}
}

func (s *jsonStore) readAll() {
	s.readAllReturnChannel <- readAllOutcome{items: s.data, ok: true}
}

func (s *jsonStore) Create(item toDo.ToDo) (value toDo.ToDo, ok bool) {
	if !s.initialised {
		return
	}
	s.createChannel <- item
	output := <-s.outChannel
	return output.item, output.ok
}

func (s *jsonStore) Read(key string) (value toDo.ToDo, ok bool) {
	if !s.initialised {
		return
	}
	s.readChannel <- key
	item := <-s.outChannel
	return item.item, item.ok
}

func (s *jsonStore) Update(item toDo.ToDo) (ok bool) {
	if !s.initialised {
		return
	}
	s.updateChannel <- item
	output := <-s.outChannel
	return output.ok
}

func (s *jsonStore) Delete(key string) (ok bool) {
	if !s.initialised {
		return
	}
	s.deleteChannel <- key
	output := <-s.outChannel
	return output.ok
}

func (s *jsonStore) ReadAll() (items map[string]toDo.ToDo, ok bool) {
	if !s.initialised {
		return
	}
	s.readAllChannel <- struct{}{}
	output := <-s.readAllReturnChannel
	return output.items, output.ok
}

func (s *jsonStore) Close(file io.Writer) error {
	jsonData, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	file.Write(jsonData)
}
