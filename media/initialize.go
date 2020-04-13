package media

import (
	"fmt"
	"os"
	"path"
)

// NewAssetDirectory creates an image storage directory and ensures all the necessary directories exist
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
	err = os.MkdirAll(path.Join(dir, "frame"), 0755)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create frame directory: %s", err.Error())
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
		path: dir,
	}, nil
}
