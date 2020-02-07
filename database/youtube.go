package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddYoutubeRecord() (err error) {
	sb.addYoutubeRecord, err = sb.db.Prepare(`
INSERT INTO youtube_videos (youtube_id, video_id, mirror_url, mirror_expire) 
VALUES (?, ?, ?, FROM_UNIXTIME(?));
	`)
	return
}

// AddYoutubeRecord adds the record to the database. it is expected that the video is already added
func (sb *SBDatabase) AddYoutubeRecord(yt *sbvision.YoutubeVideoInfo) error {
	_, err := sb.addYoutubeRecord.Exec(yt.YoutubeID, yt.Video.ID, yt.MirrorURL, yt.MirrorExp)
	if err != nil {
		return fmt.Errorf("\n\tCould not insert youtube Record: %s", err.Error())
	}
	return nil
}

func (sb *SBDatabase) prepareGetYoutubeRecord() (err error) {
	sb.getYoutubeRecord, err = sb.db.Prepare(`
SELECT	
	videos.id,
	videos.title,
	images.key,
	videos.type,
	videos.format,
	videos.duration,
	videos.fps,
	COUNT(*),
	MAX(clips.id)
FROM videos
INNER JOIN images 
		ON images.id = videos.thumbnail_id
LEFT JOIN frames
		ON frames.video_id = videos.id
LEFT JOIN clips
		ON clips.frame_id = frames.id
WHERE videos.id = ?
GROUP BY videos.id
	`)
	return
}

// GetYoutubeRecord gets the youtube contextual information about a video given it's generic videoid
func (sb *SBDatabase) GetYoutubeRecord(videoID int64) (*sbvision.YoutubeVideoInfo, error) {
	video := sb.getYoutubeRecord.QueryRow(videoID)
	yt := &sbvision.YoutubeVideoInfo{
		Video: &sbvision.Video{},
	}
	err := sb.parseVideoRow(video, yt.Video)
	if err != nil {
		return nil, fmt.Errorf("\n\tError scanning youtube record: %s", err.Error())
	}
	result := sb.db.QueryRow(`
SELECT
	youtube_id,
	mirror_url,
	mirror_expire
FROM youtube_videos
WHERE video_id = ?
LIMIT 1
	`, videoID)
	err = result.Scan(&yt.YoutubeID, &yt.MirrorURL, &yt.MirrorExp)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting youtube information for video %d", videoID)
	}
	return yt, nil
}
