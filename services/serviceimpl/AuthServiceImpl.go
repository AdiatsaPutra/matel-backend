package serviceimpl

import (
	"motor/models"
	"motor/payloads"
	"motor/repositories"
	"motor/services"
)

type AuthServiceImpl struct {
	userRepo repositories.UserRepository
}

func CreateAuthService(userRepo repositories.UserRepository) services.AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
	}
}

func (s *AuthServiceImpl) DoLogin(userReq payloads.LoginRequest) (models.User, error) {
	get, err := s.userRepo.FindByUsername(userReq.Username)

	if err != nil {
		return models.User{}, err
	}

	return get, nil
}

func (s *AuthServiceImpl) DoRegister(userReq payloads.CreateRequest) (bool, error) {
	get, err := s.userRepo.Create(userReq)

	if err != nil {
		return false, err
	}

	return get, nil
}
