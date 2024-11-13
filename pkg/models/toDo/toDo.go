package todo

type ToDo struct {
	Id          string
	Name        string
	Description string
	Completed   bool
}

func NewToDo(id string, name, description string, completed bool) ToDo {
	return ToDo{
		Id:          id,
		Name:        name,
		Description: description,
		Completed:   completed,
	}
}
