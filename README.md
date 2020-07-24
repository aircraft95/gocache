安装：
```
go get gitlab.ghzs.com/liangyh/cache
```

使用：
简单模式：
```go
package main

import (
  "gitlab.ghzs.com/liangyh/cache"
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
ty是cache类型，支持lru simple
```go
package main

import (
  "gitlab.ghzs.com/liangyh/cache"
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