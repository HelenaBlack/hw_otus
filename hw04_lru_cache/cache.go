package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	el, exist := l.items[key]
	if exist {
		l.queue.MoveToFront(el)
		el.Value = value
		return true
	}

	if l.queue.Len() >= l.capacity {
		lru := l.queue.Back()
		l.queue.Remove(lru)
		delete(l.items, lru.Value.(Key))
	}

	newItem := l.queue.PushFront(value)
	l.items[key] = newItem

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	el, exist := l.items[key]
	if exist {
		l.queue.MoveToFront(el)
		return el.Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
