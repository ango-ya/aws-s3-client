package s3

import (
	"time"

	"github.com/rs/zerolog"
)

var (
	DefaultTimeout             = 10 * time.Second
	DefaultUploaderPartSize    = int64(8 * 1024 * 1024) // 8 MB
	DefaultUploaderConcurrency = 2
)

type Option interface {
	Apply(*S3Client) error
}

type Timeout time.Duration

func (o Timeout) Apply(c *S3Client) error {
	c.timeout = time.Duration(o)
	return nil
}
func WithTimeout(timeout time.Duration) Timeout {
	return Timeout(timeout)
}

type Logger zerolog.Logger

func (o Logger) Apply(c *S3Client) error {
	c.logger = zerolog.Logger(o)
	return nil
}
func WithLogger(logger zerolog.Logger) Logger {
	return Logger(logger)
}

type UploaderPartSize int64

func (o UploaderPartSize) Apply(c *S3Client) error {
	c.uploaderPartSize = int64(o)
	return nil
}
func WithUploaderPartSize(uploaderPartSize int64) UploaderPartSize {
	return UploaderPartSize(uploaderPartSize)
}

type UploaderConcurrency int

func (o UploaderConcurrency) Apply(c *S3Client) error {
	c.uploaderConcurrency = int(o)
	return nil
}
func WithUploaderConcurrency(uploaderConcurrency int) UploaderConcurrency {
	return UploaderConcurrency(uploaderConcurrency)
}
