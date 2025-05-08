package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liur/puny-io/internal/storage"
)

type OSSHandler struct {
	storage storage.Storage
	host    string
}

func NewOSSHandler(storage storage.Storage, host string) *OSSHandler {
	return &OSSHandler{
		storage: storage,
		host:    host,
	}
}

// 桶操作
func (h *OSSHandler) ListBuckets(c *gin.Context) {
	buckets, err := h.storage.ListBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"buckets": buckets})
}

func (h *OSSHandler) CreateBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	if err := h.storage.CreateBucket(bucket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OSSHandler) DeleteBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	if err := h.storage.DeleteBucket(bucket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OSSHandler) ListObjects(c *gin.Context) {
	bucket := c.Param("bucket")
	prefix := c.Query("prefix")
	marker := c.Query("marker")
	maxKeys, _ := strconv.Atoi(c.DefaultQuery("max-keys", "1000"))

	objects, err := h.storage.ListObjects(bucket, prefix, marker, maxKeys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"objects": objects})
}

// 对象操作
func (h *OSSHandler) PutObject(c *gin.Context) {
	bucket := c.Param("bucket")
	key := c.Param("key")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	metadata := make(map[string]string)
	for k, v := range c.Request.Header {
		if strings.HasPrefix(k, "X-Amz-Meta-") {
			metadata[k] = v[0]
		}
	}

	if err := h.storage.PutObject(bucket, key, file, metadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *OSSHandler) GetObject(c *gin.Context) {
	log.Println("GetObject", c.Request.URL.Path)
	bucket := c.Param("bucket")
	key := c.Param("key")

	object, err := h.storage.GetObject(bucket, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer object.Close()

	metadata, err := h.storage.HeadObject(bucket, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", metadata.ContentType)
	c.Header("ETag", metadata.ETag)
	c.Header("Last-Modified", metadata.LastModified.Format(http.TimeFormat))
	c.Header("Content-Length", strconv.FormatInt(metadata.Size, 10))

	io.Copy(c.Writer, object)
}

func (h *OSSHandler) DeleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	key := c.Param("key")

	if err := h.storage.DeleteObject(bucket, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *OSSHandler) HeadObject(c *gin.Context) {
	bucket := c.Param("bucket")
	key := c.Param("key")

	metadata, err := h.storage.HeadObject(bucket, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", metadata.ContentType)
	c.Header("ETag", metadata.ETag)
	c.Header("Last-Modified", metadata.LastModified.Format(http.TimeFormat))
	c.Header("Content-Length", strconv.FormatInt(metadata.Size, 10))
	c.Status(http.StatusOK)
}

func (h *OSSHandler) CopyObject(c *gin.Context) {
	bucket := c.Param("bucket")
	key := c.Param("key")
	srcBucket := c.GetHeader("X-Amz-Copy-Source-Bucket")
	srcKey := c.GetHeader("X-Amz-Copy-Source-Key")

	if srcBucket == "" || srcKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing source bucket or key"})
		return
	}

	if err := h.storage.CopyObject(srcBucket, srcKey, bucket, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *OSSHandler) GetObjectURL(c *gin.Context) {
	bucket := c.Param("bucket")
	key := c.Param("key")

	// 检查对象是否存在
	_, err := h.storage.HeadObject(bucket, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成访问链接
	url := fmt.Sprintf("%s/oss/%s/%s", h.host, bucket, key)
	c.JSON(http.StatusOK, gin.H{"url": url})
}
