package database

import "github.com/kevinwylder/sbvision"

// AddImage adds the given image to the database. It is expected that it was already uploaded
func (db *SBDatabase) AddImage(image sbvision.Image, session *sbvision.Session) error {
	_, err := db.Exec(`
INSERT INTO images (key, session_id) VALUES (?, ?);
	`, image, session.ID)
	if err != nil {
		return err
	}
	return nil
}
