package auth

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
	"github.com/lestrrat-go/jwx/jwk"
)

// Database is the required functionality to lookup a user
type Database interface {
	GetUser(email string) (*sbvision.User, error)
}

// JWTVerifier is able to validate JWTs
type JWTVerifier struct {
	claimsURL string
	keyset    *jwk.Set

	db    Database
	cache map[string]*sbvision.User
}

// NewJWTVerifier creates a cache of auth data to verify tokens
// ClaimsURL is the .well-known/jwk.json directory of the auth server
func NewJWTVerifier(db Database, claimsURL string) *JWTVerifier {
	return &JWTVerifier{
		claimsURL: claimsURL,
		db:        db,
		cache:     make(map[string]*sbvision.User),
	}
}

// User checks the local cache for the given user
func (j *JWTVerifier) User(data string) (*sbvision.User, error) {
	email, err := j.Verify(data)
	if err != nil {
		return nil, err
	}

	if user, exists := j.cache[email]; exists {
		return user, nil
	}

	user, err := j.db.GetUser(email)
	if err != nil {
		return nil, fmt.Errorf("User not in database")
	}

	j.cache[email] = user
	return user, nil
}
