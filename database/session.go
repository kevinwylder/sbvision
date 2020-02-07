package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddSession() (err error) {
	sb.addSession, err = sb.db.Prepare(`
INSERT INTO sessions (start, source_ip) VALUES (FROM_UNIXTIME(?), ?);
	`)
	return
}

// AddSession puts the session in the database and fills out the session ID
func (sb *SBDatabase) AddSession(session *sbvision.Session) error {
	result, err := sb.addSession.Exec(session.Time, session.IP)
	if err != nil {
		return fmt.Errorf("\n\tError adding session to db: %s", err.Error())
	}
	session.ID, err = result.LastInsertId()
	return err
}
