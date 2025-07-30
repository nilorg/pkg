package storage

import (
	"context"
	"fmt"

	"github.com/nilorg/sdk/storage"
)

// Example_newRangeAPI 展示新的Range API设计
func Example_newRangeAPI() {
	ctx := context.Background()

	// 1. 读取前1024个字节 (0-1023)
	ctx1 := storage.NewRangeRequestContext(ctx, 0, 1023)
	if rangeReq, ok := storage.FromRangeRequestContext(ctx1); ok {
		fmt.Printf("读取前1024个字节: bytes=%d-%d, 长度=%d\n",
			rangeReq.Start, rangeReq.End, rangeReq.Length())
	}

	// 2. 从第1024个字节开始读取到文件末尾
	ctx2 := storage.NewRangeRequestContext(ctx, 1024, 0)
	if rangeReq, ok := storage.FromRangeRequestContext(ctx2); ok {
		fmt.Printf("从第1024个字节到末尾: bytes=%d-, 到末尾=%t\n",
			rangeReq.Start, rangeReq.IsToEnd())
	}

	// 3. 读取整个文件
	ctx3 := storage.NewRangeRequestContext(ctx, 0, 0)
	if rangeReq, ok := storage.FromRangeRequestContext(ctx3); ok {
		fmt.Printf("读取整个文件: bytes=%d-, 到末尾=%t\n",
			rangeReq.Start, rangeReq.IsToEnd())
	}

	// 4. 使用负数也表示到文件末尾
	ctx4 := storage.NewRangeRequestContext(ctx, 500, -1)
	if rangeReq, ok := storage.FromRangeRequestContext(ctx4); ok {
		fmt.Printf("从第500个字节到末尾: bytes=%d-, 到末尾=%t\n",
			rangeReq.Start, rangeReq.IsToEnd())
	}

	// 5. 读取单个字节
	ctx5 := storage.NewRangeRequestContext(ctx, 100, 100)
	if rangeReq, ok := storage.FromRangeRequestContext(ctx5); ok {
		fmt.Printf("读取单个字节: bytes=%d-%d, 长度=%d\n",
			rangeReq.Start, rangeReq.End, rangeReq.Length())
	}

	// Output:
	// 读取前1024个字节: bytes=0-1023, 长度=1024
	// 从第1024个字节到末尾: bytes=1024-, 到末尾=true
	// 读取整个文件: bytes=0-, 到末尾=true
	// 从第500个字节到末尾: bytes=500-, 到末尾=true
	// 读取单个字节: bytes=100-100, 长度=1
}

// Example_rangeUseCases 展示Range请求的常见用例
func Example_rangeUseCases() {
	ctx := context.Background()

	fmt.Println("=== Range请求常见用例 ===")

	// 用例1：下载文件头信息（前512字节）
	_ = storage.NewRangeRequestContext(ctx, 0, 511)
	fmt.Println("1. 下载文件头: start=0, end=511 (512字节)")

	// 用例2：断点续传 - 已下载2MB，继续下载
	resumeStart := int64(2 * 1024 * 1024) // 2MB
	_ = storage.NewRangeRequestContext(ctx, resumeStart, 0)
	fmt.Printf("2. 断点续传: start=%d, end=0 (到文件末尾)\n", resumeStart)

	// 用例3：分块下载 - 下载第2个1MB块
	chunkSize := int64(1024 * 1024) // 1MB
	chunkStart := chunkSize         // 第2个块从1MB开始
	chunkEnd := chunkStart + chunkSize - 1
	_ = storage.NewRangeRequestContext(ctx, chunkStart, chunkEnd)
	fmt.Printf("3. 分块下载: start=%d, end=%d (第2个1MB块)\n", chunkStart, chunkEnd)

	// 用例4：下载文件尾部信息（最后1KB）
	// 注意：这个需要知道文件大小，这里假设文件大小为10MB
	fileSize := int64(10 * 1024 * 1024)
	tailStart := fileSize - 1024
	_ = storage.NewRangeRequestContext(ctx, tailStart, fileSize-1)
	fmt.Printf("4. 下载文件尾部: start=%d, end=%d (最后1KB)\n", tailStart, fileSize-1)

	// Output:
	// === Range请求常见用例 ===
	// 1. 下载文件头: start=0, end=511 (512字节)
	// 2. 断点续传: start=2097152, end=0 (到文件末尾)
	// 3. 分块下载: start=1048576, end=2097151 (第2个1MB块)
	// 4. 下载文件尾部: start=10484736, end=10485759 (最后1KB)
}
