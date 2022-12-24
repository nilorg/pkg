package cache

import (
	jsoniter "github.com/json-iterator/go"
)

var jsonStandard = jsoniter.ConfigCompatibleWithStandardLibrary

// Serializer 序列化器
type Serializer interface {
	Unmarshal(data []byte, msg interface{}) (err error)
	Marshal(msg interface{}) (data []byte, err error)
	ContentType() string
}

// JSONSerialize 高性能json序列化
type JSONSerialize struct {
}

// Unmarshal ...
func (JSONSerialize) Unmarshal(data []byte, msg interface{}) (err error) {
	return jsonStandard.Unmarshal(data, msg)
}

// Marshal ...
func (JSONSerialize) Marshal(msg interface{}) (data []byte, err error) {
	return jsonStandard.Marshal(msg)
}

// ContentType ...
func (JSONSerialize) ContentType() string {
	return "application/json"
}
