## cache

### 概述
cache提供了简单的缓存机制，要注意使用的是**服务器本身的内存**

### 快速开始

```bash
go get github.com/otkinlife/cache_anything
```

```golang
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
    //TODO: handle error
}
```

### 使用方式

```golang

// 初始化缓存
if cache.Switch() {
    err := cache.Init(cache.Config{
        MaxSize:  50 * cache.SizeMB, // 设置缓存可用的最大内存，当超过时Cache会基于LRU淘汰数据
        PlanTime: "07:00:00",        // 定时清理时间，每天定时清理所有的缓存
    })
    if err != nil {
        logger.Error(err)
    }
}

// 尝试从cache中取结果
if cache.Switch() {
    cacheRes := new(XXX)
    if err := cache.GlobalCache.LoadDataFromJson(cacheKey, cacheRes); err == nil && cacheRes != nil {
        //TODO: handle error
    }
}
// 异步写缓存写入缓存
if cache.Switch() {
    go cache.GlobalCache.SetDataWithJsonWithExpiration(cacheKey, res, 60*time.Second)
}
```