package cache

import (
	"errors"
	"time"
)

const (
	defaultNotFoundExpiry = time.Minute
)

var (
	ErrDefaultNotFound = errors.New("cache: not found")
)

type ExpiryFunc func() time.Duration

type Options struct {
	Expiry         ExpiryFunc
	NotFoundExpiry time.Duration
	Prefix         string
	Serializer     Serializer
	ErrNotFound    error
}
type Option func(o *Options)

func newOptions(opts ...Option) Options {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}
	if o.Expiry == nil {
		o.Expiry = func() time.Duration {
			return RandomTimeHour(1, 24)
		}
	}
	if o.NotFoundExpiry <= 0 {
		o.NotFoundExpiry = defaultNotFoundExpiry
	}
	if o.Serializer == nil {
		o.Serializer = JSONSerialize{}
	}
	if o.ErrNotFound == nil {
		o.ErrNotFound = ErrDefaultNotFound
	}
	return o
}

func WithExpiry(expiry ExpiryFunc) Option {
	return func(o *Options) {
		o.Expiry = expiry
	}
}
func WithNotFoundExpiry(expiry time.Duration) Option {
	return func(o *Options) {
		o.NotFoundExpiry = expiry
	}
}
func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}
func WithSerializer(serializer Serializer) Option {
	return func(o *Options) {
		o.Serializer = serializer
	}
}
func WithErrNotFound(err error) Option {
	return func(o *Options) {
		o.ErrNotFound = err
	}
}
