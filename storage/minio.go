package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v6"
	"github.com/nilorg/sdk/mime"
	"github.com/nilorg/sdk/storage"
)

var (
	// ErrBucketNameNotIsNil 桶名称不能为空
	ErrBucketNameNotIsNil = errors.New("BucketName not is nil")
)

// MinioStorage minio存储
type MinioStorage struct {
	location    string
	bucketNames []string
	minioClient *minio.Client
}

// NewMinioStorage 创建minio存储
func NewMinioStorage(minioClient *minio.Client, location string, initBucket bool, bucketNames []string) (ms *MinioStorage, err error) {
	ms = &MinioStorage{
		location:    location,
		bucketNames: bucketNames,
		minioClient: minioClient,
	}
	if initBucket {
		err = ms.initBucket()
		if err != nil {
			ms = nil
		}
	}
	return
}

func (ds *MinioStorage) bucketName(parameters ...interface{}) (bucketName string, err error) {
	if len(parameters) < 1 {
		err = errors.New("Please enter bucketName")
		return
	}
	switch parameters[1].(type) {
	case string:
		bucketName = parameters[1].(string)
	default:
		err = errors.New("bucketName parameter type error")
	}
	return
}

func (ds *MinioStorage) filename(parameters ...interface{}) (filename string, err error) {
	if len(parameters) < 0 {
		err = errors.New("Please enter filename")
		return
	}
	switch parameters[0].(type) {
	case string:
		filename = parameters[0].(string)
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
func (ds *MinioStorage) Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	fullName = filename
	options := minio.PutObjectOptions{}
	contextType, contextTypeExist := FromContentTypeContext(ctx)
	if contextTypeExist {
		options.ContentType = contextType
	} else {
		suffix := filepath.Ext(filename)
		suffixContextType, suffixContextTypeExist := mime.Lookup(suffix)
		if !suffixContextTypeExist {
			err = fmt.Errorf("%s unrecognized suffix", suffix)
			return
		}
		options.ContentType = suffixContextType
	}
	md, mdExist := storage.FromIncomingContext(ctx)
	if mdExist {
		options.UserMetadata = md
	}
	_, err = ds.minioClient.PutObjectWithContext(ctx, bucketName, filename, read, -1, options)
	return
}

// Download 下载
func (ds *MinioStorage) Download(ctx context.Context, dist io.Writer, filename string) (results interface{}, err error) {
	var (
		object *minio.Object
	)
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	object, err = ds.minioClient.GetObjectWithContext(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	results, err = object.Stat()
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
