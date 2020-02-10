package amazon

import (
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
		Region: aws.String("us-west-2"),
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
func (sb *S3Bucket) PutAsset(key string, data io.Reader) error {
	if sb.cache == nil {
		// upload the file
		_, err := sb.uploader.Upload(&s3manager.UploadInput{
			Body:   data,
			Bucket: sb.bucket,
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("\n\tCould not upload asset to s3 bucket without cache: %s", err)
		}
		return nil
	}

	pr, pw := io.Pipe()
	tee := io.TeeReader(data, pw)
	sync := make(chan error)
	defer close(sync)

	go func() {
		_, err := sb.uploader.Upload(&s3manager.UploadInput{
			Body:   tee,
			Bucket: sb.bucket,
			Key:    aws.String(string(key)),
		})
		sync <- err
	}()

	go func() {
		defer pr.Close()
		err := sb.cache.PutAsset(key, pr)
		sync <- err
	}()

	err1 := <-sync
	err2 := <-sync
	if err1 != nil {
		return fmt.Errorf("\n\tCould not upload image to s3 with cache: %s", err1)
	}
	if err2 != nil {
		return fmt.Errorf("\n\tCould not put asset in cache: %s", err2)
	}
	return nil
}

type removeTmpFile struct {
	*os.File
}

// GetAsset looks gets the local image, downloading it if necessary
func (sb *S3Bucket) GetAsset(key string) (io.ReadCloser, error) {
	if sb.cache != nil {
		data, err := sb.cache.GetAsset(key)
		if err == nil {
			return data, nil
		}
	}
	// download the file
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("\n\tError opening tmp file (%s) for s3 download: %s", file.Name(), err.Error())
	}

	_, err = sb.downloader.Download(file, &s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: sb.bucket,
	})
	if err != nil {
		os.Remove(file.Name())
		return nil, fmt.Errorf("\n\tCould not download key (%s) from s3: %s", key, err.Error())
	}

	if sb.cache != nil {
		defer os.Remove(file.Name())
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
