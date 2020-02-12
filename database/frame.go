package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddFrame() (err error) {
	sb.addFrame, err = sb.db.Prepare(`
INSERT INTO frames (image_id, video_id, time) 
SELECT 
	images.id, ?, ?
FROM images
WHERE images.key = ?
	`)
	return
}

// AddFrame adds the given frame to the database and fills in the autoincrement
func (sb *SBDatabase) AddFrame(frame *sbvision.Frame) error {
	result, err := sb.addFrame.Exec(frame.VideoID, frame.Time, frame.Image)
	if err != nil {
		return fmt.Errorf("\n\tError adding frame: %s", err.Error())
	}
	frame.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("\n\tError getting added frame id: %s", err.Error())
	}
	return nil
}
