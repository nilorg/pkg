package trace

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

var (
	traceCache *cache.Cache
)

func Init() {
	traceCache = cache.New(5*time.Minute, 10*time.Minute)
}

func NewID() string {
	return uuid.New().String()
}

func NextSpanID(traceID string) string {
	value, found := traceCache.Get(traceID)
	if found {
		spanID := nextSpanID(value.(string))
		traceCache.Set(traceID, spanID, cache.DefaultExpiration)
		return spanID
	}
	return "0"
}

func StartSpanID(traceID, spanID string) string {
	value, found := traceCache.Get(traceID)
	if found {
		spanID := nextSpanID(value.(string))
		traceCache.Set(traceID, spanID, cache.DefaultExpiration)
		return spanID
	} else {
		spanID = startSpanID(spanID)
		traceCache.Set(traceID, spanID, cache.DefaultExpiration)
		return spanID
	}
}

const (
	defaultSpanID = "0"
)

func nextSpanID(spanID string) string {
	if spanID == "" {
		return defaultSpanID
	}
	spanIDLen := len(spanID)
	lastID, _ := strconv.Atoi(spanID[spanIDLen-1:])
	if lastID > 0 {
		spanID = spanID[:spanIDLen-2]
	}
	lastID++
	return fmt.Sprintf("%s.%d", spanID, lastID)
}

func startSpanID(spanID string) string {
	return fmt.Sprintf("%s.%d", spanID, 1)
}
