package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddImage() (err error) {
	sb.addImage, err = sb.db.Prepare(`
INSERT INTO images (` + "`key`" + `, session_id) VALUES (?, ?);
	`)
	return
}

// AddImage adds the given image to the database. It is expected that it was already uploaded
func (sb *SBDatabase) AddImage(image sbvision.Image, session *sbvision.Session) error {
	_, err := sb.addImage.Exec(image, session.ID)
	if err != nil {
		return fmt.Errorf("\n\tError Adding image to the database: %s", err.Error())
	}
	return nil
}
