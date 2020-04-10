package database

import (
	"database/sql"
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareDataWhereVideo() (err error) {
	sb.dataWhereVideo, err = sb.prepareFramesWhere("frames.video_id = ? AND bounds.id IS NOT NULL")
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

	page := &sbvision.FramePage{
		Frames: make([]sbvision.Frame, 0),
	}

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
