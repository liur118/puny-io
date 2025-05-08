package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &FileStorage{basePath: basePath}, nil
}

func (s *FileStorage) bucketPath(bucket string) string {
	return filepath.Join(s.basePath, bucket)
}

func (s *FileStorage) objectPath(bucket, key string) string {
	return filepath.Join(s.basePath, bucket, key)
}

func (s *FileStorage) ListBuckets() ([]string, error) {
	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory: %w", err)
	}

	var buckets []string
	for _, entry := range entries {
		if entry.IsDir() {
			buckets = append(buckets, entry.Name())
		}
	}

	return buckets, nil
}

func (s *FileStorage) CreateBucket(bucket string) error {
	return os.MkdirAll(s.bucketPath(bucket), 0755)
}

func (s *FileStorage) DeleteBucket(bucket string) error {
	return os.RemoveAll(s.bucketPath(bucket))
}

func (s *FileStorage) ListObjects(bucket, prefix, marker string, maxKeys int) ([]string, error) {
	var objects []string
	bucketPath := s.bucketPath(bucket)

	err := filepath.Walk(bucketPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(bucketPath, path)
		if err != nil {
			return err
		}

		if prefix != "" && !strings.HasPrefix(relPath, prefix) {
			return nil
		}

		if marker != "" && relPath <= marker {
			return nil
		}

		objects = append(objects, relPath)
		if len(objects) >= maxKeys {
			return filepath.SkipAll
		}
		return nil
	})

	return objects, err
}

func (s *FileStorage) PutObject(bucket, key string, content io.Reader, metadata map[string]string) error {
	path := s.objectPath(bucket, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, content)
	return err
}

func (s *FileStorage) GetObject(bucket, key string) (io.ReadCloser, error) {
	return os.Open(s.objectPath(bucket, key))
}

func (s *FileStorage) DeleteObject(bucket, key string) error {
	return os.Remove(s.objectPath(bucket, key))
}

func (s *FileStorage) HeadObject(bucket, key string) (*ObjectMetadata, error) {
	path := s.objectPath(bucket, key)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	return &ObjectMetadata{
		Size:         info.Size(),
		LastModified: info.ModTime(),
		ContentType:  getContentType(path),
		ETag:         hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

func (s *FileStorage) CopyObject(srcBucket, srcKey, dstBucket, dstKey string) error {
	src := s.objectPath(srcBucket, srcKey)
	dst := s.objectPath(dstBucket, dstKey)

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func getContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	default:
		return "application/octet-stream"
	}
}
