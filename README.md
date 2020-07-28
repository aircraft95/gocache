安装：
```
go get -u github.com/liangyouheng/gocache
```

使用：

简单模式：
```go
package main

import (
    "fmt"
    "github.com/liangyouheng/gocache"
)

func main() {
  c := cache.New()
  c.Set("key", []byte("ok"))
  value, err := c.Get("key")
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", string(value))
}
```

配置模式：
Ty是cache类型，支持Lru, Byte, Map
```go
package main

import (
    "fmt"
    "github.com/liangyouheng/gocache"
)

func main() {
  c := cache.NewWithConfig(cache.Config{
  		ShardsNum:       256,
  		DefaultSize:     1000,
  		DefaultValueLen: 5000,
  		Ty: 		cache.Byte,
  	})
  c.Set("key", []byte("ok"))
  value, err := c.Get("key")
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", string(value))
}
```

bench 测试：
```
go test -bench=. -benchmem -benchtime=4s -cpu=8
```