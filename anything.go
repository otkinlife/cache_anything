package cache_anything

import (
	"time"
)

type Handler func(params interface{}, results interface{}) error

func CacheAnything(key string, handler Handler, params interface{}, results interface{}, expire time.Duration) error {
	// 尝试从cache中取结果
	if err := GlobalCache.LoadDataFromJson(key, results); err == nil && results != nil {
		return nil
	}
	err := handler(params, results)
	if err != nil {
		return err
	}
	// 异步写缓存写入缓存
	return GlobalCache.SetDataWithJsonWithExpiration(key, results, expire)
}
