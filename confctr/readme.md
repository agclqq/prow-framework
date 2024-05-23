# confctr
config center

## 1. intro
Config center for etcd, consul.

The configuration center is often used to store, update, and listen to key-value pair information.

If you want to extend it, just implement inter.go's interface
## 2. useage
Refer to the inter.go interface file for supporting information.
```go

```
## 3. example
```go
package main

import (
	consulapi "github.com/hashicorp/consul/api"
	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/agclqq/prow/confctr/manager"
)

func main(){
	conf:=manager.Config{
		Type:       manager.CCTypeConsul,
		EtcdConf:   clientV3.Config{Endpoints: []string{"127.0.0.1:2379"}},
		ConsulConf: &consulapi.Config{Address: "127.0.0.1:8500"},
	}
	cc, err := manager.New(conf)
	if err != nil {
		return
	}
	cc.Create("test", "test")
	cc.Get("test")
	cc.Update("test", "test1")
	cc.Delete("test")
	cc.Watch("test", func(key, val string) {
        //todo
    }
}
```