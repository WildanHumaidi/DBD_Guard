package storage

import (
	"context"
	"io"
	"time"
)

type StorageService interface {
	Upload(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error)
	Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucketName, objectName string) error
	Exists(ctx context.Context, bucketName, objectName string) (bool, error)
	GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
}
