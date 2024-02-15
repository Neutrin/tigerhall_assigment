package repositiories

import "github.com/nitin/tigerhall/core/internal/model"

type UserRepo interface {
	//Will return error if user exists
	Create(model.User) (string, error)
	UserExists(email string) bool
	User(email string) model.User
}
