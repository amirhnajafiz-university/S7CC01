package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3 struct {
	Bucket  string
	Session *session.Session
}

// NewSession
// will create a new s3 session.
func NewSession(cfg Config) (*S3, error) {
	var s3 S3

	newSession, err := session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.Region),
			Credentials: credentials.NewStaticCredentials(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		},
	)
	if err != nil {
		return nil, err
	}

	s3.Session = newSession
	s3.Bucket = cfg.Bucket

	return &s3, nil
}
