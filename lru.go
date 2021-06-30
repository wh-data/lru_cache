package lru_cache

import (
	"sync"
)

type LRUCache struct {
	Head     *LRUElement
	Tail     *LRUElement
	Elements map[string]*LRUElement
	Size     int32
	Capacity int32
	MaxM     int64 //Bytes
	lock     *sync.Mutex
}

func NewLRUCache(capacity int32, maxM int64) *LRUCache {
	cache := &LRUCache{
		Head:     &LRUElement{},
		Tail:     &LRUElement{},
		Elements: make(map[string]*LRUElement),
		Size:     0,
		Capacity: capacity,
		MaxM:     maxM,
		lock:     new(sync.Mutex),
	}
	if cache.MaxM <= 0 {
		cache.MaxM = 5 * 1024 * 1024 * 1025 //default 5GB
	}
	cache.Head.Prev = nil
	cache.Head.Next = cache.Tail
	cache.Tail.Next = nil
	cache.Tail.Prev = cache.Head
	return cache
}
