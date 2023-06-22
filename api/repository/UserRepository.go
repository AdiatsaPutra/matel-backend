package repository

import (
	config "matel/configs"
	"matel/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserTotalInfo(c *gin.Context) (models.HomeUserInfo, error) {
	var user models.HomeUserInfo
	query := `
	SELECT
		SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) AS trial_members,
		SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS premium_members,
		SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) AS expired_members
	FROM
    	m_users;
	`
	result := config.InitDB().Raw(query).Scan(&user)

	if result.Error != nil {
		return models.HomeUserInfo{}, result.Error
	}

	return user, nil

}

func GetUserTotal(c *gin.Context) (uint, error) {
	var user []models.User
	result := config.InitDB().Find(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint(result.RowsAffected), nil

}

func GetUserByEmail(c *gin.Context, UserEmail string) (models.User, error) {
	var user = models.User{Email: UserEmail}
	result := config.InitDB().Where("email = ?", user.Email).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil

}

func GetUserByDeviceID(c *gin.Context, DeviceID string) (models.User, error) {
	var user = models.User{}
	result := config.InitDB().Where("device_id = ?", DeviceID).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil

}

func SetToken(c *gin.Context, user models.User) (bool, error) {
	result := config.InitDB().Model(&user).Where("id = ?", user.ID).Update("token", user.Token)

	if result.Error != nil {
		return true, result.Error
	}

	return true, nil

}

func CreateUser(c *gin.Context, user models.User) (models.User, error) {
	var newUser = models.User{}
	result := config.InitDB().Create(&user)

	newUser = user

	if result.Error != nil {
		return newUser, result.Error
	}

	return newUser, nil

}

func GetMember(c *gin.Context) ([]models.User, error) {
	var user = []models.User{}
	result := config.InitDB().Where("is_admin = 0").Find(&user)

	if len(user) == 0 {
		return user, nil
	}

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func Logout(c *gin.Context, UserID uint) error {
	var user models.User

	err := config.InitDB().Model(&user).Where("id = ?", UserID).Update("token", "").Error

	if err != nil {
		return err
	}

	e := config.InitDB().Model(&user).Where("id = ?", UserID).Update("device_id", "").Error

	if e != nil {
		return e
	}

	return nil

}

func ResetDeviceID(c *gin.Context, UserID uint, DeviceID string) error {
	var user models.User

	e := config.InitDB().Model(&user).Where("id = ?", UserID).Update("device_id", DeviceID).Error

	if e != nil {
		return e
	}

	return nil

}

func SetUser(c *gin.Context, userID uint, subscriptionMonths uint) error {
	var user models.User

	db := config.InitDB()

	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	now := time.Now()

	endSubscription := now.AddDate(0, int(subscriptionMonths), 0)

	startSubscriptionStr := now.Format("2006-01-02")
	endSubscriptionStr := endSubscription.Format("2006-01-02")

	if err := db.Model(&user).Updates(models.User{
		StartSubscription: startSubscriptionStr,
		EndSubscription:   endSubscriptionStr,
		Status:            1,
	}).Error; err != nil {
		return err
	}

	return nil
}

func AddSearchHistory(c *gin.Context, UserID uint, LeasingID uint) error {
	var user models.User

	err := config.InitDB().Model(&user).Where("id = ?", UserID).First(&user).Error
	if err != nil {
		return err
	}

	numbersViewed := user.NoPolHistory

	// Check if LeasingID already exists in nopol_history
	if !containsNumber(numbersViewed, int(LeasingID)) {
		if numbersViewed != "" {
			numbersViewed += ","
		}

		numbersViewed += strconv.Itoa(int(LeasingID))

		err = config.InitDB().Model(&user).Where("id = ?", UserID).Update("nopol_history", numbersViewed).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to check if a number exists in the given string
func containsNumber(numbersViewed string, number int) bool {
	numbers := strings.Split(numbersViewed, ",")
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err == nil && num == number {
			return true
		}
	}
	return false
}

func UserProfile(c *gin.Context, userId uint) (models.UserDetail, error) {
	var user models.UserDetail
	query := `
		SELECT u.*, p.name AS province_name, k.name AS kabupaten_name, kc.name AS kecamatan_name
		FROM m_users AS u
		JOIN m_province AS p ON u.province_id = p.id
		JOIN m_kabupaten AS k ON u.kabupaten_id = k.id
		JOIN m_kecamatan AS kc ON u.kecamatan_id = kc.id
		WHERE u.id = ?
	`

	result := config.InitDB().Raw(query, userId).Scan(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
