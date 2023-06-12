package models

type HomeUserInfo struct {
	TrialMembers   uint `json:"trial_members"`
	ExpiredMembers uint `json:"expired_members"`
	PremiumMembers uint `json:"premium_members"`
}
