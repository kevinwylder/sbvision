package cdn

import (
	"fmt"
)

// VideoDirectory returns the path for a video directory
func VideoDirectory(id string) string {
	return fmt.Sprintf("/video/%s", id)
}

// ClipDirectory returns the path for a clip
func ClipDirectory(id string) string {
	return fmt.Sprintf("/clip/%s", id)
}
