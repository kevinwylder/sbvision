package sbvision

// Frame is a frame of a video
type Frame struct {
	ID    int64
	Image Image  `json:"image"`
	Video *Video `json:"video"`
	Time  int64  `json:"time"`
}

// Clip is the image that the user clipped
type Clip struct {
	ID      int64
	Session *Session
	Frame   *Frame
	Width   int64   `json:"width"`
	Height  int64   `json:"height"`
	X       int64   `json:"x"`
	Y       int64   `json:"y"`
	R       float64 `json:"r"`
	I       float64 `json:"i"`
	J       float64 `json:"j"`
	K       float64 `json:"k"`
}
