package helper

import (
	"matel/models"
	"time"
)

func GetUserStatus(user models.User) int {

	currentTime := time.Now()
	if user.StartSubscription != "" && user.EndSubscription != "" {

		startSubscriptionTime, err := time.Parse("2006-01-02", user.StartSubscription)
		if err != nil {
			return 0
		}
		endSubscriptionTime, err := time.Parse("2006-01-02", user.EndSubscription)
		if err != nil {
			return 0
		}

		if currentTime.After(startSubscriptionTime) && currentTime.Before(endSubscriptionTime) {
			return 1
		} else if currentTime.After(endSubscriptionTime) {
			return 2
		}
	} else {
		if currentTime.Before(user.CreatedAt.AddDate(0, 0, 1)) {
			return 0
		} else if currentTime.After(user.CreatedAt.AddDate(0, 0, 1)) {
			return 2
		}
	}

	return 0
}
