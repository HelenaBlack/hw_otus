package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}
type CacheItem struct {
	key   Key
	value interface{}
}
type lruCache struct {
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
	item, exist := l.items[key]
	if exist {
		cacheItem := item.Value.(*CacheItem)
		cacheItem.value = value
		l.queue.MoveToFront(item)
		return true
	}

	if l.queue.Len() >= l.capacity {
		lastItem := l.queue.Back()
		lastCacheItem := lastItem.Value.(*CacheItem)
		delete(l.items, lastCacheItem.key)
		l.queue.Remove(lastItem)
	}
	newCacheItem := &CacheItem{key: key, value: value}
	newItem := l.queue.PushFront(newCacheItem)
	l.items[key] = newItem

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, exist := l.items[key]
	if exist {
		l.queue.MoveToFront(item)
		cacheItem := item.Value.(*CacheItem)
		return cacheItem.value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
