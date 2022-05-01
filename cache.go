package go-localcache

import (
	"math/bits"
	"runtime"
	"sync"
	"time"
)

const (
	DEFAULT_EXPIRATION     = 10 * time.Minute
	DEFAULT_CLEAN_DURATION = 10 * time.Minute
	DEFAULT_CAP            = 1024
	DEFAULT_LRU_CLEAN_SIZE = 20
)

type Cache struct {
	defaultExpiration time.Duration
	elements          map[string]Elem
	capacity          int64
	size              int64
	lock              *sync.RWMutex
	pool              *sync.Pool
	cleaner           *cleaner
}

type Elem struct {
	K          string
	V          interface{}
	Expiration int64
	LastHit    int64
}

type cleaner struct {
	Interval time.Duration
	stop     chan bool
}

func (c *Cache) Get(k string) (v interface{}, err error) {
	ele := c.pool.Get()
	if item, ok := ele.(Elem); ok {
		if item.K == k {
			return item.V, nil
		}
	}
	expire := time.Now().Add(DEFAULT_EXPIRATION).UnixNano()
	lastHit := time.Now().UnixNano()
	c.lock.RLock()
	defer c.lock.RUnlock()
	if ele, ok := c.elements[k]; ok {
		ele.Expiration = expire
		ele.LastHit = lastHit
		return ele.V, nil
	}
	return nil, nil
}

func (c *Cache) Put(k string, v interface{}) error {
	expire := time.Now().Add(DEFAULT_EXPIRATION).UnixNano()
	lastHit := time.Now().UnixNano()
	if c.size+1 > c.capacity {
		// LRU kicks in
		if err := c.removeLeastVisited(); err != nil {
			return err
		}
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if found := c.update(k, v, expire, lastHit); found {
		return nil
	}

	ele := Elem{
		V:          v,
		Expiration: expire,
		LastHit:    lastHit,
	}
	c.pool.Put(ele)
	c.elements[k] = ele
	c.size = c.size + 1
	return nil
}

func (c *Cache) update(k string, v interface{}, expire int64, lastHit int64) bool {
	if ele, ok := c.elements[k]; ok {
		ele.V = v
		ele.Expiration = expire
		ele.LastHit = lastHit
		return true
	}
	return false
}

func (c *Cache) removeLeastVisited() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var lastTime int64 = 1<<(bits.UintSize-1) - 1 // MaxInt
	t := time.Now().UnixNano()
	lastItems := make([]string, DEFAULT_LRU_CLEAN_SIZE)
	liCount := 0
	full := false

	for k, v := range c.elements {
		if v.Expiration > t { // not expiring
			atime := v.LastHit
			if full == false || atime < lastTime {
				lastTime = atime
				if liCount < DEFAULT_LRU_CLEAN_SIZE {
					lastItems[liCount] = k
					liCount++
				} else {
					lastItems[0] = k
					liCount = 1
					full = true
				}
			}
		}
	}

	for i := 0; i < len(lastItems) && lastItems[i] != ""; i++ {
		lastName := lastItems[i]
		delete(c.elements, lastName)
	}
	return nil
}

func (c *Cache) Remove(k string) (isFound bool, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	v := c.pool.Get()
	if v != nil && v.(Elem).K != k {
		c.pool.Put(v)
	}
	for key := range c.elements {
		if key == k {
			delete(c.elements, key)
			return true, nil
		}
	}
	return false, nil
}

func (c *Cache) Flush() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.pool.Get()
	c.elements = make(map[string]Elem, DEFAULT_CAP)
	return nil
}

// Used in cleaning job
func (c *Cache) RemoveExpired() {
	now := time.Now().UnixNano()
	c.lock.Lock()
	defer c.lock.Unlock()
	for k, v := range c.elements {
		if v.Expiration > 0 && now > v.Expiration {
			_, _ = c.Remove(k)
		}
	}
}

// Run cleaning job
func (cl *cleaner) Run(c *Cache) {
	ticker := time.NewTicker(cl.Interval)
	for {
		select {
		case <-ticker.C:
			c.RemoveExpired()
		case <-cl.stop:
			ticker.Stop()
			return
		}
	}
}

func stopCleaner(c *Cache) {
	c.cleaner.stop <- true
}

func NewCache() (*Cache, error) {
	return newCache(DEFAULT_CAP, DEFAULT_EXPIRATION, DEFAULT_CLEAN_DURATION)
}

func newCache(cap int64, expiration time.Duration, clean_duration time.Duration) (*Cache, error) {
	c := &Cache{
		defaultExpiration: expiration,
		elements:          make(map[string]Elem, cap),
		capacity:          cap,
		lock:              new(sync.RWMutex),
		cleaner: &cleaner{
			Interval: clean_duration,
			stop:     make(chan bool),
		},
		pool: &sync.Pool{},
	}

	go c.cleaner.Run(c)
	runtime.SetFinalizer(c, stopCleaner)
	return c, nil
}
