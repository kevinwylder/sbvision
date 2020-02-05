package database

import (
	"fmt"

	"github.com/kevinwylder/sbvision"
)

// AddYoutubeRecord adds the record to the database. it is expected that the video is already added
func (db *SBDatabase) AddYoutubeRecord(yt *sbvision.YoutubeVideoInfo) error {
	_, err := db.Exec(`
INSERT INTO youtube_videos (youtube_id, video_id, mirror_url, mirror_expire) 
VALUES (?, ?, ?, FROM_UNIXTIME(?));
	`, yt.YoutubeID, yt.Video.ID, yt.MirrorURL, yt.MirrorExp)
	if err != nil {
		return fmt.Errorf("\n\tCould not insert youtube Record: %s", err.Error())
	}
	return nil
}

// GetYoutubeRecord gets the youtube contextual information about a video given it's generic videoid
func (db *SBDatabase) GetYoutubeRecord(videoID int64) (*sbvision.YoutubeVideoInfo, error) {
	videos, err := db.getVideos(`WHERE videos.id = ?`, "", videoID)
	if err != nil {
		return nil, fmt.Errorf("\n\tError getting videos with id %d: %s", videoID, err.Error())
	}
	if len(videos) == 0 {
		return nil, fmt.Errorf("\n\tNo results for video %d", videoID)
	}
	yt := &sbvision.YoutubeVideoInfo{
		Video: &videos[0],
	}
	result := db.QueryRow(`
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
