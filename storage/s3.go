package storage

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nilorg/sdk/mime"
	"github.com/nilorg/sdk/storage"
)

var _ storage.Storager = (*AwsS3Storage)(nil)

// AwsS3Storage awsS3存储
type AwsS3Storage struct {
	location                    string
	bucketNames                 []string
	s3Client                    *s3.S3
	CheckAndCreateBucketEnabled bool
}

// NewAwsS3Storage 创建AWS S3存储
func NewAwsS3Storage(s3Client *s3.S3, location string, initBucket bool, bucketNames []string) (ms *AwsS3Storage, err error) {
	ms = &AwsS3Storage{
		location:                    location,
		bucketNames:                 bucketNames,
		s3Client:                    s3Client,
		CheckAndCreateBucketEnabled: false,
	}
	if initBucket {
		err = ms.initBucket()
		if err != nil {
			ms = nil
		}
	}
	return
}

// initBucket 初始化桶
func (ds *AwsS3Storage) initBucket() (err error) {
	for _, bucketName := range ds.bucketNames {
		err = ds.CheckAndCreateBucket(bucketName)
		if err != nil {
			return
		}
	}
	return
}

// CheckAndCreateBucket 检查并创建桶
func (ds *AwsS3Storage) CheckAndCreateBucket(bucketName string) (err error) {
	// 检查存储桶是否已经存在
	_, err = ds.s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		// 桶已存在
		return nil
	}

	// 创建桶
	_, err = ds.s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(ds.location),
		},
	})
	return
}

// Upload 上传
func (ds *AwsS3Storage) Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	if ds.CheckAndCreateBucketEnabled {
		err = ds.CheckAndCreateBucket(bucketName)
		if err != nil {
			return
		}
	}
	if rename, ok := storage.FromRenameContext(ctx); ok {
		filename = rename(filename)
	}
	fullName = filename

	contentType, contentTypeExist := FromContentTypeContext(ctx)
	if !contentTypeExist {
		contentType, err = mime.DetectContentType(filename)
		if err != nil {
			return
		}
	}

	// 准备上传参数
	params := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		Body:        aws.ReadSeekCloser(read),
		ContentType: aws.String(contentType),
	}

	// 添加自定义元数据
	if md, mdExist := storage.FromIncomingContext(ctx); mdExist {
		metadata := make(map[string]*string)
		for k, v := range md {
			metadata[k] = aws.String(v)
		}
		params.Metadata = metadata
	}

	_, err = ds.s3Client.PutObjectWithContext(ctx, params)
	return
}

// Download 下载
func (ds *AwsS3Storage) Download(ctx context.Context, dist io.Writer, filename string) (info storage.DownloadFileInfoer, err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}

	// 获取对象
	result, err := ds.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return
	}
	defer result.Body.Close()

	// 准备元数据
	md := storage.Metadata{
		"Content-Type": *result.ContentType,
	}

	// 添加用户元数据
	for k, v := range result.Metadata {
		if v != nil {
			md.Set(k, *v)
		}
	}

	var (
		downloadFilename      string
		downloadFilenameExist bool
	)
	if downloadFilename, downloadFilenameExist = storage.FromDownloadFilenameContext(ctx); !downloadFilenameExist {
		downloadFilename = filepath.Base(filename)
	}

	info = &downloadFileInfo{
		filename: downloadFilename,
		size:     *result.ContentLength,
		metadata: md,
	}

	if downloadBefore, downloadBeforeExist := storage.FromDownloadBeforeContext(ctx); downloadBeforeExist {
		downloadBefore(info)
	}

	_, err = io.Copy(dist, result.Body)
	if err != nil {
		info = nil
		return
	}
	return
}

// Remove 删除
func (ds *AwsS3Storage) Remove(ctx context.Context, filename string) (err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}

	_, err = ds.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	return
}

// Exist 是否存在
func (ds *AwsS3Storage) Exist(ctx context.Context, filename string) (exist bool, err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}

	_, err = ds.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		// 检查错误类型
		// 如果是对象不存在的错误，不返回错误
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NotFound" || aerr.Code() == "NoSuchKey" {
				err = nil
				return false, nil
			}
		}
		return false, err
	}

	exist = true
	return
}
