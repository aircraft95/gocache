package cache

import (
	"container/list"
	"errors"
	"sync"
)

type lruCacheShard struct {
	lock         sync.RWMutex
	items     	 map[interface{}]*list.Element
	list 		*list.List
	maxSize      	 int
}

type lruItem struct {
	key        uint32
	value      []byte
}

func initNewLruShard(config Config) *lruCacheShard {
	return &lruCacheShard{
		items: make(map[interface{}]*list.Element),
		list: list.New(),
		maxSize: config.DefaultSize,
	}
}

func (s *lruCacheShard) set(hashedKey uint32, value []byte,) {
	s.lock.Lock()

	//获取map里面是否存在该key
	if element, ok := s.items[hashedKey]; ok {
		//如果存在，把该元素移到链表的最前方
		s.list.MoveToFront(element)
		//元素重新赋值value 必须使用引用的方式，否则无法修改值
		element.Value.(*lruItem).value = value
		s.lock.Unlock()
		return
	}
	//在链表的最前方增加一个元素，元素的类型是引用的lruItem，如果不使用引用的形式，后续是无法修改这个值的
	element := s.list.PushFront(&lruItem{hashedKey, value})
	//把元素的引用赋值到map类型的items上
	s.items[hashedKey] = element
	//判断链表的长度，如果长度超过规定值，移除链表最后的数据
	if s.maxSize != 0 && s.list.Len() > s.maxSize {
		s.removeOldest()
	}
	s.lock.Unlock()
}


func (s *lruCacheShard) get(hashedKey uint32) ([]byte, error) {
	s.lock.RLock()
	//获取map里面是否存在该key
	if element, ok := s.items[hashedKey]; ok {
		//如果存在 把该元素移到链表的头部
		s.list.MoveToFront(element)
		s.lock.RUnlock()
		return element.Value.(*lruItem).value, nil
	}
	s.lock.RUnlock()
	return []byte{}, errors.New("key not found")
}

func (s *lruCacheShard) del(hashedKey uint32) (bool, error) {
	s.lock.Lock()
	if element, ok := s.items[hashedKey]; ok {
		s.removeElement(element)
	}
	s.lock.Unlock()
	return true, nil
}

// 删除链表的尾部的元素
func (s *lruCacheShard) removeOldest() {
	//获取链表的尾部的元素
	element := s.list.Back()
	if element != nil {
		//删除尾部元素
		s.removeElement(element)
	}
}

//删除链表指定元素
func (s *lruCacheShard) removeElement(e *list.Element) {
	//删除链表的元素
	s.list.Remove(e)
	//获取元素数据并转换成lruItem
	kv := e.Value.(*lruItem)
	//删除在items上的记录
	delete(s.items, kv.key)
}

// 返回链表长度
func (s *lruCacheShard) Len() int {
	return s.list.Len()
}