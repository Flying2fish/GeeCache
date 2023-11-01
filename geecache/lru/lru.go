package lru

import "container/list"

type Value interface {
	Len() int
}

type Cache struct {
	maxBytes int64
	nowBytes int64

	ll       *list.List
	elements map[string]*list.Element

	OnEvicted func(key string, value Value)
}

type pair struct {
	key   string
	value Value
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		elements: make(map[string]*list.Element),

		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.elements[key]; ok {
		c.ll.MoveToFront(ele)

		p := ele.Value.(*pair)
		c.nowBytes += int64(value.Len()) - int64(p.value.Len())
		p.value = value
	} else {
		ele := c.ll.PushFront(&pair{key, value})
		c.elements[key] = ele

		c.nowBytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.elements[key]; ok {
		p := ele.Value.(*pair)
		c.ll.MoveToFront(ele)

		return p.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		p := ele.Value.(*pair)

		delete(c.elements, p.key)

		c.nowBytes -= int64(len(p.key)) + int64(p.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(p.key, p.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
