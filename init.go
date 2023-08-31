package cache_anything

import (
	"time"
)

var GlobalCache *Cache     // 全局缓存
var GlobalWatcher *Watcher // 全局缓存监控
var LimitCh chan int       // 监听是否触发LRU的channel

type GetSwitch func() bool

type Config struct {
	MaxSize  int64
	PlanTime string //定时清理缓存的时间，格式: HH:MM:SS
}

// Init 显示调用
func Init(c Config) error {
	GlobalCache = NewCache()
	GlobalWatcher = NewWatcher(GlobalCache)
	GlobalWatcher.SetMaxSize(c.MaxSize)
	LimitCh = make(chan int, 1)
	//每天7点清理
	t, err := time.Parse("15:04:05", c.PlanTime)
	if err != nil {
		return err
	}
	GlobalWatcher.SetClearPlanTime(t.Hour(), t.Minute(), t.Second())

	// 启动协程定时清理缓存
	go GlobalWatcher.WatchClear()

	// 启动协程监听是否需要LRU
	go GlobalWatcher.WatchLimit()

	// 启动协程点听是否过期
	go GlobalWatcher.WatchExpiration()

	return nil
}

func Switch(switchController GetSwitch) bool {
	return switchController()
}
