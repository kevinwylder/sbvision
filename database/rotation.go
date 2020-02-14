package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

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
