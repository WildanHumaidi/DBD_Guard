package storage

import (
	"context"
	"io"
	"time"
)

type MockStorageService struct {
	UploadFunc          func(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error)
	DownloadFunc        func(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	DeleteFunc          func(ctx context.Context, bucketName, objectName string) error
	ExistsFunc          func(ctx context.Context, bucketName, objectName string) (bool, error)
	GetPresignedURLFunc func(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
}

func (m *MockStorageService) Upload(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(ctx, bucketName, objectName, reader, objectSize, contentType)
	}
	return objectName, nil
}

func (m *MockStorageService) Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	if m.DownloadFunc != nil {
		return m.DownloadFunc(ctx, bucketName, objectName)
	}
	return nil, nil
}

func (m *MockStorageService) Delete(ctx context.Context, bucketName, objectName string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, bucketName, objectName)
	}
	return nil
}

func (m *MockStorageService) Exists(ctx context.Context, bucketName, objectName string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, bucketName, objectName)
	}
	return true, nil
}

func (m *MockStorageService) GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	if m.GetPresignedURLFunc != nil {
		return m.GetPresignedURLFunc(ctx, bucketName, objectName, expiry)
	}
	return "http://mock-minio/" + bucketName + "/" + objectName, nil
}
