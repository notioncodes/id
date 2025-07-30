//go:build !ristretto

package id

import (
	"container/list"
	"sync"
)

// cacheEntry represents a single cache entry with its key for easy removal.
type cacheEntry struct {
	key   string
	value string
}

// NoOpCache implements a simple in-memory cache with LRU eviction to maintain a 1MB size limit.
// It tracks insertion order using a doubly-linked list and removes the oldest entries
// when the cache exceeds the maximum size.
type NoOpCache struct {
	mu       sync.RWMutex
	cache    map[string]*list.Element
	order    *list.List
	maxSize  int64
	currSize int64
}

// NewNoOpCache creates a new NoOpCache instance with a 1MB size limit.
//
// Returns:
// - A Cache interface implementation that maintains size constraints through LRU eviction
func NewNoOpCache() Cache {
	return &NoOpCache{
		cache:   make(map[string]*list.Element),
		order:   list.New(),
		maxSize: 1024 * 1024, // 1MB
	}
}

// Get retrieves a value from the cache by key.
//
// Arguments:
// - key: The cache key to look up
//
// Returns:
// - The cached value and true if found, empty string and false otherwise
func (c *NoOpCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if elem, ok := c.cache[key]; ok {
		// Move to front (most recently used)
		c.order.MoveToFront(elem)
		entry := elem.Value.(*cacheEntry)
		return entry.value, true
	}
	return "", false
}

// Set stores a key-value pair in the cache, evicting oldest entries if necessary to maintain size limit.
//
// Arguments:
// - key: The cache key to store
// - value: The value to associate with the key
func (c *NoOpCache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	entrySize := int64(len(key) + len(value))

	// If this single entry exceeds max size, don't cache it
	if entrySize > c.maxSize {
		return
	}

	// If key already exists, update it.
	if elem, ok := c.cache[key]; ok {
		entry := elem.Value.(*cacheEntry)
		oldSize := int64(len(entry.key) + len(entry.value))
		entry.value = value
		c.currSize = c.currSize - oldSize + entrySize
		c.order.MoveToFront(elem)
		c.evictIfNeeded()
		return
	}

	// Else, add new entry.
	entry := &cacheEntry{key: key, value: value}
	elem := c.order.PushFront(entry)
	c.cache[key] = elem
	c.currSize += entrySize

	c.evictIfNeeded()
}

// evictIfNeeded removes the oldest entries until the cache size is within the limit.
func (c *NoOpCache) evictIfNeeded() {
	for c.currSize > c.maxSize {
		// Remove least recently used (back of list).
		elem := c.order.Back()
		if elem == nil {
			break
		}

		entry := elem.Value.(*cacheEntry)
		entrySize := int64(len(entry.key) + len(entry.value))

		c.order.Remove(elem)
		delete(c.cache, entry.key)
		c.currSize -= entrySize
	}
}
