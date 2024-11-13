package json

import (
	"encoding/json"
	"io"

	toDo "github.com/dan-harwood-bjss/toDoApp/pkg/models/toDo"
)

const (
	CREATE    string = "CREATE"
	READ      string = "READ"
	UPDATE    string = "UPDATE"
	DELETE    string = "DELETE"
	FILEWRITE string = "FILEWRITE"
)

type constError string

func (err constError) Error() string {
	return string(err)
}

const (
	ItemAlreadyExists = constError("Can not create item as it already exists.")
	ItemDoesNotExist  = constError("Can not update item as it does not exist.")
)

type storeInput struct {
	action string
	item   toDo.ToDo
	file   io.ReadWriter
}

type storeResponse struct {
	items map[string]toDo.ToDo
	err   error
}

type jsonStore struct {
	data            map[string]toDo.ToDo
	inputChannel    chan storeInput
	responseChannel chan storeResponse
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
		data:            data,
		inputChannel:    make(chan storeInput),
		responseChannel: make(chan storeResponse),
	}
	go loop(store)
	return store, nil
}

func loop(s *jsonStore) {
	for input := range s.inputChannel {
		switch input.action {
		case CREATE:
			create(s, input.item)
		case READ:
			read(s)
		case UPDATE:
			update(s, input.item)
		case DELETE:
			deleteItem(s, input.item)
		case FILEWRITE:
			writeToFile(s, input.file)
		}
	}
}

func create(s *jsonStore, item toDo.ToDo) {
	_, ok := s.data[item.Id]
	if ok {
		s.responseChannel <- storeResponse{nil, ItemAlreadyExists}
		return
	}
	s.data[item.Id] = item
	s.responseChannel <- storeResponse{}
}

func read(s *jsonStore) {
	s.responseChannel <- storeResponse{s.data, nil}
}

func update(s *jsonStore, item toDo.ToDo) {
	_, ok := s.data[item.Id]
	if !ok {
		s.responseChannel <- storeResponse{nil, ItemDoesNotExist}
		return
	}
	s.data[item.Id] = item
	s.responseChannel <- storeResponse{}
}

func deleteItem(s *jsonStore, item toDo.ToDo) {
	delete(s.data, item.Id)
	s.responseChannel <- storeResponse{}
}

func writeToFile(s *jsonStore, file io.ReadWriter) {
	jsonData, err := json.Marshal(s.data)
	if err != nil {
		s.responseChannel <- storeResponse{nil, err}
		return

	}
	file.Write(jsonData)
}

func Create(s *jsonStore, item toDo.ToDo) error {
	input := storeInput{action: CREATE, item: item}
	s.inputChannel <- input
	output := <-s.responseChannel
	return output.err
}

func Read(s *jsonStore) map[string]toDo.ToDo {
	s.inputChannel <- storeInput{action: READ, item: toDo.ToDo{}}
	output := <-s.responseChannel
	return output.items
}

func Update(s *jsonStore, item toDo.ToDo) error {
	s.inputChannel <- storeInput{action: UPDATE, item: item}
	output := <-s.responseChannel
	return output.err
}

func Delete(s *jsonStore, item toDo.ToDo) {
	s.inputChannel <- storeInput{action: DELETE, item: item}
	<-s.responseChannel
}

func WriteToFile(s *jsonStore, file io.ReadWriter) error {
	s.inputChannel <- storeInput{action: FILEWRITE, file: file}
	output := <-s.responseChannel
	return output.err
}
