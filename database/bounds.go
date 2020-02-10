package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

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
