package lru_cache

type LRUElement struct {
	Prev      interface{}
	Next      interface{}
	Key       string
	Val       interface{}
	Expire    float64 //check only before "Get" or "ViewExpire
	TimeStamp int64
}

type LRUCache struct {
	Head     *LRUElement
	Tail     *LRUElement
	Elements map[string]*LRUElement
	Size     int32
	Capacity int32
	MaxM     int32 //Bytes
}

func NewLRUCache(capacity, maxM int32) *LRUCache {
	cache := &LRUCache{
		Head:     &LRUElement{},
		Tail:     &LRUElement{},
		Elements: make(map[string]*LRUElement),
		Size:     0,
		Capacity: capacity,
		MaxM:     maxM,
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
