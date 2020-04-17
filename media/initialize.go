package media

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

// AssetDirectory is a folder that holds stuff and should be put on a CDN
type AssetDirectory struct {
	path      string
	ServeHTTP http.HandlerFunc
}

// NewAssetDirectory creates a storage directory and ensures all the necessary directories exist
func NewAssetDirectory(dir string) (*AssetDirectory, error) {
	if dir == "" {
		return nil, fmt.Errorf(`Cannot use "" for data storage`)
	}
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("\n\tCannot create the given dir (%s): %s", dir, err.Error())
		}
	} else if err != nil {
		return nil, fmt.Errorf("\n\tCannot stat the given dir (%s): %s", dir, err.Error())
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("\n\tGiven path (%s) is a file", dir)
	}
	err = os.MkdirAll(path.Join(dir, "thumbnail"), 0755)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create thumbnail directory: %s", err.Error())
	}
	err = os.MkdirAll(path.Join(dir, "bound"), 0755)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create bound directory: %s", err.Error())
	}
	err = os.MkdirAll(path.Join(dir, "video"), 0755)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create video directory: %s", err.Error())
	}
	return &AssetDirectory{
		path:      dir,
		ServeHTTP: http.FileServer(http.Dir(dir)).ServeHTTP,
	}, nil
}
