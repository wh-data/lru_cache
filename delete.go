package lru_cache

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

func (l *LRUCache) deleteLru() (string, error) {
	//1. delete ele in link
	key, err := l.deleteLinkHead()
	if err != nil {
		return key, err
	}
	//2. delete in map
	delete(l.Elements, key)
	l.Size--
	return key, nil
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
