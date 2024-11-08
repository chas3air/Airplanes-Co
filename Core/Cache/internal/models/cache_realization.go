package models

import (
	"log"
	"time"
)

// NewCarrotCache initializes a new CarrotCache instance and starts a cleanup goroutine.
func NewCarrotCache() *CarrotCache {
	cc := &CarrotCache{
		cache: make(map[string]CacheItem), // Initialize the cache map.
	}
	go cc.startCleanup() // Start the cleanup process in a separate goroutine.
	return cc
}

// Get retrieves an item from the cache using the provided key.
// It returns the value and a boolean indicating whether the item was found.
// If the item has expired, it returns nil and false.
func (cc *CarrotCache) Get(key string) (interface{}, bool) {
	cc.mu.RLock() // Acquire a read lock for safe concurrent access.
	defer cc.mu.RUnlock()

	item, found := cc.cache[key]
	if !found {
		return nil, false // Item not found in the cache.
	}

	// Check if the item has expired.
	if time.Since(item.SetedTime).Seconds() > float64(item.Expiration) {
		return nil, false // Item has expired.
	}

	return item.Value, true // Return the item's value and indicate success.
}

// Set stores a value in the cache with the specified key and time-to-live (TTL) in seconds.
// It updates the item's set time to the current time.
func (cc *CarrotCache) Set(key string, value interface{}, ttl int) {
	cc.mu.Lock() // Acquire a write lock for safe concurrent access.
	defer cc.mu.Unlock()

	// Create a new CacheItem with the current time and expiration.
	cacheItem := CacheItem{
		Value:      value,
		Expiration: ttl,
		SetedTime:  time.Now(),
	}

	cc.cache[key] = cacheItem // Store the item in the cache.
}

// Delete removes an item from the cache using the specified key.
// It returns the value of the deleted item, or nil if the item was not found.
func (cc *CarrotCache) Delete(key string) interface{} {
	cc.mu.Lock() // Acquire a write lock for safe concurrent access.
	defer cc.mu.Unlock()

	out, ok := cc.cache[key]
	if !ok {
		log.Println("Element doesn't exist") // Log if the item was not found.
		return nil
	}
	delete(cc.cache, key) // Remove the item from the cache.

	return out // Return the deleted item.
}

// cleanUp removes expired items from the cache.
func (cc *CarrotCache) cleanUp() {
	cc.mu.Lock() // Acquire a write lock for safe concurrent access.
	defer cc.mu.Unlock()

	for i, v := range cc.cache {
		// Check if the item has expired.
		if time.Since(v.SetedTime).Seconds() > float64(v.Expiration) {
			delete(cc.cache, i) // Delete the expired item.
		}
	}
}

// startCleanup periodically invokes the cleanUp method to remove expired items from the cache.
func (cc *CarrotCache) startCleanup() {
	for {
		cc.cleanUp()                // Call the cleanup method.
		time.Sleep(5 * time.Second) // Wait for 5 seconds before the next cleanup cycle.
	}
}
