package cache_anything

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheAnything(t *testing.T) {
	// 初始化缓存和监控器
	err := Init(Config{
		MaxSize:  100,
		PlanTime: "18:30:00",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	var res string
	err = CacheAnything("test", CaseHandler, "hello world", &res, 10*time.Second)
	if err != nil {
		t.Fatalf("CacheAnything failed: %v", err)
	}
	err = CacheAnything("test", CaseHandler, "hello world", &res, 10*time.Second)
	if err != nil {
		t.Fatalf("CacheAnything failed: %v", err)
	}
	time.Sleep(11 * time.Second)
}

func CaseHandler(params interface{}, results interface{}) error {
	fmt.Println("do something")
	// 将 results 参数转换为指针类型
	ptr := results.(*string)
	// 使用解引用指针的方式将 params 的值写入到指针指向的位置
	*ptr = params.(string)
	return nil
}

func TestCache(t *testing.T) {
	// 初始化缓存和监控器
	err := Init(Config{
		MaxSize:  100,
		PlanTime: "18:30:00",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 写入缓存
	for i := 0; i < 10; i++ {
		GlobalCache.SetString(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// 读取缓存
	for i := 0; i < 10; i++ {
		value, err := GlobalCache.GetString(fmt.Sprintf("key%d", i))
		if err != nil {
			t.Fatalf("GetString failed: %v", err)
		}
		if value != fmt.Sprintf("value%d", i) {
			t.Fatalf("Expected value%d, got %s", i, value)
		}
	}

	// 写入超出限制的缓存，触发LRU淘汰
	for i := 10; i < 20; i++ {
		GlobalCache.SetString(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// 检查被淘汰的缓存项
	for i := 0; i < 5; i++ {
		value, err := GlobalCache.GetString(fmt.Sprintf("key%d", i))
		if err == nil && value != "" {
			t.Fatalf("Expected key%d to be evicted, but got %s", i, value)
		}
	}

	// 检查未被淘汰的缓存项
	for i := 6; i < 10; i++ {
		value, err := GlobalCache.GetString(fmt.Sprintf("key%d", i))
		if err != nil {
			t.Fatalf("GetString failed: %v", err)
		}
		if value != fmt.Sprintf("value%d", i) {
			t.Fatalf("Expected value%d, got %s", i, value)
		}
	}

	key := "kkk"
	GlobalCache.SetStringWithExpiration(key, "vvv", 3*time.Second)
	value, err := GlobalCache.GetString(key)
	fmt.Println(value)
	time.Sleep(5 * time.Second)
	value, err = GlobalCache.GetString(key)
	fmt.Println(value)
}

func TestDeleteKey(t *testing.T) {
	// 初始化缓存和监控器
	err := Init(Config{
		MaxSize:  100,
		PlanTime: "18:30:00",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	var res string
	err = CacheAnything("test", CaseHandler, "hello world", &res, 10*time.Second)
	if err != nil {
		t.Fatalf("CacheAnything failed: %v", err)
	}
	err = DeleteKey("test")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
