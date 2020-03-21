package amazon

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kevinwylder/sbvision"
)

// S3Bucket is an S3 and local storage synchronized filesystem
type S3Bucket struct {
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	cache      sbvision.KeyValueStore
	bucket     *string
}

// NewS3BucketManager is a constructor for the image bucket
func NewS3BucketManager(bucket string, cache sbvision.KeyValueStore) (*S3Bucket, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		return nil, fmt.Errorf("\n\tCould not create AWS session: %s", err.Error())
	}
	return &S3Bucket{
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
		cache:      cache,
		bucket:     aws.String(bucket),
	}, nil
}

// PutAsset puts the given data in the bucket
func (sb *S3Bucket) PutAsset(key sbvision.Key, data io.Reader) error {
	if sb.cache != nil {
		err := sb.cache.PutAsset(key, data)
		if err != nil {
			return fmt.Errorf("\n\tCould not store asset in cache: %s", err)
		}
		data, err = sb.cache.GetAsset(key)
		if err != nil {
			return fmt.Errorf("\n\tCould not read cached asset: %s", err)
		}
	}
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, data)
	if err != nil {
		return fmt.Errorf("\n\tCould not read data into buffer: %s", err)
	}
	_, err = sb.uploader.Upload(&s3manager.UploadInput{
		Body:   &buffer,
		Bucket: sb.bucket,
		Key:    aws.String(string(key)),
	})
	if err != nil {
		return fmt.Errorf("\n\tCould not upload asset to s3 bucket: %s", err)
	}

	return nil
}

type removeTmpFile struct {
	*os.File
}

// GetAsset looks gets the local image, downloading it if necessary
func (sb *S3Bucket) GetAsset(key sbvision.Key) (io.ReadCloser, error) {
	if sb.cache != nil {
		data, err := sb.cache.GetAsset(key)
		if err == nil {
			return data, nil
		}
	}
	// download the file
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("\n\tError opening tmp file for s3 download: %s", err.Error())
	}

	_, err = sb.downloader.Download(file, &s3.GetObjectInput{
		Key:    aws.String(string(key)),
		Bucket: sb.bucket,
	})
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return nil, fmt.Errorf("\n\tCould not download key (%s) from s3: %s", key, err.Error())
	}

	if sb.cache != nil {
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()
		err := sb.cache.PutAsset(key, file)
		if err != nil {
			return nil, fmt.Errorf("\n\tDownloaded key (%s) from S3 lost by cache write: %s", key, err)
		}
		data, err := sb.cache.GetAsset(key)
		if err != nil {
			return nil, fmt.Errorf("\n\tDownloaded key (%s) from S3 lost by cache read: %s", key, err)
		}
		return data, err
	}

	return &removeTmpFile{file}, nil
}

func (f *removeTmpFile) Read(data []byte) (int, error) {
	return f.Read(data)
}

func (f *removeTmpFile) Close() error {
	f.Close()
	return os.Remove(f.Name())
}
