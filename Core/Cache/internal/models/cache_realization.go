package models

import (
	"log"
	"time"
)

func NewCarrotCache() *CarrotCache {
	cc := &CarrotCache{
		cache: make(map[string]CacheItem),
	}
	go cc.startCleanup()
	go cc.startShowing()
	return cc
}

func (cc *CarrotCache) Get(key string) (interface{}, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	item, found := cc.cache[key]
	if !found {
		return nil, false
	}

	if time.Since(item.SetedTime).Seconds() > float64(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}
func (cc *CarrotCache) Set(key string, value interface{}, ttl int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cacheItem := CacheItem{
		Value:      value,
		Expiration: ttl,
		SetedTime:  time.Now(),
	}

	cc.cache[key] = cacheItem
}

func (cc *CarrotCache) Delete(key string) interface{} {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	out, ok := cc.cache[key]
	if !ok {
		log.Println("Element dont exists")
		return nil
	}
	delete(cc.cache, key)

	return out
}

func (cc *CarrotCache) Showing() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for i, v := range cc.cache {
		log.Println(i, v)
	}
}

func (cc *CarrotCache) startShowing() {
	for {
		cc.Showing()
		time.Sleep(4 * time.Second)
	}
}

func (cc *CarrotCache) cleanUp() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for i, v := range cc.cache {
		if time.Since(v.SetedTime).Seconds() > float64(v.Expiration) {
			delete(cc.cache, i)
		}
	}
}

func (cc *CarrotCache) startCleanup() {
	for {
		cc.cleanUp()
		time.Sleep(5 * time.Second)
	}
}
