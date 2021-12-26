package golang_localcache

import (
	"fmt"
	"testing"
	"time"

	hashicorp "github.com/hashicorp/golang-lru"
	gocache "github.com/patrickmn/go-cache"
)

func BenchmarkGoCache(b *testing.B) {
	c := gocache.New(1*time.Minute, 5*time.Minute)

	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Add(toKey(i), toKey(i), gocache.DefaultExpiration)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, found := c.Get(toKey(i))
			if found {
				_ = value
			}
		}
	})
}

func BenchmarkCache(b *testing.B) {
	c, _ := NewCache()
	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Put(toKey(i), toKey(i))
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			value, _ := c.Get(toKey(i))
			if value != nil {
				_ = value
			}
		}
	})
}

func BenchmarkHashicorpLRU(b *testing.B) {
	// c := cache2go.Cache("test")
	c, _ := hashicorp.New(1024)

	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Add(toKey(i), toKey(i))
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			value, err := c.Get(toKey(i))
			if err == true {
				_ = value
			}
		}
	})
}

func toKey(i int) string {
	return fmt.Sprintf("item:%d", i)
}