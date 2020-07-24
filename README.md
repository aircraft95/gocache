安装：
```
go get github.com/liangyouheng/gocache
```

使用：

简单模式：
```go
package main

import (
  "github.com/liangyouheng/gocache"
  "fmt"
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
ty是cache类型，支持lru, byte, map
```go
package main

import (
  "github.com/liangyouheng/gocache"
  "fmt"
)

func main() {
  c := cache.NewWithConfig(cache.Config{
  		ShardsNum:       256,
  		DefaultSize:     1000,
  		DefaultValueLen: 5000,
  		Ty: 		"lru",
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