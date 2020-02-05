package images

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kevinwylder/sbvision"
)

// ImageBucket is an S3 and local storage synchronized filesystem
type ImageBucket struct {
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	cache      *ImageDirectory
	bucket     *string
}

// NewImageBucket is a constructor for the image bucket
func NewImageBucket(bucket string) (*ImageBucket, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("\n\tCould not create S3 bucket cache dir: %s", err.Error())
	}
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return nil, fmt.Errorf("\n\tCould not create AWS session: %s", err.Error())
	}
	cache, err := NewImageDirectory(dir)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot create image cache directory: %s", err)
	}
	return &ImageBucket{
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
		cache:      cache,
		bucket:     aws.String(bucket),
	}, nil
}

// PutImage puts the given data in the bucket
func (sb *ImageBucket) PutImage(data io.Reader, key sbvision.Image) error {
	_, err := sb.uploader.Upload(&s3manager.UploadInput{
		Body:   data,
		Bucket: sb.bucket,
		Key:    aws.String(string(key)),
	})
	if err != nil {
		return fmt.Errorf("\n\tCould not upload image: %s", err.Error())
	}
	return nil
}

// GetImage looks gets the local image, downloading it if necessary
func (sb *ImageBucket) GetImage(image sbvision.Image) (io.ReadCloser, error) {
	file, err := sb.cache.GetImage(image)
	if err != nil {
		// download the file
		filePath := path.Join(sb.cache.path, string(image))
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return nil, fmt.Errorf("\n\tError opening file (%s) for s3 download: %s", filePath, err.Error())
		}
		_, err = sb.downloader.Download(file, &s3.GetObjectInput{
			Key:    aws.String(string(image)),
			Bucket: sb.bucket,
		})
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("\n\tCould not download key (%s) from s3: %s", image, err.Error())
		}
		file, err = os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("\n\tError re-opening downloaded file (%s): %s", filePath, err.Error())
		}
	}
	return file, nil
}

// ClearCache removes the temporary cache directory
func (sb *ImageBucket) ClearCache() {
	os.RemoveAll(sb.cache.path)
}
