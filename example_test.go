package localcache_test

func Example() {
	cache, _ := NewCache()
	cache.Put("a", 1)
	a, _ := cache.Get("a")
}

func Example_custom() {
	cache, _ := NewCache(1000)
	cache.Put("a", 1)
	a, _ := cache.Get("a")
}
