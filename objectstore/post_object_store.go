package objectstore

type PostObjectStore interface {
	GetPost(location string) (string, error)
	DeletePost(location string) (string, error)
}
