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

func (sb *SBDatabase) prepareDataWhereHasBound() (err error) {
	sb.dataWhereHasBound, err = sb.prepareFramesWhere("bounds.id IS NOT NULL")
	return
}

// DataWhereHasBound returns a page of all the data where a bound is known
func (sb *SBDatabase) DataWhereHasBound(offset int64) (*sbvision.FramePage, error) {
	results, err := sb.dataWhereHasBound.Query(offset)
	if err != nil {
		return nil, err
	}
	return parseFrames(results, offset)
}

func (sb *SBDatabase) prepareDataNearestRotation() (err error) {
	sb.dataNearestRotation, err = sb.db.Prepare(fmt.Sprintf(`
SELECT %s
FROM %s
WHERE rotations.id IS NOT NULL
ORDER BY 
	(rotations.r - ?) * (rotations.r - ?) +
	(rotations.i - ?) * (rotations.i - ?) +
	(rotations.j - ?) * (rotations.j - ?) +
	(rotations.k - ?) * (rotations.k - ?) ASC
LIMIT ?`, frameColumns, frameJoin))
	return
}

// DataNearestRotation looks up the closest rotation to the given quaternion
func (sb *SBDatabase) DataNearestRotation(rot *sbvision.Rotation, count int64) (*sbvision.Frame, error) {
	result, err := sb.dataNearestRotation.Query(rot.R, rot.R, rot.I, rot.I, rot.J, rot.J, rot.K, rot.K, count)
	if err != nil {
		return nil, fmt.Errorf("\n\tError looking up nearby rotation: %s", err.Error())
	}
	page, err := parseFrames(result, 0)
	if err != nil {
		return nil, fmt.Errorf("\n\tError parsing frame")
	}
	if len(page.Frames) == 0 {
		return nil, fmt.Errorf("\n\tNo data found")
	}
	return &page.Frames[0], nil
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
ORDER BY frames.video_id DESC, frames.time, bounds.id, rotations.id
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
				Bounds:  make([]sbvision.Bound, 0),
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
				ID:        boundID.Int64,
				FrameID:   frame.ID,
				Height:    height.Int64,
				Width:     width.Int64,
				X:         x.Int64,
				Y:         y.Int64,
				Rotations: make([]sbvision.Rotation, 0),
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
