package session

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/kevinwylder/sbvision"
)

type sessionManager struct {
	key *rsa.PrivateKey
}

// NewRSASessionManager generates an rsa key to sign session tokens with
func NewRSASessionManager(db sbvision.SessionTracker) (sbvision.SessionManager, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	return &sessionManager{key: key}, nil
}

// SignSession uses the RSA key to generate a JWT
func (s *sessionManager) CreateSession(session *sbvision.Session) (sbvision.SessionJWT, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return sbvision.SessionJWT(base64.URLEncoding.EncodeToString(append(signature, data...))), nil
}

func (s *sessionManager) ValidateSession(header sbvision.SessionJWT) (*sbvision.Session, error) {
	data, err := base64.URLEncoding.DecodeString(string(header))
	if err != nil {
		return nil, err
	}

	if len(data) < 130 {
		return nil, fmt.Errorf("\n\tHeader is not long enough, must be at least 130 bytes (was %d)", len(data))
	}

	signature := data[:128]

	hashed := sha256.Sum256(data[128:])

	err = rsa.VerifyPKCS1v15(&s.key.PublicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return nil, fmt.Errorf("\n\tError from verification: %s", err)
	}

	var session sbvision.Session
	err = json.Unmarshal(data[128:], &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
