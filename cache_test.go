package localcache

import "testing"

func TestCache(t *testing.T) {
	c, _ := NewCache()

	a, _ := c.Get("no")
	if a != nil {
		t.Error("Getting a value that shouldn't exist:", a)
	}

	c.Put("yes", 1)

	af, _ := c.Get("yes")
	if af == nil {
		t.Error("Should find a value")
	} else if af.(int) != 1 {
		t.Error("Expected 1, actual: ", af.(int))
	}

	found, _ := c.Remove("yes")
	if !found {
		t.Error("Should find a value")
	}
	af, _ = c.Get("yes")
	if af != nil {
		t.Error("Should delete the value")
	}
}
