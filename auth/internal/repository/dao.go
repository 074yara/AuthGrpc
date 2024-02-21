package repository

import "auth/internal/domain/entities"

type DAO interface {
	CreateUser(user *entities.UserDTO) error
	UpdateUser(user *entities.UserDTO) error
	DeleteUser(user *entities.UserDTO) error
	GetUser(user *entities.UserDTO) (*entities.UserDTO, error)
}
