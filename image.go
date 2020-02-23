package sbvision

// Image is handled by the /image route
type Image string

// ImageHash is a hash of the image data
type ImageHash string

// ImageTracker adds images to the database
type ImageTracker interface {
	AddImage(Image, *Session) error
}
