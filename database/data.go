package database

import (
	"database/sql"
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareDataByBoundID() (err error) {
	sb.dataByBoundID, err = sb.db.Prepare(`
SELECT ` + frameColumns + `
FROM ` + frameJoin + `
WHERE bounds.id = ?
	`)
	return
}

// DataByBoundID gets a frame and it's rotations by the given bound id
func (sb *SBDatabase) DataByBoundID(id int64) (*sbvision.Frame, error) {
	result, err := sb.dataByBoundID.Query(id)
	if err != nil {
		return nil, fmt.Errorf("\n\tError querying frames for ByBoundID: %s", err.Error())
	}
	defer result.Close()
	frames, err := parseFrames(result)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing frames for ByBoundID: %s", err.Error())
	}
	if len(frames) != 1 {
		return nil, fmt.Errorf("\n\tBound not found")
	}
	return &frames[0], nil
}

func (sb *SBDatabase) prepareDataRotationFrames() (err error) {
	sb.dataRotationFrames, err = sb.db.Prepare(`
SELECT ` + frameColumns + `
FROM ` + frameJoin + `
WHERE rotations.id IS NULL AND bounds.id IS NOT NULL
ORDER BY frames.video_id ASC, frames.time ASC
LIMIT 100;
	`)
	return
}

// DataRotationFrames gets the frames that need an orientation
func (sb *SBDatabase) DataRotationFrames() ([]sbvision.Frame, error) {
	results, err := sb.dataRotationFrames.Query()
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting rotation frames: %s", err.Error())
	}
	defer results.Close()
	frames, err := parseFrames(results)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing rotation frames: %s", err.Error())
	}
	return frames, nil
}

func (sb *SBDatabase) prepareGetFrame() (err error) {
	sb.getFrame, err = sb.db.Prepare(`
SELECT ` + frameColumns + `
FROM ` + frameJoin + `
WHERE (
	frames.time = ? 
	OR
	(
		frames.image_hash = ?
		AND 
		ABS(frames.time - ?) < 1000
	)
) 
AND
frames.video_id = ?
	`)
	return
}

// GetFrame returns a frame object for the given video
func (sb *SBDatabase) GetFrame(video int64, frameNum int64, hash int64) (*sbvision.Frame, error) {
	result, err := sb.getFrame.Query(frameNum, hash, frameNum, video)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting frame %d of video %d: %s", frameNum, video, err)
	}
	defer result.Close()
	frames, err := parseFrames(result)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing frame %d of video %d: %s", frameNum, video, err)
	}
	if len(frames) == 0 {
		return nil, nil
	}
	return &frames[0], nil
}

func (sb *SBDatabase) prepareDataAllFrames() (err error) {
	sb.dataAllFrames, err = sb.db.Prepare(`
SELECT ` + frameColumns + `
FROM ` + frameJoin + `
ORDER BY frames.video_id ASC, frames.time ASC
	`)
	return
}

// DataAllFrames gets all teh datas
func (sb *SBDatabase) DataAllFrames() ([]sbvision.Frame, error) {
	rows, err := sb.dataAllFrames.Query()
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting data frames: %s", err.Error())
	}
	frames, err := parseFrames(rows)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing all frames: %s", err.Error())
	}
	return frames, nil
}

func (sb *SBDatabase) prepareDataVideoFrames() (err error) {
	sb.dataVideoFrames, err = sb.db.Prepare(`
SELECT ` + frameColumns + `
FROM ` + frameJoin + `
WHERE frames.video_id = ?
ORDER BY frames.time ASC`)
	return
}

// DataVideoFrames gets all the data for a given video
func (sb *SBDatabase) DataVideoFrames(videoID int64) ([]sbvision.Frame, error) {
	result, err := sb.dataVideoFrames.Query(videoID)
	if err != nil {
		return nil, fmt.Errorf("\n\tError querying video data frames: %s", err.Error())
	}
	defer result.Close()
	frames, err := parseFrames(result)
	if err != nil {
		return nil, fmt.Errorf("\n\tError scanning video data frame results: %s", err.Error())
	}
	return frames, nil
}

const frameColumns = `
	frames.id,
	frames.time,
	frames.video_id,
	bounds.id,
	bounds.x,
	bounds.y,
	bounds.width,
	bounds.height,
	rotations.id,
	rotations.r,
	rotations.i,
	rotations.j,
	rotations.k `

const frameJoin = ` frames
LEFT JOIN bounds
		ON bounds.frame_id = frames.id
LEFT JOIN rotations
		ON rotations.bounds_id = bounds.id
`

func parseFrames(results *sql.Rows) ([]sbvision.Frame, error) {
	var frameID, frameTime, videoID int64
	var boundID, x, y, width, height, rotationID sql.NullInt64
	var r, i, j, k sql.NullFloat64
	var frames []sbvision.Frame
	var frame *sbvision.Frame
	var bound *sbvision.Bound
	for results.Next() {
		err := results.Scan(
			&frameID,
			&frameTime,
			&videoID,
			&boundID,
			&x,
			&y,
			&width,
			&height,
			&rotationID,
			&r,
			&i,
			&j,
			&k,
		)
		if err != nil {
			return nil, err
		}

		if frame == nil || frame.ID != frameID {
			frames = append(frames, sbvision.Frame{
				ID:      frameID,
				VideoID: videoID,
				Time:    frameTime,
			})
			frame = &frames[len(frames)-1]
		}
		if !boundID.Valid {
			bound = nil
			continue
		}
		if bound == nil || bound.ID != boundID.Int64 {
			frame.Bounds = append(frame.Bounds, sbvision.Bound{
				ID:      boundID.Int64,
				FrameID: frame.ID,
				Height:  height.Int64,
				Width:   width.Int64,
				X:       x.Int64,
				Y:       y.Int64,
			})
			bound = &frame.Bounds[len(frame.Bounds)-1]
		}
		if !rotationID.Valid {
			continue
		}
		bound.Rotations = append(bound.Rotations, sbvision.Rotation{
			ID:      rotationID.Int64,
			BoundID: bound.ID,
			I:       i.Float64,
			J:       j.Float64,
			K:       k.Float64,
			R:       r.Float64,
		})
	}
	return frames, nil
}
