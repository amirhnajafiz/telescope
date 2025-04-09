package cache

import "testing"

func TestSegmentCache(t *testing.T) {
	c := NewCache()

	if c.IsCached("abc") {
		t.Error("expected abc to not be cached yet")
	}

	c.MarkCached("abc")

	if !c.IsCached("abc") {
		t.Error("expected abc to be cached")
	}

	if size := c.Size(); size != 1 {
		t.Errorf("expected size to be 1, got %d", size)
	}
}
