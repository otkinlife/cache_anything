package cache_anything

import (
	"time"
)

// Item 表示缓存中的一项
type Item struct {
	value        interface{} // 缓存项的值
	valueType    string      // 缓存项值的类型
	lastUsedTime time.Time   // 缓存项最后一次使用的时间
	size         int64       // 缓存项的大小（字节）
	accessCount  int64       // 缓存项的访问次数
	expiration   time.Time   // 缓存项的过期时间
}

// NewItem 创建一个带有默认字段值的新缓存项
func NewItem() *Item {
	return &Item{
		lastUsedTime: time.Time{},
		size:         0,
	}
}

// Get 返回缓存项的值
func (i *Item) Get() interface{} {
	return i.value
}

// Set 设置缓存项的值
func (i *Item) Set(v interface{}) {
	i.value = v
}

// SetExpiration 设置过期时间
func (i *Item) SetExpiration(t time.Time) {
	i.expiration = t
}

// GetExpiration 返回缓存项的过期时间
func (i *Item) GetExpiration() time.Time {
	return i.expiration
}

// GetValueType 返回缓存项值的类型
func (i *Item) GetValueType() string {
	return i.valueType
}

// SetValueType 设置缓存项值的类型
func (i *Item) SetValueType(v string) {
	i.valueType = v
}

// GetLastUsedTime 返回缓存项最后一次使用的时间
func (i *Item) GetLastUsedTime() time.Time {
	return i.lastUsedTime
}

// UpdateLastUsedTime 更新缓存项最后一次使用的时间为当前时间
func (i *Item) UpdateLastUsedTime() {
	i.lastUsedTime = time.Now()
}

// GetSize 返回缓存项的大小（字节）
func (i *Item) GetSize() int64 {
	return i.size
}

// SetSize 设置缓存项的大小（字节）
func (i *Item) SetSize(s int64) {
	i.size = s
}

// IncreaseAccessCount 增加缓存项的访问次数
func (i *Item) IncreaseAccessCount() {
	i.accessCount++
}
