package lru_cache

import (
	"fmt"
	"time"
)

func (l *LRUCache) Get(key string) (interface{}, error) {
	if ele, ok := l.Elements[key]; ok {
		//delete (have to do it at beginning)
		err := l.deleteLinkElement(ele)
		if err != nil {
			return nil, err
		}
		//if expired, also delete in map and return
		if expired := ele.isExpired(); expired {
			delete(l.Elements, key)
			l.Size--
			return nil, err
		}
		//set to tail
		err = l.setLinkTail(ele)
		return ele.Val, err
	}
	return nil, nil
}

func (l *LRUCache) GetSize() (int32, error) {
	if l.Size != int32(len(l.Elements)) {
		return l.Size, fmt.Errorf("link size not equal map size, need to rebuild lru cache")
	}
	return l.Size, nil
}

func (l *LRUCache) GetCapacity() int32 {
	return l.Capacity
}

//ViewLinkedList view all elements in list
func (l *LRUCache) ViewLinkedList() ([]interface{}, error) {
	list := make([]interface{}, 0)
	next := l.Head.Next.(*LRUElement)
	for {
		if next == l.Tail {
			return list, nil
		}
		list = append(list, next.Val)
		next = next.Next.(*LRUElement)
	}
}

//ViewMap view all elements in map
func (l *LRUCache) ViewMap() (map[string]interface{}, error) {
	elements := make(map[string]interface{})
	for k, e := range l.Elements {
		elements[k] = e.Val
	}
	return elements, nil
}

//CheckExpire check whether key is expired, if not expired, return left life time
func (l *LRUCache) CheckExpire(key string) (string, error) {
	if ele, ok := l.Elements[key]; ok {
		//if expired, delete
		if ele.isExpired() {
			err := l.deleteElement(ele)
			return "", fmt.Errorf("key is not eixst %v", err)
		}
		//if not expire, calculate left time
		if ele.Expire > 0 {
			return fmt.Sprintf("%.2f seconds", ele.Expire-time.Now().Sub(time.Unix(ele.TimeStamp, 0)).Seconds()), nil
		} else {
			return fmt.Sprintf("%.2f seconds", ele.Expire), nil
		}
	}
	//not exist
	return "", fmt.Errorf("key is not eixst")
}

//ClearExpire clear all expired key in cache
func (l *LRUCache) ClearExpire() error {
	for _, ele := range l.Elements {
		if ele.isExpired() {
			err := l.deleteElement(ele)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
