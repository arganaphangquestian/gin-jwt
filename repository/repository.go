package repository

import (
	"github.com/arganaphangquestian/gin-jwt/model"
)

type (
	// UserRepository interface
	UserRepository interface {
		Register(register model.InputUser) (*model.User, error)
		Login(login model.Login) (*model.User, error)
		Users() ([]*model.User, error)
	}
)
