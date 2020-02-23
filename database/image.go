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

// GetAllImages returns all the image Keys from the database
func (sb *SBDatabase) GetAllImages() ([]string, error) {
	rows, err := sb.db.Query(` SELECT ` + "`key`" + ` FROM images WHERE ` + "`key`" + `NOT LIKE "%.jpg" `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var image string
		rows.Scan(&image)
		images = append(images, image)
	}
	return images, nil
}

func (sb *SBDatabase) prepareSetFrameHash() (err error) {
	sb.setFrameHash, err = sb.db.Prepare(`
UPDATE frames
SET frames.image_hash = ?
WHERE frames.id = ?
	`)
	return
}

// SetFrameHash updates the hash of the frame
func (sb *SBDatabase) SetFrameHash(hash int64, frameID int64) error {
	_, err := sb.setFrameHash.Exec(hash, frameID)
	if err != nil {
		return fmt.Errorf("\n\tError updating hash: %s", err.Error())
	}
	return nil
}
