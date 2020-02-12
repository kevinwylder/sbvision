package database

import (
	"database/sql"
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareGetFrame() (err error) {
	sb.getFrame, err = sb.db.Prepare(`
SELECT ` + frameColumns + `
FROM ` + frameJoin + `
WHERE frames.time = ? AND frames.video_id = ?
	`)
	return
}

// GetFrame returns a frame object for the given video
func (sb *SBDatabase) GetFrame(video int64, frameNum int64) (*sbvision.Frame, error) {
	result, err := sb.getFrame.Query(frameNum, video)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting frame %d of video %d: %s", frameNum, video, err)
	}
	frames, err := parseFrames(result)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing frame %d of video %d: %s", frameNum, video, err)
	}
	if len(frames) == 0 {
		return nil, nil
	}
	return &frames[0], nil
}

func (sb *SBDatabase) prepareDataVideoFrames() (err error) {
	sb.dataVideoFrames, err = sb.db.Prepare(`
SELECT ` + frameColumns + `
FROM ` + frameJoin + `
WHERE frames.video_id = ?`)
	return
}

// DataVideoFrames gets all the data for a given video
func (sb *SBDatabase) DataVideoFrames(videoID int64) ([]sbvision.Frame, error) {
	result, err := sb.dataVideoFrames.Query(videoID)
	if err != nil {
		return nil, fmt.Errorf("\n\tError querying video data frames: %s", err.Error())
	}
	frames, err := parseFrames(result)
	if err != nil {
		return nil, fmt.Errorf("\n\tError scanning video data frame results: %s", err.Error())
	}
	return frames, nil
}

const frameColumns = `
	frames.id,
	images.key,
	frames.time,
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
INNER JOIN images
		ON images.id = frames.image_id
LEFT JOIN bounds
		ON bounds.frame_id = frames.id
LEFT JOIN rotations
		ON rotations.bounds_id = bounds.id
`

func parseFrames(results *sql.Rows) ([]sbvision.Frame, error) {
	var frameID, frameTime int64
	var image string
	var boundID, x, y, width, height, rotationID sql.NullInt64
	var r, i, j, k sql.NullFloat64
	var frames []sbvision.Frame
	var frame *sbvision.Frame
	var bound *sbvision.Bound
	for results.Next() {
		err := results.Scan(
			&frameID,
			&image,
			&frameTime,
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
				ID:    frameID,
				Time:  frameTime,
				Image: sbvision.Image(image),
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
