package lru_cache

import (
	"fmt"
	"runtime"
	"time"
)

func (l *LRUCache) Set(key string, val interface{}, expire float64) error {
	//check exceed capacity (key count)
	if l.Size >= l.Capacity {
		key, err := l.deleteLru()
		if err != nil {
			return fmt.Errorf("err when lru to deleteElement old key, err: %v", err)
		}
		fmt.Printf("lru deleteElement old key: %s\n", key)
	}
	//check exceed max mem
	if l.exceedMaxMem() {
		runtime.GC()
		time.Sleep(1 * time.Second)
		//lru
		for l.exceedMaxMem() {
			key, err := l.deleteLru()
			if err != nil {
				return fmt.Errorf("err when lru to delete old key: %s, err: %v", key, err)
			}
			runtime.GC()
			fmt.Printf("lru delete old key: %s, new mem: %d\n", key, realtimeMem)
		}
	}
	//check exist
	if ele, exist := l.Elements[key]; exist {
		//delete in both link and map
		err := l.deleteElement(ele)
		if err != nil {
			return err
		}
	}
	//new obj
	ele := &LRUElement{
		Prev:      nil,
		Next:      nil,
		Key:       key,
		Val:       val,
		Expire:    expire,
		TimeStamp: time.Now().Unix(),
	}
	//set new val in list tail
	err := l.setLinkTail(ele)
	if err != nil {
		return err
	}
	//set new val in map
	err = l.setMap(key, ele)
	if err != nil {
		return err
	}
	//adjust size
	l.Size++
	return nil
}

//SetLinkTail serves Get
//it also an base part for Set
func (l *LRUCache) setLinkTail(ele *LRUElement) error {
	old := l.Tail.Prev.(*LRUElement)
	l.Tail.Prev = ele
	ele.Prev = old
	old.Next = ele
	ele.Next = l.Tail
	return nil
}

//SetMap serves Get
//it also an base part for Set
func (l *LRUCache) setMap(key string, ele *LRUElement) error {
	l.Elements[key] = ele
	return nil
}
