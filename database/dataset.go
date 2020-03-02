package database

import (
	"database/sql"
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareDataWhereVideo() (err error) {
	sb.dataWhereVideo, err = sb.prepareFramesWhere("frames.video_id = ?")
	return
}

// DataWhereVideo returns a page of data for the given video
func (sb *SBDatabase) DataWhereVideo(videoID int64, offset int64) (*sbvision.FramePage, error) {
	rows, err := sb.dataWhereVideo.Query(videoID, offset)
	if err != nil {
		return nil, fmt.Errorf("\n\tError querying SBDatabase.DataWhereVideo: %s", err.Error())
	}
	result, err := parseFrames(rows, offset)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing SBDatabase.DataWhereVideo: %s", err.Error())
	}
	return result, nil
}

func (sb *SBDatabase) prepareDataWhereFrame() (err error) {
	sb.dataWhereFrame, err = sb.prepareFramesWhere(`
frames.image_hash = ? AND ABS(frames.time - ?) < 1000 AND frames.video_id = ?
`)
	return
}

// DataWhereFrame looks up a frame and returns annotated data about it
func (sb *SBDatabase) DataWhereFrame(hash int64, time int64, videoID int64, offset int64) (*sbvision.FramePage, error) {
	rows, err := sb.dataWhereFrame.Query(hash, time, videoID, offset)
	if err != nil {
		return nil, fmt.Errorf("\n\tError querying SBDatabase.DataWhereFrame: %s", err.Error())
	}
	result, err := parseFrames(rows, offset)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing SBDatabase.DataWhereFrame: %s", err.Error())
	}
	return result, nil
}

func (sb *SBDatabase) prepareDataWhereNoRotation() (err error) {
	sb.dataWhereNoRotation, err = sb.prepareFramesWhere("rotations.id IS NULL AND bounds.id IS NOT NULL")
	return
}

// DataWhereNoRotation gets a page of data where there is a bound but no rotation
func (sb *SBDatabase) DataWhereNoRotation(offset int64) (*sbvision.FramePage, error) {
	rows, err := sb.dataWhereNoRotation.Query(offset)
	if err != nil {
		return nil, fmt.Errorf("\n\tError querying SBDatabase.DataWhereNoRotation: %s", err.Error())
	}
	result, err := parseFrames(rows, offset)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing SBDatabase.DataWhereNoRotation: %s", err.Error())
	}
	return result, nil
}

/**
 * Below are the "generic" parts of extracting a video frame.
**/

const parseLimit int64 = 500

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

func (sb *SBDatabase) prepareFramesWhere(where string) (*sql.Stmt, error) {
	return sb.db.Prepare(fmt.Sprintf(`
SELECT %s
FROM %s
WHERE %s
ORDER BY frames.video_id, frames.time, bounds.id, rotations.id
LIMIT ?, %d`, frameColumns, frameJoin, where, parseLimit))
}

func parseFrames(results *sql.Rows, offset int64) (*sbvision.FramePage, error) {

	page := &sbvision.FramePage{}

	var frameID, frameTime, videoID int64
	var boundID, x, y, width, height, rotationID sql.NullInt64
	var r, i, j, k sql.NullFloat64
	var frame *sbvision.Frame
	var bound *sbvision.Bound

	var resultCount int64
	var frameCount int64

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
			return nil, fmt.Errorf("\n\tError scanning when parsing Data Rows: \n\t\t%s", err.Error())
		}

		if frame == nil || frame.ID != frameID {
			page.Frames = append(page.Frames, sbvision.Frame{
				ID:      frameID,
				VideoID: videoID,
				Time:    frameTime,
			})
			frame = &page.Frames[len(page.Frames)-1]
			resultCount += frameCount
			frameCount = 0
		} else {
			frameCount++
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

	if frameCount+resultCount < 500 {
		return page, nil
	}

	if frameCount != 0 {
		page.Frames = page.Frames[:len(page.Frames)-1]
	}
	page.IsTruncated = true
	page.NextOffset = offset + resultCount
	return page, nil
}
