package storage

import (
	"context"
	"errors"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v6"
	"github.com/nilorg/sdk/mime"
)

// MinioStorage minio存储
type MinioStorage struct {
	location    string
	bucketNames []string
	minioClient *minio.Client
}

// NewMinioStorage 创建minio存储
func NewMinioStorage(minioClient *minio.Client, location string, bucketNames []string) (ms *MinioStorage, err error) {
	ms = &MinioStorage{
		location:    location,
		bucketNames: bucketNames,
		minioClient: minioClient,
	}
	err = ms.initBucket()
	if err != nil {
		ms = nil
	}
	return
}

func (ds *MinioStorage) bucketName(parameters ...interface{}) (bucketName string, err error) {
	if len(parameters) < 0 {
		err = errors.New("Please enter bucketName")
		return
	}
	switch parameters[0].(type) {
	case string:
		bucketName = parameters[0].(string)
	default:
		err = errors.New("bucketName parameter type error")
	}
	return
}

func (ds *MinioStorage) filename(parameters ...interface{}) (filename string, err error) {
	if len(parameters) < 1 {
		err = errors.New("Please enter filename")
		return
	}
	switch parameters[1].(type) {
	case string:
		filename = parameters[1].(string)
	default:
		err = errors.New("filename parameter type error")
	}
	return
}

// initBucket 初始化桶
func (ds *MinioStorage) initBucket() (err error) {
	for _, bucketName := range ds.bucketNames {
		var exists bool
		// 检查存储桶是否已经存在。
		exists, err = ds.minioClient.BucketExists(bucketName)
		if err != nil {
			return
		}
		if exists {
			continue
		}
		// 创建桶
		err = ds.minioClient.MakeBucket(bucketName, ds.location)
		if err != nil {
			return
		}
	}
	return
}

// Upload 上传
func (ds *MinioStorage) Upload(ctx context.Context, read io.Reader, parameters ...interface{}) (filename string, err error) {
	var (
		bucketName string
	)
	bucketName, err = ds.bucketName(parameters...)
	if err != nil {
		return
	}
	filename, err = ds.filename(parameters...)
	if err != nil {
		return
	}
	suffix := filepath.Ext(filename)
	contextType, exist := mime.Lookup(suffix)
	if !exist {
		err = errors.New("unrecognized suffix")
		return
	}
	_, err = ds.minioClient.PutObjectWithContext(ctx, bucketName, filename, read, -1, minio.PutObjectOptions{
		ContentType: contextType,
	})
	return
}

// Download 下载
func (ds *MinioStorage) Download(ctx context.Context, dist io.Writer, parameters ...interface{}) (err error) {
	var (
		bucketName string
		filename   string
		object     *minio.Object
	)
	bucketName, err = ds.bucketName(parameters...)
	if err != nil {
		return
	}
	filename, err = ds.filename(parameters...)
	if err != nil {
		return
	}
	object, err = ds.minioClient.GetObjectWithContext(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	_, err = io.Copy(dist, object)
	return
}

// Remove 删除
func (ds *MinioStorage) Remove(_ context.Context, fullPath string, parameters ...interface{}) (err error) {
	var (
		bucketName string
	)
	bucketName, err = ds.bucketName(parameters...)
	if err != nil {
		return
	}
	err = ds.minioClient.RemoveObject(bucketName, fullPath)
	return
}
