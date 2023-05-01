package payloads

type UserRequest struct {
	UserName string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
	Token    string `json:"token"`
}

type MemberRequest struct {
	UserID uint
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDetail struct {
	Username  string `json:"username"`
	Authorize string `json:"authorized"`
}
