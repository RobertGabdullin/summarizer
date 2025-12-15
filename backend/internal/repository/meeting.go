package repository

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func NewS3Client(region, bucket, accessKeyID, secretAccessKey, endpoint string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKeyID,
				secretAccessKey,
				"",
			),
		))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	s3Client := &S3Client{
		client: client,
		bucket: bucket,
	}

	err = s3Client.ensureBucket(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Не получилось проверить наличие бакета в S3: %w", err)
	}

	return s3Client, nil
}

func (c *S3Client) ensureBucket(ctx context.Context) error {
	_, err := c.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &c.bucket})
	if err == nil {
		return nil
	}

	_, err = c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &c.bucket,
	})
	return err
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
