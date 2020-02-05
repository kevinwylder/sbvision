package s3

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
	tmpdir     string
	bucket     *string
}

// NewImageBucket is a constructor for the image bucket
func NewImageBucket() (*ImageBucket, error) {
	if os.Getenv("S3_BUCKET") == "" {
		return nil, fmt.Errorf("\n\tEmpty S3_BUCKET env variable")
	}
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
	return &ImageBucket{
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
		tmpdir:     dir,
		bucket:     aws.String(os.Getenv("S3_BUCKET")),
	}, nil
}

// UploadImage puts the given data in the bucket
func (sb *ImageBucket) UploadImage(data io.Reader, key string) (sbvision.Image, error) {
	_, err := sb.uploader.Upload(&s3manager.UploadInput{
		Body:   data,
		Bucket: sb.bucket,
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("\n\tCould not upload image: %s", err.Error())
	}
	return sbvision.Image(key), nil
}

// DownloadImage stores the file locally and serves that
func (sb *ImageBucket) DownloadImage(image sbvision.Image) (io.ReadCloser, error) {
	filePath := path.Join(sb.tmpdir, string(image))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// download the file
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return nil, fmt.Errorf("Error opening file for s3 download: %s", err.Error())
		}
		_, err = sb.downloader.Download(file, &s3.GetObjectInput{
			Key:    aws.String(string(image)),
			Bucket: sb.bucket,
		})
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("Could not download file from s3: %s", err.Error())
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to (re)open file: %s", err.Error())
	}
	return file, err
}
