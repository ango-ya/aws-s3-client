package s3

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	ErrUnsupported = errors.New("unsupported")
)

type S3Client struct {
	cfg    aws.Config
	client *s3.Client

	logger zerolog.Logger

	uploader            *manager.Uploader
	uploaderPartSize    int64
	uploaderConcurrency int

	timeout time.Duration
}

func NewS3Client(ctx context.Context, confOpts []func(*config.LoadOptions) error, opts ...Option) (c S3Client, err error) {
	c.logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	c.timeout = DefaultTimeout
	c.uploaderPartSize = DefaultUploaderPartSize
	c.uploaderConcurrency = DefaultUploaderConcurrency

	for i := range opts {
		opts[i].Apply(&c)
	}

	r := retry.NewStandard(func(o *retry.StandardOptions) {
		o.MaxAttempts = retry.DefaultMaxAttempts // 3
	})

	confOpts = append(confOpts, config.WithRetryer(func() aws.Retryer { return r }))

	if c.cfg, err = config.LoadDefaultConfig(ctx, confOpts...); err != nil {
		err = errors.Wrap(err, "failed to load aws config")
		return
	}

	c.client = s3.NewFromConfig(c.cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	c.uploader = manager.NewUploader(c.client, func(u *manager.Uploader) {
		u.PartSize = c.uploaderPartSize
		u.Concurrency = c.uploaderConcurrency
		u.LeavePartsOnError = false
	})

	return
}

func (c *S3Client) Put(ctx context.Context, bucket, key string, body io.Reader) (location string, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	output, err := c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
		// ContentType:
	})
	if err != nil {
		err = errors.Wrap(err, "failed to uplad in s3")
		return
	}

	location = output.Location

	return
}

func (c *S3Client) Delete(ctx context.Context, bucket, key string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err = c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return
}
