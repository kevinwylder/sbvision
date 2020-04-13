package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

func (j *JWTVerifier) getKey(token *jwt.Token) (interface{}, error) {

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

	if key := j.keyset.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, fmt.Errorf("unable to find key %q", keyID)
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
