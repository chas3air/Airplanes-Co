package models

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{} `json:"value" xml:"value"`
	Expiration int         `json:"expiration" xml:"expiration"`
	SetedTime  time.Time   `json:"setedTime" xml:"setedTime"`
}

type CarrotCache struct {
	cache map[string]CacheItem
	mu    sync.RWMutex
}

type Message struct {
	Key   string    `json:"key"`
	Value CacheItem `json:"cacheItem"`
}
