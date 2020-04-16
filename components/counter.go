package components

import "sync"

type Counter struct {
	mx       sync.RWMutex
	counters map[string]int
}

func NewCounter() *Counter {
	return &Counter{
		counters: make(map[string]int),
	}
}

func (c *Counter) List() map[string]int {
	c.mx.RLock()
	defer c.mx.RUnlock()
	newMap := make(map[string]int)
	for k, v := range c.counters {
		newMap[k] = v
	}
	return newMap
}

func (c *Counter) Load(key string) (int, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.counters[key]
	return val, ok
}

func (c *Counter) Store(key string, value int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.counters[key] = value
}

func (c *Counter) Inc(key string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.counters[key]++
}
