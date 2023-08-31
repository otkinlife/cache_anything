package cache_anything

import "errors"

type StringItem struct {
	*Item
}

const TypeString = "string"

func NewStringItem() *StringItem {
	item := NewItem()
	// 设置值
	item.SetValueType(TypeString)

	return &StringItem{
		item,
	}
}

func (s *StringItem) SetString(v string) {
	// 设置大小
	size := len(v)
	s.SetSize(int64(size))
	// 设置时间
	s.UpdateLastUsedTime()
	s.Set(v)
}

func (s *StringItem) GetString() string {
	if v, ok := s.Item.value.(string); ok {
		return v
	}
	return ""
}

func (s *StringItem) GetItem() *Item {
	return s.Item
}

// Load 将item类型加载为StringItem类型
func (s *StringItem) Load(i *Item) error {
	if i.GetValueType() != TypeString {
		return errors.New("load type error, want string but get a " + i.GetValueType())
	}
	s.Item = i
	return nil
}
