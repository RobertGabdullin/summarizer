package repository

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func NewS3Client(bucket string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	s3Client := &S3Client{
		client: client,
		bucket: bucket,
	}

	return s3Client, nil
}

func (c *S3Client) StoreMeeting(ctx context.Context, id string, file io.Reader) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(id),
		Body:   file,
	})

	if err != nil {
		return err
	}

	return nil
}
