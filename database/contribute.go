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

func (sb *SBDatabase) prepareAddBounds() (err error) {
	sb.addBounds, err = sb.db.Prepare(`
INSERT INTO bounds (session_id, frame_id, x, y, width, height) VALUES (?, ?, ?, ?, ?, ?);
	`)
	return
}

// AddBounds stores the bounds in the database and updateas the bounds pointer with the new ID
func (sb *SBDatabase) AddBounds(bounds *sbvision.Bound, session *sbvision.Session) error {
	result, err := sb.addBounds.Exec(session.ID, bounds.FrameID, bounds.X, bounds.Y, bounds.Width, bounds.Height)
	if err != nil {
		return fmt.Errorf("\n\tError executing addBounds: %s", err.Error())
	}
	bounds.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("\n\tError getting addBounds insert id: %s ", err.Error())
	}
	return nil
}

func (sb *SBDatabase) prepareAddRotation() (err error) {
	sb.addRotation, err = sb.db.Prepare(`
INSERT INTO rotations (bounds_id, session_id, r, i, j, k) VALUES (?, ?, ?, ?, ?, ?);
	`)
	return
}

// AddRotation adds the rotation to the database and fills out it's id
func (sb *SBDatabase) AddRotation(rotation *sbvision.Rotation, session *sbvision.Session) error {
	result, err := sb.addRotation.Exec(rotation.BoundID, session.ID, rotation.R, rotation.I, rotation.J, rotation.K)
	if err != nil {
		return fmt.Errorf("\n\tError executing add rotation: %s", err.Error())
	}
	rotation.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("\n\tError assigning rotation id: %s", err.Error())
	}
	return nil
}
