## event简介
event主要是为了解耦业务逻辑，实现一个事件，多个业务消费的模型。event基于channel，默认监听队列容量为1000。

使用方式：
1. **事件名配置**：在.env中配置event_names，以英文逗号分隔
2. **事件消费配置**：默认在application/event/中新增事件消费文件，当然也可以放其他目录，需实现以下接口：
   ```go
    type Eventer interface { 
        GetName() string
        Handle(data []byte)
    }
   ```
3. **事件注册**：在provider/event.go中配置要消费事件的处理方法，如下：
    ```go
    //事件消费的默认并发数量为1
    event.Register("test", &handler.Test{})
   
    //修改事件消费并发数量
    event.Register("test", &handler.Test{}, 3)
    ```
4. **事件触发**：
   ```go
      event.Fire("test", []byte("event message!"))
   ```

