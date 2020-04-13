package media

import (
	"fmt"
	"io"
	"io/ioutil"
)

// putAsset reads the given source and writes it to the file
func (sd *AssetDirectory) putAsset(key string, data io.Reader) error {
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		return fmt.Errorf("\n\tCannot read image (%s) from reader: %s", key, err)
	}
	err = ioutil.WriteFile(key, bytes, 0755)
	if err != nil {
		return fmt.Errorf("\n\tCannot create file for image (%s): %s", key, err)
	}
	return nil
}

// PutBound stores the data against the given data id
func (sd *AssetDirectory) PutBound(id int64, data io.Reader) error {
	return sd.putAsset(sd.bound(id), data)
}

// PutFrame stores the data against the given data id
func (sd *AssetDirectory) PutFrame(id int64, data io.Reader) error {
	return sd.putAsset(sd.frame(id), data)
}

// PutThumbnail stores the data against the given data id
func (sd *AssetDirectory) PutThumbnail(id int64, data io.Reader) error {
	return sd.putAsset(sd.thumbnail(id), data)
}

// PutVideo stores the data against the given data id
func (sd *AssetDirectory) PutVideo(id int64, data io.Reader) error {
	return sd.putAsset(sd.video(id), data)
}
