package sbvision

// Frame is a frame of a video
type Frame struct {
	ID      int64 `json:"id"`
	Image   Image `json:"image"`
	VideoID int64 `json:"video"`
	Time    int64 `json:"time"`
}

// Bound is an area on a frame
type Bound struct {
	ID        int64      `json:"id"`
	FrameID   int64      `json:"frameId"`
	Width     int64      `json:"width"`
	Height    int64      `json:"height"`
	X         int64      `json:"x"`
	Y         int64      `json:"y"`
	Rotations []Rotation `json:"rotations"`
}

// Rotation is the angle that a user has voted for a bound
type Rotation struct {
	BoundID int64   `json:"boundId"`
	ID      int64   `json:"id"`
	R       float64 `json:"r"`
	I       float64 `json:"i"`
	J       float64 `json:"j"`
	K       float64 `json:"k"`
}
