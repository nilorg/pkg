package viperx

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetViperSliceMap(conf *viper.Viper, key string) (results []map[string]interface{}) {
	values := conf.Get(key).([]interface{})
	for _, v := range values {
		value := v.(map[interface{}]interface{})
		newValue := make(map[string]interface{})
		for vk, vv := range value {
			newValue[fmt.Sprint(vk)] = vv
		}
		results = append(results, newValue)
	}
	return
}

func GetViperSliceMapKeyValues(conf *viper.Viper, key, mkey string) (mkeyValues []interface{}) {
	values := GetViperSliceMap(conf, key)
	for _, v := range values {
		if v, ok := v[mkey]; ok {
			mkeyValues = append(mkeyValues, v)
		}
	}
	return
}
