package localcache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	_ "github.com/allegro/bigcache"
	"github.com/golang/groupcache/lru"
	hashicorp "github.com/hashicorp/golang-lru"
	gocache "github.com/patrickmn/go-cache"
)

func BenchmarkGoCache(b *testing.B) {
	c := gocache.New(1*time.Minute, 5*time.Minute)

	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Add(toKey(i), randomVal(), gocache.DefaultExpiration)
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

// bigcache
func BenchmarkBigCache(b *testing.B) {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Set(toKey(i), []byte(fmt.Sprintf("%v", randomVal())))
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

// groupcache
func BenchmarkGroupCache(b *testing.B) {
	c := lru.New(102400)
	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Add(toKey(i), randomVal())
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

func BenchmarkGolangLRU(b *testing.B) {
	// c := cache2go.Cache("test")
	c, _ := hashicorp.New(102400)

	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Add(toKey(i), randomVal())
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

type Val struct {
	i int
	k string
	l []int
	b bool
	f float64
}

func randomVal() Val {
	t := rand.Int()
	l := make([]int, 1)
	l[0] = t
	return Val{
		i: rand.Int(),
		l: l,
		b: true,
		f: rand.Float64(),
	}
}
