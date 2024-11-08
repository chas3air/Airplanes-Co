package models

import (
	"sync"
	"time"
)

// CacheItem represents a single item stored in the cache.
// It includes the value, the expiration time, and the time it was set.
type CacheItem struct {
	Value      interface{} `json:"value" xml:"value"`           // The value stored in the cache.
	Expiration int         `json:"expiration" xml:"expiration"` // Time-to-live (TTL) in seconds.
	SetedTime  time.Time   `json:"setedTime" xml:"setedTime"`   // The time when the item was set.
}

// CarrotCache represents a simple in-memory cache with expiration capabilities.
// It uses a read-write mutex to manage concurrent access to the underlying cache.
type CarrotCache struct {
	cache map[string]CacheItem // The underlying map for storing cache items.
	mu    sync.RWMutex         // Mutex to manage concurrent access to the cache.
}

// Message represents a structure for handling JSON requests to set cache items.
// It contains a key for the cache and a CacheItem representing the value to be stored.
type Message struct {
	Key   string    `json:"key"`       // The key under which the cache item will be stored.
	Value CacheItem `json:"cacheItem"` // The CacheItem to be stored in the cache.
}

