package sbvision

// User comes from the cognito user pool
type User struct {
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Videos   []string `json:"videos"`
	Clips    []string `json:"clips"`
}

// VideoType is an enum of subclasses of the video abstraction
type VideoType int64

const (
	// YoutubeVideo type means that this video has a record in the youtube_videos table
	YoutubeVideo VideoType = 1
	// RedditVideo type means that this video has a record in the reddit_videos table
	RedditVideo VideoType = 2
	// UploadedVideo type means this video was uploaded
	UploadedVideo VideoType = 3
)

// Video is a generic video source
type Video struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Width         int64     `json:"width"`
	Height        int64     `json:"height"`
	FPS           float64   `json:"fps"`
	Duration      string    `json:"duration"`
	Type          VideoType `json:"type"`
	UploadedAt    string    `json:"uploaded_at"`
	UploadedBy    string    `json:"uploaded_by"`
	UploaderEmail string    `json:"-" dynamodbav:"uploader_email,string"`
	SourceURL     string    `json:"source" dynamodbav:"source_url,string"`
}

// Clip is part of a video that has a trick
type Clip struct {
	ID             string               `json:"id"`
	VideoID        string               `json:"videoId"`
	Username       string               `json:"clipped_by"`
	Trick          string               `json:"trick"`
	UploadedAt     string               `json:"uploaded_at"`
	OriginalSource string               `json:"source"`
	Start          int64                `json:"startFrame"`
	End            int64                `json:"endFrame"`
	Bounds         map[int64]Bound      `json:"boxes"`
	Rotations      map[int64]Quaternion `json:"rotations"`
}

// Frame is one annotated frame of data
type Frame struct {
	Image    string     `json:"image"`
	Bound    Bound      `json:"bound"`
	Rotation Quaternion `json:"rotation"`
}

// Bound is an area on a frame
type Bound struct {
	Width  int64 `json:"w"`
	Height int64 `json:"h"`
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
}

type Quaternion [4]float64
