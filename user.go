package sbvision

// User comes from the cognito user pool
type User struct {
	ID       int64
	Email    string `json:"email"`
	Username string `json:"username"`
}
