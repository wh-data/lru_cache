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
		//if expire, also delete in map and return
		if expired := ele.isExpired(); expired {
			delete(l.Elements, key)
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

func (l *LRUCache) ClearExpire() error {
	for _, ele := range l.Elements {
		if ele.isExpired() {
			err := l.delete(ele)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *LRUCache) GetCapacity() int32 {
	return l.Capacity
}

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

func (l *LRUCache) ViewMap() (map[string]interface{}, error) {
	elements := make(map[string]interface{})
	for k, e := range l.Elements {
		elements[k] = e.Val
	}
	return elements, nil
}

func (l *LRUCache) ViewExpire(key string) (string, error) {
	if ele, ok := l.Elements[key]; ok {
		//if expired, delete
		if ele.isExpired() {
			err := l.delete(ele)
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

/*
Methods for LRUElement
*/
func (e *LRUElement) isExpired() bool {
	if e.Expire >= 0 && time.Now().Sub(time.Unix(e.TimeStamp, 0)).Seconds() >= e.Expire {
		return true
	}
	return false
}
