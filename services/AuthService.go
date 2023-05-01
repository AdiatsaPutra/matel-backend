package services

import (
	"motor/models"
	"motor/payloads"
)

type AuthService interface {
	DoLogin(payloads.LoginRequest) (models.User, error)
	DoRegister(payloads.CreateRequest) (bool, error)
}
