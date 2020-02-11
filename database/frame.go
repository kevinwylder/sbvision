package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareGetFrame() (err error) {
	sb.getFrame, err = sb.db.Prepare(`
SELECT 
	frames.id,
	images.key
FROM frames 
INNER JOIN images
		ON images.id = frames.image_id
WHERE frames.time = ? AND frames.video_id = ?
	`)
	return
}

// GetFrame returns a frame object for the given video
func (sb *SBDatabase) GetFrame(video int64, frameNum int64) (*sbvision.Frame, error) {
	result := sb.getFrame.QueryRow(frameNum, video)
	frame := &sbvision.Frame{
		Time:    frameNum,
		VideoID: video,
	}
	err := result.Scan(&frame.ID, &frame.Image)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting frame %d of video %d: %s", frameNum, video, err)
	}
	return frame, nil
}

func (sb *SBDatabase) prepareAddFrame() (err error) {
	sb.addFrame, err = sb.db.Prepare(`
INSERT INTO frames (image_id, video_id, time) 
SELECT 
	images.id, ?, ?
FROM images
WHERE images.key = ?
	`)
	return
}

// AddFrame adds the given frame to the database and fills in the autoincrement
func (sb *SBDatabase) AddFrame(frame *sbvision.Frame) error {
	result, err := sb.addFrame.Exec(frame.VideoID, frame.Time, frame.Image)
	if err != nil {
		return fmt.Errorf("\n\tError adding frame: %s", err.Error())
	}
	frame.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("\n\tError getting added frame id: %s", err.Error())
	}
	return nil
}
