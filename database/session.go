package database

import "github.com/kevinwylder/sbvision"

// AddSession puts the session in the database and fills out the session ID
// this makes it comply with the SessionTracker interface
func (db *SBDatabase) AddSession(session *sbvision.Session) error {
	result, err := db.Exec(`
INSERT INTO sessions (start, source_ip) VALUES (?, ?);
	`, session.Time, session.IP)
	if err != nil {
		return err
	}
	session.ID, err = result.LastInsertId()
	return err
}
