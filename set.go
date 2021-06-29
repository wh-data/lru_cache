package lru_cache

import (
	"fmt"
	"runtime"
	"time"
)

func (l *LRUCache) Set(key string, val interface{}, expire float64) error {
	//check max mem
	if l.exceedMaxMem() {
		runtime.GC()
		time.Sleep(1 * time.Second)
		//lru
		for l.exceedMaxMem() {
			key, err := l.deleteLinkHead()
			if err != nil {
				return fmt.Errorf("err when lru to delete old key, err: %v", err)
			}
			runtime.GC()
			fmt.Println("lru delete old key: ", key)
		}
	}
	//check exist
	if ele, exist := l.Elements[key]; exist {
		err := l.delete(ele)
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
	//check exceed capacity, it happen after added new element, but at most it will exceed 1 ele
	if l.Size > l.Capacity {
		//delete in list
		key, err := l.deleteLinkHead()
		if err != nil {
			return err
		}
		//delete in map
		delete(l.Elements, key)
		//adjust size
		l.Size = l.Capacity
	}
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

func (l *LRUCache) delete(ele *LRUElement) error {
	//1. delete ele in link
	err := l.deleteLinkElement(ele)
	if err != nil {
		return err
	}
	//2. delete in map
	delete(l.Elements, ele.Key)
	l.Size--
	return nil
}

//deleteLinkElement is a base part of Delete
//gc will restore the deleted obj space
func (l *LRUCache) deleteLinkElement(ele *LRUElement) error {
	if l.Size < 1 {
		return nil
	}
	prev := ele.Prev.(*LRUElement)
	next := ele.Next.(*LRUElement)
	prev.Next = next
	next.Prev = prev
	return nil
}

func (l *LRUCache) deleteLinkHead() (string, error) {
	if l.Size < 1 {
		return "", nil
	}
	old := l.Head.Next.(*LRUElement)
	oldNext := old.Next.(*LRUElement)
	l.Head.Next = oldNext
	oldNext.Prev = l.Head
	return old.Key, nil
}
