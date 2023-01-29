package lru_cache

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	fmt.Println("START CACHE")
	cache := NewCache()
	for _, word := range []string{"parrot", "avocado", "dragonfruit", "tree", "potato", "tomato", "tree", "dog"} {
		cache.Check(word)
		cache.Display()
	}
}
