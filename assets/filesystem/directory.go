package filesystem

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kevinwylder/sbvision"
)

// AssetDirectory is a folder that holds images
type AssetDirectory struct {
	path string
}

// NewAssetDirectory creates an image storage directory that fulfils the sbvision.Image interfaces
func NewAssetDirectory(dir string) (*AssetDirectory, error) {
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
	err = os.MkdirAll(path.Join(dir, "videos"), 0755)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create videos directory: %s", err.Error())
	}
	return &AssetDirectory{
		path: dir,
	}, nil
}

// PutAsset reads the given source and writes it to the file
func (sd *AssetDirectory) PutAsset(key sbvision.Key, data io.Reader) error {
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		return fmt.Errorf("\n\tCannot read image (%s) from reader: %s", key, err)
	}
	err = ioutil.WriteFile(path.Join(sd.path, string(key)), bytes, 0755)
	if err != nil {
		return fmt.Errorf("\n\tCannot create file for image (%s): %s", key, err)
	}
	return nil
}

// GetAsset returns the open file
func (sd *AssetDirectory) GetAsset(image sbvision.Key) (io.ReadCloser, error) {
	file, err := os.Open(path.Join(sd.path, string(image)))
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot open image (%s): %s", image, err)
	}
	return file, nil
}
