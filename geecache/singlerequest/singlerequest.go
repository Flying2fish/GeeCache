package singlerequest

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Doer struct {
	mu sync.Mutex
	m  map[string]*call
}

func (d *Doer) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	d.mu.Lock()
	if d.m == nil {
		d.m = make(map[string]*call)
	}
	if c, ok := d.m[key]; ok {
		d.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	d.m[key] = c
	d.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	d.mu.Lock()
	delete(d.m, key)
	d.mu.Unlock()

	return c.val, c.err
}
