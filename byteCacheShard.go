package cache

import (
	"encoding/binary"
	"errors"
	"sync"
)
const headerBufferSize  = 8

//增加keyItems，valueItems，headerBuffer，valueLen
//headerBuffer是用来在计算值的长度转换成[]byte时使用的， uint32转换成[]byte长度是8 ，所以headerBufferSize设定是8，valueLen是该valueItems的长度，用作扩容
type byteCacheShard struct {
	keyItems        map[uint32]uint32
	valueItems      []byte
	index           int
	lock         sync.RWMutex
	headerBuffer []byte
	valueLen  int
}

func initNewByteShard(config Config) *byteCacheShard {
	return &byteCacheShard{
		keyItems:        make(map[uint32]uint32),
		valueItems:     make([]byte, config.DefaultValueLen),		//创建默认长度为defaultValueLen 的[]byte
		index:			1,
		headerBuffer:    make([]byte, headerBufferSize),		//创建长度为8的[]byte数组，用作转换value长度变成[]byte
		valueLen : config.DefaultValueLen,			//记录该分片的默认长度
	}
}

func (s *byteCacheShard) set(hashedKey uint32, value []byte) {
	s.lock.Lock()
	s.keyItems[hashedKey] = uint32(s.index)

	//获取value的长度
	dataLen := len(value)

	//如果这个分片的valueItems 长度已经不能满足这一次的存放，我们存放数据大概占用的是 value的值的长度跟value值占用[]byte的长度
	if len(s.valueItems) < s.index + headerBufferSize + dataLen {

		//如果该次增加的长度，比整个分片默认值长还大，就不是创建多一倍出来，而是根据这个倍数+1
		i := (headerBufferSize + dataLen) / s.valueLen

		count :=2
		if i > 1 {
			count = i + 1
		}

		//把分片的数据取出来
		oldItem := s.valueItems

		//分片的valueLen值更新
		s.valueLen = s.valueLen * count

		//根据新的valueLen 创建更大的valueItems来存放数据
		s.valueItems = make([]byte, s.valueLen)
		//把旧数据index之前的数据全部复制回新的valueItems，扩容完成
		copy(s.valueItems, oldItem[:s.index])
	}


	//value值的长度原本是int 利用headerBuffer 转换成uint32的[]byte类型的值
	binary.LittleEndian.PutUint32(s.headerBuffer, uint32(dataLen))

	//把该值放到valueItems
	s.index += copy(s.valueItems[s.index:], s.headerBuffer[:headerBufferSize])

	//然后把[]byte 的value放进valueItems，更新index
	s.index += copy(s.valueItems[s.index:], value[:dataLen])

	//实际上每一次set(),就是把value的长度转换成[]byte ,再加上[]byte类型的value ,存放进valueItems这个大仓库
	s.lock.Unlock()
}


func (s *byteCacheShard) get(hashedKey uint32) ([]byte, error) {
	s.lock.RLock()
	//读取这个key在valueItems开始的下标
	index, ok := s.keyItems[hashedKey]
	if !ok {
		s.lock.RUnlock()
		return nil, errors.New("key not found")
	}

	//读取这个value的长度
	valueLen := binary.LittleEndian.Uint32(s.valueItems[index : index + headerBufferSize])

	dst := make([]byte, valueLen)

	//在c.shards[shardIndex].valueItems 中 index+headerBufferSize 到 index+headerBufferSize+valueLen 区间 就是这个value的[]byte的值
	copy(dst,s.valueItems[index + headerBufferSize: index + headerBufferSize + valueLen])
	s.lock.RUnlock()
	return dst,nil
}

func (s *byteCacheShard) del(hashedKey uint32) (bool, error) {
	s.lock.Lock()
	//读取这个key在valueItems开始的下标
	_, ok := s.keyItems[hashedKey]
	if !ok {
		s.lock.Unlock()
		return false,nil
	}

	delete(s.keyItems, hashedKey)
	s.lock.Unlock()
	return true, nil
}