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
}
