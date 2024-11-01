package todo

type ToDo struct {
	Name        string
	Description string
	Completed   bool
}

func NewToDo(name, description string, completed bool) ToDo {
	return ToDo{
		Name:        name,
		Description: description,
		Completed:   completed,
	}
}
