package store

type Store interface {
	Create()
	Read()
	Update()
	Delete()
	ReadAll()
}
