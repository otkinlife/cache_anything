package cache_anything

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Cache 缓存结构
type Cache struct {
	data     sync.Map      // 缓存的数据，这里使用sync.Map防止并发读写问题
	size     int64         // 缓存的数据大小，用于做限制，防止内存溢出
	keyList  []*cacheEntry // 全局的key列表，用于维护key的优先级。LRU淘汰的时候讲从该结构中选取优先级低的key淘汰
	keyMutex sync.Mutex    // 用于保护keyList的互斥锁
}

type cacheEntry struct {
	key          interface{}
	accessCount  int64
	lastUsedTime time.Time
}

func NewCache() *Cache {
	return &Cache{
		keyList: make([]*cacheEntry, 0),
	}
}

func (c *Cache) Clear() {
	fmt.Println("clear cache")
	// 直接赋空值，利用GC回收旧的Cache
	c.data = sync.Map{}
	c.size = 0
	c.keyList = make([]*cacheEntry, 0)
}

// LoadDataFromJson 如果存的值是json格式的字符串，可以通过该方法load到data里
func (c *Cache) LoadDataFromJson(params string, data interface{}) error {
	str, err := c.GetString(params)
	if err != nil {
		return err
	}
	if str == "" {
		return errors.New("get nothing")
	}
	return json.Unmarshal([]byte(str), data)
}

// SetDataWithJson data会以json字符串的形式存到cache中
func (c *Cache) SetDataWithJson(params string, data interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.SetString(params, string(jsonStr))
	return nil
}

// SetDataWithJsonWithExpiration 带有过期时间的缓存
func (c *Cache) SetDataWithJsonWithExpiration(params string, data interface{}, d time.Duration) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.SetStringWithExpiration(params, string(jsonStr), d)
	return nil
}

func (c *Cache) SetStringWithExpiration(params string, v string, d time.Duration) {
	key := NewKeyBuilder().SetParams(params).Build()
	stringItem := NewStringItem()
	stringItem.SetString(v)
	stringItem.SetExpiration(time.Now().Add(d))
	c.data.Store(key, stringItem.GetItem())
	c.size += int64(len(v))
	// 更新keyList
	c.updateKeyList(key)
	log.Printf("set cache:%s", key)

	// 给watcher发信号，校验是否超出size限制
	LimitCh <- 1
}

// SetString 写入cache
// params 用于生成key的因素
// v 存入cache的值
func (c *Cache) SetString(params string, v string) {
	key := NewKeyBuilder().SetParams(params).Build()
	stringItem := NewStringItem()
	stringItem.SetString(v)
	stringItem.SetExpiration(time.Time{})
	c.data.Store(key, stringItem.GetItem())
	c.size += int64(len(v))
	// 更新keyList
	c.updateKeyList(key)
	log.Printf("set cache:%s", key)
	// 给watcher发信号，校验是否超出size限制
	LimitCh <- 1
}

// GetString 获取cache
// params 用于生成key的因素
func (c *Cache) GetString(params string) (string, error) {
	key := NewKeyBuilder().SetParams(params).Build()
	if v, ok := c.data.Load(key); ok {
		if item, ok := v.(*Item); ok {
			stringItem := NewStringItem()
			err := stringItem.Load(item)
			if err != nil {
				return "", err
			}
			item.IncreaseAccessCount()
			return stringItem.GetString(), nil
		}
	}
	return "", nil
}

func (c *Cache) Delete(params string) error {
	// 构建缓存键
	cacheKey := NewKeyBuilder().SetParams(params).Build()

	// 检查缓存中是否存在该键
	value, ok := c.data.Load(cacheKey)
	if !ok {
		return nil
	}

	// 删除缓存中的键值对
	c.data.Delete(cacheKey)

	// 更新缓存大小
	if item, ok := value.(*Item); ok {
		c.size -= item.size
	}

	// 更新keyList
	c.keyMutex.Lock()
	defer c.keyMutex.Unlock()
	for i, entry := range c.keyList {
		if entry.key == cacheKey {
			c.keyList = append(c.keyList[:i], c.keyList[i+1:]...)
			break
		}
	}

	log.Printf("deleted cache key: %s", params)
	return nil
}

// 更新keyList，将新的key插入到合适的位置
func (c *Cache) updateKeyList(key interface{}) {
	c.keyMutex.Lock()
	defer c.keyMutex.Unlock()

	value, ok := c.data.Load(key)
	if !ok {
		return
	}

	item := value.(*Item)
	entry := &cacheEntry{
		key:          key,
		accessCount:  item.accessCount,
		lastUsedTime: item.lastUsedTime,
	}

	// 寻找插入位置
	insertIndex := -1
	for i, e := range c.keyList {
		if e.accessCount < entry.accessCount || (e.accessCount == entry.accessCount && e.lastUsedTime.Before(entry.lastUsedTime)) {
			insertIndex = i
			break
		}
	}

	// 插入新的key
	if insertIndex == -1 {
		c.keyList = append(c.keyList, entry)
	} else {
		c.keyList = append(c.keyList[:insertIndex], append([]*cacheEntry{entry}, c.keyList[insertIndex:]...)...)
	}
}

// removeLeastAccessed 根据keyList淘汰缓存数据
func (c *Cache) removeLeastAccessed() {
	c.keyMutex.Lock()
	defer c.keyMutex.Unlock()

	// 淘汰数量为当前缓存的1/3
	removeCount := len(c.keyList) / 3

	// 删除最后的n个元素，并更新Cache结构的size值
	for i := 0; i < removeCount; i++ {
		entry := c.keyList[len(c.keyList)-1-i]
		value, ok := c.data.Load(entry.key)
		if ok {
			item := value.(*Item)
			c.size -= item.size
			c.data.Delete(entry.key)
		}
	}

	// 更新keyList
	c.keyList = c.keyList[:len(c.keyList)-removeCount]
}
