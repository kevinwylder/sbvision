package cdn

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const bucketName = "skateboardvision.net"
const distributionID = "E37VGOJSB08WSG"

// Uploader is able to add things to the website
type Uploader struct {
	s3 *s3manager.Uploader
	cf *cloudfront.CloudFront
}

// NewUploader creates an uploader to add to the CDN
func NewUploader(sess *session.Session) *Uploader {
	return &Uploader{
		s3: s3manager.NewUploader(sess),
		cf: cloudfront.New(sess),
	}
}

// Add adds this data to the given path
func (u *Uploader) Add(data io.Reader, path string) error {
	_, err := u.s3.Upload(&s3manager.UploadInput{
		Body:   data,
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})
	return err
}

// Invalidate makes the changes made by Add live
func (u *Uploader) Invalidate(path string) error {
	_, err := u.cf.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().Format(time.RFC1123Z)),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items:    []*string{aws.String(path)},
			},
		},
	})
	return err
}
