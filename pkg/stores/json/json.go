package json

import (
	"encoding/json"
	"io"

	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
)

const (
	CREATE     string = "CREATE"
	READ       string = "READ"
	UPDATE     string = "UPDATE"
	DELETE     string = "DELETE"
	FILE_WRITE string = "FILE_WRITE"
	GET_ITEM   string = "GET_ITEM"
)

type constError string

func (err constError) Error() string {
	return string(err)
}

const (
	ItemAlreadyExists   = constError("Can not create item as it already exists.")
	ItemDoesNotExist    = constError("Can not update item as it does not exist.")
	StoreNotInitialised = constError("Store has not been initialised.")
)

type storeInput struct {
	action string
	item   toDo.ToDo
	key    string
	file   io.ReadWriter
}

type storeResponse struct {
	items map[string]toDo.ToDo
	err   error
}

type JsonStore struct {
	data            map[string]toDo.ToDo
	inputChannel    chan storeInput
	responseChannel chan storeResponse
	initialised     bool
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

func NewJsonStore(file io.Reader) (*JsonStore, error) {
	data, err := readFile(file)
	if err != nil {
		return nil, err
	}
	store := &JsonStore{
		data:            data,
		inputChannel:    make(chan storeInput),
		responseChannel: make(chan storeResponse),
		initialised:     true,
	}
	go loop(store)
	return store, nil
}

func loop(s *JsonStore) {
	for input := range s.inputChannel {
		switch input.action {
		case CREATE:
			create(s, input.item)
		case READ:
			read(s)
		case UPDATE:
			update(s, input.item)
		case DELETE:
			deleteItem(s, input.key)
		case FILE_WRITE:
			writeToFile(s, input.file)
		case GET_ITEM:
			getItem(s, input.key)
		}
	}
}

func create(s *JsonStore, item toDo.ToDo) {
	_, ok := s.data[item.Id]
	if ok {
		s.responseChannel <- storeResponse{nil, ItemAlreadyExists}
		return
	}
	s.data[item.Id] = item
	s.responseChannel <- storeResponse{}
}

func read(s *JsonStore) {
	s.responseChannel <- storeResponse{s.data, nil}
}

func update(s *JsonStore, item toDo.ToDo) {
	_, ok := s.data[item.Id]
	if !ok {
		s.responseChannel <- storeResponse{nil, ItemDoesNotExist}
		return
	}
	s.data[item.Id] = item
	s.responseChannel <- storeResponse{}
}

func deleteItem(s *JsonStore, key string) {
	delete(s.data, key)
	s.responseChannel <- storeResponse{}
}

func writeToFile(s *JsonStore, file io.ReadWriter) {
	jsonData, err := json.Marshal(s.data)
	if err != nil {
		s.responseChannel <- storeResponse{nil, err}
		return

	}
	file.Write(jsonData)
	s.responseChannel <- storeResponse{nil, nil}
}

func getItem(s *JsonStore, key string) {
	item, ok := s.data[key]
	if !ok {
		s.responseChannel <- storeResponse{nil, ItemDoesNotExist}
		return
	}
	s.responseChannel <- storeResponse{map[string]toDo.ToDo{item.Id: item}, nil}
}

func Create(s *JsonStore, item toDo.ToDo) error {
	if !s.initialised {
		return StoreNotInitialised
	}
	input := storeInput{action: CREATE, item: item}
	s.inputChannel <- input
	output := <-s.responseChannel
	return output.err
}

func Read(s *JsonStore) (map[string]toDo.ToDo, error) {
	if !s.initialised {
		return nil, StoreNotInitialised
	}
	s.inputChannel <- storeInput{action: READ, item: toDo.ToDo{}}
	output := <-s.responseChannel
	return output.items, output.err
}

func Update(s *JsonStore, item toDo.ToDo) error {
	if !s.initialised {
		return StoreNotInitialised
	}
	s.inputChannel <- storeInput{action: UPDATE, item: item}
	output := <-s.responseChannel
	return output.err
}

func Delete(s *JsonStore, key string) error {
	if !s.initialised {
		return StoreNotInitialised
	}
	s.inputChannel <- storeInput{action: DELETE, key: key}
	<-s.responseChannel
	return nil
}

func WriteToFile(s *JsonStore, file io.ReadWriter) error {
	if !s.initialised {
		return StoreNotInitialised
	}
	s.inputChannel <- storeInput{action: FILE_WRITE, file: file}
	output := <-s.responseChannel
	return output.err
}

func GetItem(s *JsonStore, key string) (toDo.ToDo, error) {
	if !s.initialised {
		return toDo.ToDo{}, StoreNotInitialised
	}
	s.inputChannel <- storeInput{action: GET_ITEM, key: key}
	output := <-s.responseChannel
	return output.items[key], output.err
}
