package sbvision

// Session is a visit to the website
type Session struct {
	ID   int64  `json:"id"`
	Time int64  `json:"time"`
	IP   string `json:"ip"`
}

// SessionManager can create and validate sessions
type SessionManager interface {
	NewSession() (*Session, error)
	ValidateSession(*Session) error
}
