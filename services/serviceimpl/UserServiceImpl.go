package serviceimpl

import (
	"motor/models"
	"motor/repositories"
	"motor/services"
)

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func CreateUserService(userRepo repositories.UserRepository) services.UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// func (r *UserServiceImpl) ViewUser() ([]models.User, error) {
// 	get, err := r.userRepo.View()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return get, nil
// }

func (r *UserServiceImpl) FindUser(username string) (models.User, error) {
	get, err := r.userRepo.FindByUsername(username)

	if err != nil {
		return models.User{}, err
	}

	return get, nil
}

func (r *UserServiceImpl) FindById(userId uint) (models.User, error) {
	get, err := r.userRepo.FindById(userId)

	if err != nil {
		return models.User{}, err
	}

	return get, nil
}

// func (r *UserServiceImpl) UpdateUser(userReq payloads.CreateRequest, userId uint) (bool, error) {
// 	get, err := r.userRepo.Update(userReq, userId)

// 	if err != nil {
// 		return false, err
// 	}

// 	return get, nil
// }

// func (r *UserServiceImpl) DeleteUser(userId uint) (bool, error) {
// 	get, err := r.userRepo.Delete(userId)

// 	if err != nil {
// 		return false, err
// 	}

// 	return get, nil
// }
