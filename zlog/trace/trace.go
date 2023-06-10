package trace

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	sdkStrings "github.com/nilorg/sdk/strings"
	"github.com/patrickmn/go-cache"
)

var (
	traceCache *cache.Cache
	once       sync.Once
)

func NewID() string {
	return uuid.New().String()
}

func getSingleton() *cache.Cache {
	once.Do(func() {
		traceCache = cache.New(5*time.Minute, 10*time.Minute)
	})
	return traceCache
}

func NextSpanID(traceID string) string {
	value, found := getSingleton().Get(traceID)
	if found {
		spanID := nextSpanID(value.(string))
		getSingleton().Set(traceID, spanID, cache.DefaultExpiration)
		return spanID
	}
	return "0"
}

func StartSpanID(traceID, spanID string) string {
	value, found := getSingleton().Get(traceID)
	if found {
		spanID := nextSpanID(value.(string))
		getSingleton().Set(traceID, spanID, cache.DefaultExpiration)
		return spanID
	} else {
		spanID = startSpanID(spanID)
		getSingleton().Set(traceID, spanID, cache.DefaultExpiration)
		return spanID
	}
}

const (
	defaultSpanID = "0"
)

func nextSpanID(spanID string) string {
	spanIDs := sdkStrings.Split(spanID, ".")
	if len(spanIDs) == 0 {
		return defaultSpanID
	}
	spanIDNumbers := make([]int, len(spanIDs))
	for i, v := range spanIDs {
		spanIDNumbers[i], _ = strconv.Atoi(v)
	}
	spanIDNumbers[len(spanIDNumbers)-1] = spanIDNumbers[len(spanIDNumbers)-1] + 1

	newSpanIDs := make([]string, len(spanIDNumbers))
	for i, v := range spanIDNumbers {
		newSpanIDs[i] = strconv.Itoa(v)
	}
	return strings.Join(newSpanIDs, ".")
}

func startSpanID(spanID string) string {
	return fmt.Sprintf("%s.%d", spanID, 1)
}
