package cache_anything

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type KeyBuilder struct {
	Params         string //需要参与编码的参数，是一个json字符串
	Key            string
	LastUpdateTime time.Time
}

func NewKeyBuilder() *KeyBuilder {
	return &KeyBuilder{}
}

func (k *KeyBuilder) SetParams(params string) *KeyBuilder {
	k.Params = params
	return k
}

func (k *KeyBuilder) Build() string {
	if k.Params == "" {
		return ""
	}
	md5Hash := md5.Sum([]byte(k.Params))
	md5Str := hex.EncodeToString(md5Hash[:])
	k.LastUpdateTime = time.Now()
	timeStr := k.LastUpdateTime.Format("20060102")
	k.Key = fmt.Sprintf("%s_%s", md5Str, timeStr)
	return k.Key
}
