package sbvision

// Image is an object stored in s3
type Image struct {
	ID  int64
	Key string
	URL string `json:"url"`
}
