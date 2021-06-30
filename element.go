package lru_cache

import "time"

type LRUElement struct {
	Prev      interface{}
	Next      interface{}
	Key       string
	Val       interface{}
	Expire    float64 //check only before "Get" or "CheckExpire
	TimeStamp int64
}

/*
Methods for LRUElement
*/
func (e *LRUElement) isExpired() bool {
	if e.Expire >= 0 && time.Now().Sub(time.Unix(e.TimeStamp, 0)).Seconds() >= e.Expire {
		return true
	}
	return false
}
