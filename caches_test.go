package cache

import (
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/coocood/freecache"
	"sync"
	"testing"
	"time"
)

const (
	NumOfReader = 90
	NumOfWriter = 10
)

func BenchmarkCacheMap(b *testing.B) {
	newCache := New()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for i := 0; i < NumOfWriter; i++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					newCache.Set(key(i), []byte(value(i)))
				}
				wg.Done()
			}()
		}
		for i := 0; i < NumOfReader; i++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					newCache.Get(key(i))
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkBigCacheMap(b *testing.B) {
	newCache := initBigCache(1000 * 10 * 60)
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for i := 0; i < NumOfWriter; i++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					newCache.Set(key(i), []byte(value(i)))
				}
				wg.Done()
			}()
		}
		for i := 0; i < NumOfReader; i++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					newCache.Get(key(i))
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkFreeCacheMap(b *testing.B) {
	newCache := freecache.NewCache(256)
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for i := 0; i < NumOfWriter; i++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					newCache.Set([]byte(key(i)), []byte(value(i)),0)
				}
				wg.Done()
			}()
		}
		for i := 0; i < NumOfReader; i++ {
			wg.Add(1)
			go func() {
				for i := 0; i < 100; i++ {
					newCache.Get([]byte(key(i)))
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func value(i int) string {
	return fmt.Sprintf("value-%010d", i)
}

func initBigCache(entriesInWindow int) *bigcache.BigCache {
	cache, _ := bigcache.NewBigCache(bigcache.Config{
		Shards:             256,
		LifeWindow:         10 * time.Minute,
		MaxEntriesInWindow: entriesInWindow,
		MaxEntrySize:       256,
		Verbose:            false,
	})

	return cache
}