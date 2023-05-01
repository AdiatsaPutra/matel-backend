package repositories

import (
	"motor/models"
	"motor/payloads"
)

type UserRepository interface {
	// View() ([]models.User, error)
	FindById(uint) (models.User, error)
	FindByUsername(string) (models.User, error)
	Create(payloads.CreateRequest) (bool, error)
	// Update(payloads.CreateRequest, uint) (bool, error)
	// Delete(uint) (bool, error)
}
