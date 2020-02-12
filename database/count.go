package database

import "fmt"

func (sb *SBDatabase) prepareDataCounts() (err error) {
	sb.dataCounts, err = sb.db.Prepare(`
SELECT 
	COUNT(DISTINCT frames.id, frames.id IS NOT NULL),
	COUNT(DISTINCT bounds.id, bounds.id IS NOT NULL),
	COUNT(DISTINCT rotations.id, rotations.id IS NOT NULL)
FROM frames
LEFT JOIN bounds
		ON bounds.frame_id = frames.id
LEFT JOIN rotations
		ON rotations.bounds_id = bounds.id;
	`)
	return
}

// DataCounts gives a high level aggregation of the data
func (sb *SBDatabase) DataCounts() (frames int64, bounds int64, rotations int64, err error) {
	result := sb.dataCounts.QueryRow()
	err = result.Scan(&frames, &bounds, &rotations)
	return
}

func (sb *SBDatabase) prepareGetVideoCount() (err error) {
	sb.getVideoCount, err = sb.db.Prepare(`
SELECT COUNT(*) FROM videos
	`)
	return
}

// GetVideoCount gets the total number of videos in the database
func (sb *SBDatabase) GetVideoCount() (int64, error) {
	result := sb.getVideoCount.QueryRow()
	var count int64
	err := result.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("\n\tError getting video count from db: %s", err.Error())
	}
	return count, nil
}
