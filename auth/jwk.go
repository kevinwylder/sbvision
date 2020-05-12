package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

func (j *JWTVerifier) getKey(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Invalid signing method")
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	if j.keyset == nil {
		keyset, err := jwk.FetchHTTP(j.claimsURL)
		if err != nil {
			return nil, err
		}
		j.keyset = keyset
	}

	key := j.keyset.LookupKeyID(keyID)
	if len(key) == 0 {
		return nil, fmt.Errorf("Could not find that key type")
	}

	var raw interface{}
	return raw, key[0].Raw(&raw)
}

// Verify gets the email from the given key
func (j *JWTVerifier) Verify(data string) (string, error) {

	jwt.TimeFunc = func() time.Time {
		return time.Now().Add(time.Second)
	}

	token, err := jwt.Parse(data, j.getKey)
	if err != nil {
		return "", err
	}

	claims, safe := token.Claims.(jwt.MapClaims)
	if !safe || !token.Valid {
		return "", fmt.Errorf("token is not of type jwt.MapClaims")
	}

	email, exists := claims["email"]
	if !exists {
		return "", fmt.Errorf("JWT has no email claim")
	}

	value, safe := email.(string)
	if !safe {
		return "", fmt.Errorf("Email claim is not a string")
	}

	return value, nil
}
