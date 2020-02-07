package sbvision

// Session is a visit to the website
type Session struct {
	ID   int64 `json:"id"`
	Time int64 `json:"time"`
	IP   string
}

// SessionJWT is a base64encoded payload and signature of the session record
type SessionJWT string

// SessionManager can create and validate sessions
type SessionManager interface {
	CreateSession(*Session) (SessionJWT, error)
	ValidateSession(SessionJWT) (*Session, error)
}

// SessionTracker sets the session's auto-increment id
type SessionTracker interface {
	AddSession(*Session) error
}
