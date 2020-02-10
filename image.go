package sbvision

// Image is handled by the /image route
type Image string

// ImageTracker adds images to the database
type ImageTracker interface {
	AddImage(Image, *Session) error
}
