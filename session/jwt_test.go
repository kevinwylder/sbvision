package session_test

import (
	"testing"
	"time"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/session"
)

func TestJWTSignAndValidate(t *testing.T) {
	// establish an identity
	manager, err := session.NewRSASessionManager()
	if err != nil {
		t.Fatal(err)
	}

	// create and sign a session
	s := &sbvision.Session{
		ID:   10,
		IP:   "192.168.1.1",
		Time: time.Now().Unix(),
	}
	header, err := manager.CreateSession(s)
	if err != nil {
		t.Fatal(err)
	}

	// make sure that we can decode the signature
	decoded, err := manager.ValidateSession(header)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.IP != s.IP {
		t.Fail()
	}
	if decoded.ID != s.ID {
		t.Fail()
	}
	if decoded.Time != s.Time {
		t.Fail()
	}

	// establish a different identity
	manager2, err := session.NewRSASessionManager()
	if err != nil {
		t.Fatal(err)
	}
	// make sure they do not think the original session is valid
	_, err = manager2.ValidateSession(header)
	if err == nil {
		t.Fail()
	}

}
