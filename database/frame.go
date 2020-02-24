package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddFrame() (err error) {
	sb.addFrame, err = sb.db.Prepare(`
INSERT INTO frames (video_id, time, image_hash, session_id) 
VALUES (?, ?, ?, ?);
	`)
	return
}

// AddFrame adds the given frame to the database and fills in the autoincrement
func (sb *SBDatabase) AddFrame(frame *sbvision.Frame, session *sbvision.Session, hash int64) error {
	result, err := sb.addFrame.Exec(frame.VideoID, frame.Time, hash, session.ID)
	if err != nil {
		return fmt.Errorf("\n\tError adding frame: %s", err.Error())
	}
	frame.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("\n\tError getting added frame id: %s", err.Error())
	}
	return nil
}
