package cache

var shardsNum = 256

type cache struct {
	shards 		[]shard
	Config
}

type initShard func(config Config) shard

type Config struct {
	ShardsNum int
	DefaultSize int
	DefaultValueLen int
	Ty     initShard
}

var (
	Map = initNewMapShard
 	Lru = initNewLruShard
 	Byte = initNewByteShard
)



func New() *cache {
	config := Config{
		ShardsNum:       shardsNum,
		DefaultSize:     1000,
		DefaultValueLen: 5000,
		Ty:        Map,
	}
	c := new(cache)
	c.Config = config
	return c.initShard(config)
}

func NewWithConfig(config Config) *cache {
	if config.ShardsNum == 0 {
		config.ShardsNum = 256
	}

	if config.DefaultSize == 0 {
		config.DefaultSize = 1000
	}

	if config.DefaultValueLen == 0 {
		config.DefaultValueLen = 5000
	}

	if config.Ty == nil {
		config.Ty = Map
	}

	c := new(cache)
	c.Config = config
	return c.initShard(config)
}


func (c *cache) initShard(config Config) *cache {
	c.shards = make([]shard, c.ShardsNum)
	for i := 0; i < c.ShardsNum; i++ {
		c.shards[i] = c.Ty(config)

	}
	return c
}

func (c *cache) getShard(hashedKey uint32) shard {
	return c.shards[hashedKey & uint32(c.ShardsNum - 1)]
}

func (c *cache) Set(key string, value []byte) {
	hashedKey := fnv32(key)
	c.getShard(hashedKey).set(hashedKey, value)
}

func (c *cache) Get(key string) ([]byte, error) {
	hashedKey := fnv32(key)
	return c.getShard(hashedKey).get(hashedKey)
}

func (c *cache) Del(key string) (bool, error) {
	hashedKey := fnv32(key)
	return c.getShard(hashedKey).del(hashedKey)
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

type shard interface {
	set(hashedKey uint32, value []byte)
	get(hashedKey uint32) ([]byte, error)
	del(hashedKey uint32) (bool, error)
}