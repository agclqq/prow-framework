## event简介
event主要是为了解耦业务逻辑，实现一个事件，多个业务消费的模型。event基于channel，默认监听队列容量为1000。

使用方式：
1. **事件名初始化**：

注意：env name不能为空，且不重复；capacity为队列容量，默认为1000
   ```go
   event.InitEnvName(name string, capacity int)
   ```
2. **事件消费配置**：

事件消费者可以类比为队列中的消费者或广播接收者。事件消费者需实现Eventer接口，如下：
   ```go
   type Eventer interface {
      ListenName() string
      Concurrence() int64
      Handle(ctx context.Context, data []byte)
   }
   ```
3. **事件注册**：

事件消费者需要将自己注册到需要消费的事件中，如下：
   ```go
   event.Register(&handler.Test{})
   ```
4. **事件触发**：

事件触发者可以类比为队列中的生产者或广播发送者。事件触发者通过事件名触发事件，如下：
   ```go
      event.Fire("test", []byte("event message!"))
   ```

