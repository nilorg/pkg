package cache

import (
	"math/rand"
	"time"
)

// RandomTimeSecond 随机多少秒
// 防止缓存雪崩：设置缓存时间为范围随机，解决大量缓存在同一时间失效的问题。
func RandomTimeSecond(min, max int64) time.Duration {
	n := randInt64(min, max)
	return time.Duration(n) * time.Second
}

// RandomTimeMinute 随机多少分钟
// 防止缓存雪崩：设置缓存时间为范围随机，解决大量缓存在同一时间失效的问题。
func RandomTimeMinute(min, max int64) time.Duration {
	n := randInt64(min, max)
	return time.Duration(n) * time.Minute
}

// RandomTimeHour 随机多少小时
// 防止缓存雪崩：设置缓存时间为范围随机，解决大量缓存在同一时间失效的问题。
func RandomTimeHour(min, max int64) time.Duration {
	n := randInt64(min, max)
	return time.Duration(n) * time.Hour
}

func randInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
