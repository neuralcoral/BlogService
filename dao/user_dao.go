package dao

import "github.com/neuralcoral/BlogService/model"

type UserDao interface {
	GetUser(username string) (model.User, error)
}
