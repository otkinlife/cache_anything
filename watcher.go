package cache_anything

import (
	"log"
	"time"
)

const (
	SizeKB = 1024
	SizeMB = 1024 * 1024
	SizeGB = 1024 * 1024 * 1024
)

type PlanTime struct {
	Hour int
	Min  int
	Sec  int
}

type Watcher struct {
	cache    *Cache    // 缓存对象指针
	maxSize  int64     // 最大缓存限制
	planTime *PlanTime // 定时清理缓存的时间点
	time     *time.Ticker
}

// NewWatcher 创建监视器
func NewWatcher(cache *Cache) *Watcher {
	return &Watcher{
		cache:   cache,
		maxSize: 128 * SizeMB, // 默认128M
	}
}

func (w *Watcher) SetMaxSize(size int64) {
	w.maxSize = size
}

func (w *Watcher) SetClearPlanTime(h, m, s int) {
	w.planTime = &PlanTime{
		Hour: h,
		Min:  m,
		Sec:  s,
	}
}

// WatchClear 启动定时清理计划
func (w *Watcher) WatchClear() {
	if w.planTime == nil {
		log.Fatalf("watch cache error: time plan not config")
	}
	// 计算到下一个执行时间点的时间间隔
	now := time.Now()
	executionTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		w.planTime.Hour,
		w.planTime.Min,
		w.planTime.Sec,
		0, now.Location(),
	)
	duration := executionTime.Sub(now)
	if duration < 0 {
		next := now.Add(time.Hour * 24)
		executionTime = time.Date(
			next.Year(),
			next.Month(),
			next.Day(),
			w.planTime.Hour,
			w.planTime.Min,
			w.planTime.Sec,
			0, next.Location(),
		)
		duration = executionTime.Sub(now)
	}

	// 使用time.AfterFunc()在指定时间间隔后执行Clear()方法
	time.AfterFunc(duration, func() {
		w.cache.Clear()
		// 之后每24小时执行一次
		w.time = time.NewTicker(24 * time.Hour)
		for range w.time.C {
			w.cache.Clear()
		}
	})
}

func (w *Watcher) WatchLimit() {
	for {
		select {
		case <-LimitCh:
			w.CheckSize()
		}
	}
}

func (w *Watcher) WatchExpiration() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		now := time.Now()
		w.cache.data.Range(func(key, value interface{}) bool {
			item := value.(*Item)
			if !item.GetExpiration().IsZero() && item.GetExpiration().Before(now) {
				log.Printf("cache expire:%s", key)
				w.cache.data.Delete(key)
				w.cache.size -= item.GetSize()
				w.cache.updateKeyList(key)
			}
			return true
		})
	}
}

// CheckSize 检查缓存大小是否超过最大值
func (w *Watcher) CheckSize() {
	if w.cache.size > w.maxSize {
		w.cache.removeLeastAccessed() // 移除访问次数最少的缓存数据
	}
}
