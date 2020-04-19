package mysqldb

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareDataNearestRotation() (err error) {
	sb.dataNearestRotation, err = sb.db.Prepare(`
SELECT id, bound_id, r, i, j, k
FROM rotations
ORDER BY 
	(rotations.r * ?) +
	(rotations.i * ?) +
	(rotations.j * ?) +
	(rotations.k * ?) 
DESC
LIMIT ?`)
	return
}

// DataNearestRotation looks up the closest rotation to the given quaternion
func (sb *SBDatabase) DataNearestRotation(rot *sbvision.Rotation, dst []sbvision.Rotation) error {
	results, err := sb.dataNearestRotation.Query(rot.R, rot.I, rot.J, rot.K, len(dst))
	if err != nil {
		return fmt.Errorf("\n\tError looking up nearby rotation: %s", err.Error())
	}

	i := 0
	for results.Next() {
		err = results.Scan(&dst[i].ID, &dst[i].BoundID, &dst[i].R, &dst[i].I, &dst[i].J, &dst[i].K)
		if err != nil {
			return err
		}
		i++
	}

	return nil
}
