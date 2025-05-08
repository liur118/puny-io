package storage

import (
	"io"
	"time"
)

type ObjectMetadata struct {
	Size         int64
	LastModified time.Time
	ContentType  string
	ETag         string
}

type Storage interface {
	// 桶操作
	ListBuckets() ([]string, error)
	CreateBucket(bucket string) error
	DeleteBucket(bucket string) error
	ListObjects(bucket string, prefix string, marker string, maxKeys int) ([]string, error)

	// 对象操作
	PutObject(bucket, key string, content io.Reader, metadata map[string]string) error
	GetObject(bucket, key string) (io.ReadCloser, error)
	DeleteObject(bucket, key string) error
	HeadObject(bucket, key string) (*ObjectMetadata, error)
	CopyObject(srcBucket, srcKey, dstBucket, dstKey string) error
}
