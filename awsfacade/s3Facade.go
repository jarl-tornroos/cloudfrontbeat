package awsfacade

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws"
)

// S3Facade handle simple S3 calls
type S3Facade struct {
	client *s3.S3
	downloader *s3manager.Downloader
	Bucket string
}

// NewS3 is a construct function for creating the object with auth as argument
func NewS3(auth *Auth) *S3Facade {
	return &S3Facade{
		client: s3.New(auth.Session, auth.Config),
		downloader: s3manager.NewDownloader(auth.SessionWithConf),
	}
}

// SetBucket is a setter for the S3 bucket to operate on
func (s *S3Facade) SetBucket(bucket string) {
	s.Bucket = bucket
}

// Download object from bucket
func (s *S3Facade) Download(src string, dst string) (int64, error) {
	file, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	numBytes, err := s.downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(src),
		})
	if err != nil {
		return 0, err
	}

	return numBytes, nil
}

// ListFiles from the bucket with prefix as the argument
func (s *S3Facade) ListFiles(prefix string) (*[]string, error) {
	var files []string

	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(prefix),
	}

	err := s.client.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {

		for _, object := range page.Contents {
			files = append(files, *object.Key)
		}

		return true
	})

	return &files, err
}
