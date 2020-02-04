package session_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kevinwylder/sbvision"
	"github.com/kevinwylder/sbvision/session"
)

func TestJWTSignAndValidate(t *testing.T) {
	manager, err := session.NewRSASessionManager()
	if err != nil {
		t.Fatal(err)
	}

	session := &sbvision.Session{
		ID:   10,
		IP:   "192.168.1.1",
		Time: time.Now().Unix(),
	}
	header, err := manager.SignSession(session)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := manager.ValidateSession(header)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.IP != session.IP {
		t.Fail()
	}
	if decoded.ID != session.ID {
		t.Fail()
	}
	if decoded.Time != session.Time {
		t.Fail()
	}

	fmt.Println(header)
	t.Fail()
}
