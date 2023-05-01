package services

import (
	"motor/models"
)

type UserService interface {
	// ViewUser() ([]models.User, error)
	FindById(uint) (models.User, error)
	FindUser(string) (models.User, error)
	// UpdateUser(payloads.CreateRequest, uint) (bool, error)
	// DeleteUser(uint) (bool, error)
}
