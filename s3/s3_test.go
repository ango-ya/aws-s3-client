package s3

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/require"
)

const (
	TestRegion     = "us-east-1"
	TestProfile    = "sto-dev"
	TestImagePath  = "./sample.jpeg"
	TestBucketName = "tak-sandbox"
	TestKey        = "sample.jpeg"
)

func TestAll(t *testing.T) {
	file, err := os.Open(TestImagePath)
	require.NoError(t, err)
	defer file.Close()

	var (
		ctx      = context.Background()
		confOpts = []func(*config.LoadOptions) error{
			config.WithDefaultRegion(TestRegion),
			config.WithSharedConfigProfile(TestProfile),
		}
	)

	c, err := NewS3Client(ctx, confOpts, WithTimeout(5*time.Second))
	require.NoError(t, err)

	path, err := c.Put(ctx, TestBucketName, TestKey, file)
	require.NoError(t, err)

	require.Equal(t, "https://s3.us-east-1.amazonaws.com/tak-sandbox/sample.jpeg", path)

	err = c.Delete(ctx, TestBucketName, TestKey)
	require.NoError(t, err)
}
