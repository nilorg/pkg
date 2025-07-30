package storage

import (
	"testing"

	"github.com/nilorg/sdk/storage"
	"github.com/stretchr/testify/assert"
)

// TestRangeRequestValidation 测试Range请求的各种边界情况
func TestRangeRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		start     int64
		end       int64
		wantValid bool
		wantToEnd bool
		wantLen   int64
		desc      string
	}{
		{
			name:      "读取前100个字节",
			start:     0,
			end:       99,
			wantValid: true,
			wantToEnd: false,
			wantLen:   100,
			desc:      "bytes=0-99，读取100个字节",
		},
		{
			name:      "读取单个字节",
			start:     50,
			end:       50,
			wantValid: true,
			wantToEnd: false,
			wantLen:   1,
			desc:      "bytes=50-50，读取1个字节",
		},
		{
			name:      "从第101个字节到文件末尾（end=0）",
			start:     100,
			end:       0,
			wantValid: true,
			wantToEnd: true,
			wantLen:   -1,
			desc:      "bytes=100-，读取到文件末尾",
		},
		{
			name:      "从第101个字节到文件末尾（end=-1）",
			start:     100,
			end:       -1,
			wantValid: true,
			wantToEnd: true,
			wantLen:   -1,
			desc:      "bytes=100-，读取到文件末尾",
		},
		{
			name:      "从第101个字节到文件末尾（end=-100）",
			start:     100,
			end:       -100,
			wantValid: true,
			wantToEnd: true,
			wantLen:   -1,
			desc:      "bytes=100-，读取到文件末尾",
		},
		{
			name:      "读取中间的一段",
			start:     100,
			end:       199,
			wantValid: true,
			wantToEnd: false,
			wantLen:   100,
			desc:      "bytes=100-199，读取100个字节",
		},
		{
			name:      "从文件开始到末尾",
			start:     0,
			end:       0,
			wantValid: true,
			wantToEnd: true,
			wantLen:   -1,
			desc:      "bytes=0-，读取整个文件",
		},
		{
			name:      "无效：start为负数",
			start:     -1,
			end:       100,
			wantValid: false,
			wantToEnd: false,
			wantLen:   0,
			desc:      "start不能为负数",
		},
		{
			name:      "无效：end小于start且为正数",
			start:     100,
			end:       50,
			wantValid: false,
			wantToEnd: false,
			wantLen:   0,
			desc:      "end不能小于start（当end为正数时）",
		},
		{
			name:      "无效：start为负数且end为负数",
			start:     -5,
			end:       -1,
			wantValid: false,
			wantToEnd: false,
			wantLen:   0,
			desc:      "start不能为负数",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rangeReq := &storage.RangeRequest{
				Start: tt.start,
				End:   tt.end,
			}

			assert.Equal(t, tt.wantValid, rangeReq.IsValid(), "IsValid() = %v, want %v (%s)", rangeReq.IsValid(), tt.wantValid, tt.desc)

			if tt.wantValid {
				assert.Equal(t, tt.wantToEnd, rangeReq.IsToEnd(), "IsToEnd() = %v, want %v (%s)", rangeReq.IsToEnd(), tt.wantToEnd, tt.desc)
				assert.Equal(t, tt.wantLen, rangeReq.Length(), "Length() = %v, want %v (%s)", rangeReq.Length(), tt.wantLen, tt.desc)
			}
		})
	}
}
