package cache

import (
	"errors"
	"sync"
)

type mapCache struct {
	items        map[uint32][]byte
	lock         sync.RWMutex
}

func initNewMapShard(config Config) *mapCache {
	return &mapCache{
		items:        make(map[uint32][]byte),
	}
}

func (s *mapCache) set(hashedKey uint32, value []byte) {
	s.lock.Lock()
	s.items[hashedKey] = value
	s.lock.Unlock()
}


func (s *mapCache) get(hashedKey uint32) ([]byte, error) {
	s.lock.RLock()
	value, ok := s.items[hashedKey]
	s.lock.RUnlock()

	if ok {
		return value, nil
	} else {
		return []byte{}, errors.New("key not found")
	}
}

func (s *mapCache) del(hashedKey uint32) (bool, error) {
	s.lock.Lock()
	_, ok := s.items[hashedKey]
	if !ok {
		s.lock.Unlock()
		return false,nil
	}
	delete(s.items, hashedKey)
	s.lock.Unlock()
	return true, nil
}