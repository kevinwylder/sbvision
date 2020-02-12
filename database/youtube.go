package database

import (
	"fmt"
	"time"

	"github.com/kevinwylder/sbvision"
)

func (sb *SBDatabase) prepareAddYoutubeRecord() (err error) {
	sb.addYoutubeRecord, err = sb.db.Prepare(`
INSERT INTO youtube_videos (youtube_id, video_id, mirror_url, mirror_expire) 
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	mirror_url = VALUES(mirror_url),
	mirror_expire = VALUES(mirror_expire);
	`)
	return
}

// AddYoutubeRecord adds the record to the database. it is expected that the video is already added
func (sb *SBDatabase) AddYoutubeRecord(yt *sbvision.YoutubeVideoInfo) error {
	_, err := sb.addYoutubeRecord.Exec(yt.YoutubeID, yt.VideoID, yt.MirrorURL, yt.MirrorExp)
	if err != nil {
		return fmt.Errorf("\n\tCould not insert youtube Record: %s", err.Error())
	}
	return nil
}

func (sb *SBDatabase) prepareGetYoutubeRecord() (err error) {
	sb.getYoutubeRecord, err = sb.db.Prepare(`
SELECT	
	youtube_id, 
	mirror_url,
	UNIX_TIMESTAMP(mirror_expire)
FROM youtube_videos
WHERE video_id = ?
	`)
	return
}

// GetYoutubeRecord gets the youtube contextual information about a video given it's generic videoid
func (sb *SBDatabase) GetYoutubeRecord(videoID int64) (*sbvision.YoutubeVideoInfo, error) {
	video := sb.getYoutubeRecord.QueryRow(videoID)
	yt := &sbvision.YoutubeVideoInfo{
		VideoID: videoID,
	}
	var unixExpire int64
	err := video.Scan(&yt.YoutubeID, &yt.MirrorURL, &unixExpire)
	if err != nil {
		return nil, fmt.Errorf("\n\tError scanning youtube record: %s", err.Error())
	}
	yt.MirrorExp = time.Unix(unixExpire, 0)
	return yt, nil
}
