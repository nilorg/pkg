package storage

import (
	"context"
	"testing"

	"github.com/nilorg/sdk/storage"
	"github.com/stretchr/testify/assert"
)

// TestRangeRequestContext 测试Range请求的Context操作
func TestRangeRequestContext(t *testing.T) {
	ctx := context.Background()

	// 测试创建有效的Range请求Context
	rangeCtx := storage.NewRangeRequestContext(ctx, 100, 200)
	rangeReq, ok := storage.FromRangeRequestContext(rangeCtx)

	assert.True(t, ok)
	assert.NotNil(t, rangeReq)
	assert.Equal(t, int64(100), rangeReq.Start)
	assert.Equal(t, int64(200), rangeReq.End)
	assert.True(t, rangeReq.IsValid())
	assert.False(t, rangeReq.IsToEnd())
	assert.Equal(t, int64(101), rangeReq.Length()) // 200-100+1

	// 测试从普通context中获取Range请求（应该失败）
	_, ok = storage.FromRangeRequestContext(ctx)
	assert.False(t, ok)

	// 测试创建到文件末尾的Range请求
	rangeToEndCtx := storage.NewRangeRequestContext(ctx, 100, 0)
	rangeToEndReq, ok := storage.FromRangeRequestContext(rangeToEndCtx)

	assert.True(t, ok)
	assert.NotNil(t, rangeToEndReq)
	assert.Equal(t, int64(100), rangeToEndReq.Start)
	assert.Equal(t, int64(0), rangeToEndReq.End)
	assert.True(t, rangeToEndReq.IsValid())
	assert.True(t, rangeToEndReq.IsToEnd())
	assert.Equal(t, int64(-1), rangeToEndReq.Length())

	// 测试边界情况：End等于0的情况（表示到文件末尾）
	rangeFromZeroCtx := storage.NewRangeRequestContext(ctx, 0, 0)
	rangeFromZeroReq, ok := storage.FromRangeRequestContext(rangeFromZeroCtx)

	assert.True(t, ok)
	assert.NotNil(t, rangeFromZeroReq)
	assert.Equal(t, int64(0), rangeFromZeroReq.Start)
	assert.Equal(t, int64(0), rangeFromZeroReq.End)
	assert.True(t, rangeFromZeroReq.IsValid())
	assert.True(t, rangeFromZeroReq.IsToEnd())
	assert.Equal(t, int64(-1), rangeFromZeroReq.Length())

	// 测试无效的Range请求：Start > End
	invalidCtx := storage.NewRangeRequestContext(ctx, 200, 100)
	_, ok = storage.FromRangeRequestContext(invalidCtx)
	assert.False(t, ok) // 无效的Range请求不会被存储

	// 测试无效的Range请求：Start < 0
	invalidStartCtx := storage.NewRangeRequestContext(ctx, -5, 100)
	_, ok = storage.FromRangeRequestContext(invalidStartCtx)
	assert.False(t, ok) // 无效的Range请求不会被存储
}

// Example_rangeRequest 示例：如何在实际代码中使用Range请求进行断点续传
func Example_rangeRequest() {
	ctx := context.Background()

	// 假设我们有一个Storage实例和要下载的文件
	// var storage storage.Storager
	// filename := "large-file.zip"
	// bucketName := "my-bucket"

	// 1. 设置bucket名称
	ctx = NewBucketNameContext(ctx, "my-bucket")

	// 2. 设置Range请求 - 从字节1024开始下载到字节2047
	ctx = storage.NewRangeRequestContext(ctx, 1024, 2047)

	// 3. 执行下载（这将使用Range请求）
	// var buffer bytes.Buffer
	// info, err := storage.Download(ctx, &buffer, filename)
	// if err != nil {
	//     // 处理错误
	//     return
	// }

	// 4. 检查下载的内容（buffer.Bytes()包含指定范围的数据）
	// downloadedData := buffer.Bytes()
	// fmt.Printf("Downloaded %d bytes\n", len(downloadedData))
}

// Example_resumeDownload 示例：如何实现断点续传下载
func Example_resumeDownload() {
	ctx := context.Background()

	// 假设我们已经下载了一部分文件，现在想继续下载
	alreadyDownloadedBytes := int64(1024)

	// 1. 设置bucket名称
	ctx = NewBucketNameContext(ctx, "my-bucket")

	// 2. 从已下载的位置继续下载到文件末尾
	ctx = storage.NewRangeRequestContext(ctx, alreadyDownloadedBytes, 0)

	// 3. 执行断点续传下载
	// var buffer bytes.Buffer
	// info, err := storage.Download(ctx, &buffer, "large-file.zip")
	// if err != nil {
	//     // 处理错误
	//     return
	// }

	// 4. 将新下载的数据追加到已有文件
	// newData := buffer.Bytes()
	// 将newData追加到本地文件...
}
