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
}

func NewLRUCache(capacity int32) *LRUCache {
	cache := &LRUCache{
		Head:     &LRUElement{},
		Tail:     &LRUElement{},
		Elements: make(map[string]*LRUElement),
		Size:     0,
		Capacity: capacity,
	}
	cache.Head.Prev = nil
	cache.Head.Next = cache.Tail
	cache.Tail.Next = nil
	cache.Tail.Prev = cache.Head
	return cache
}
