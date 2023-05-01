package repoimpl

import (
	"motor/models"
	"motor/payloads"
	"motor/repositories"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func CreateUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

// func (r *UserRepositoryImpl) View() ([]models.User, error) {
// 	row, err := r.DB.Raw("SELECT * FROM m_users").Rows()

// 	if err != nil {
// 		return nil, err
// 	}

// 	var users []models.User
// 	for row.Next() {
// 		var user models.User

// 		err := row.Scan(
// 			&user.ID,
// 			&user.FirstName,
// 			&user.LastName,
// 			&user.Username,
// 			&user.Password)

// 		if err != nil {
// 			return nil, err
// 		}

// 		users = append(users, user)
// 	}

// 	return users, nil
// }

func (r *UserRepositoryImpl) FindByUsername(username string) (models.User, error) {
	var user models.User

	r.DB.Model(user).Where("user_name=?", username).Scan(&user)

	return user, nil
}

func (r *UserRepositoryImpl) FindById(userId uint) (models.User, error) {
	var user models.User

	r.DB.Model(user).Where("user_id=?", userId).Scan(&user)

	return user, nil
}

func (r *UserRepositoryImpl) Create(userReq payloads.CreateRequest) (bool, error) {
	err := r.DB.Exec("INSERT INTO m_users (user_name, email, phone, password, device_id, token, is_admin) VALUES(?, ?, ?, ?, ?, ?, ?)",
		userReq.UserName,
		userReq.Email,
		userReq.Phone,
		userReq.Password,
		userReq.DeviceID,
		userReq.Token,
		0,
	).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

// func (r *UserRepositoryImpl) Update(userReq payloads.CreateRequest, userId uint) (bool, error) {
// 	err := r.DB.Exec("UPDATE m_users SET first_name=?, last_name=?, username=?, password=? WHERE user_id=?",
// 		userReq.FirstName,
// 		userReq.LastName,
// 		userReq.Username,
// 		userReq.Password,
// 		userId,
// 	).Error

// 	if err != nil {
// 		return false, err
// 	}

// 	return true, nil
// }

// func (r *UserRepositoryImpl) Delete(userId uint) (bool, error) {
// 	err := r.DB.Where("user_id=?", userId).Delete(&models.User{}).Error

// 	if err != nil {
// 		return false, err
// 	}

// 	return true, nil
// }
