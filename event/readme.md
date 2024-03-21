# event
Events are meant to separate time-consuming operations from the main thread.  This project implements a distribution subscription pattern and provides a staging of messages.
## intro
- Produce once, consume many times
- Support for concurrent consumers
- The event is based on the channel, and the default listening queue capacity is 1000.
- Suitable for single machine
## usage
### As publisher
1. **init event**：
- The name cannot be repeated.
- The capacity means the capacity of the queue. If it is less than or equal to 0, the capacity is set to 1000.
   ```go
   event.InitEvent(name string, capacity int)
   ```
2. **trigger event**：

An event trigger can be analogous to a producer or broadcast sender in a queue.  Event trigger Triggers an event by event name.
   ```go
      event.Fire("eventName", []byte("event message!"))
   ```
### As subscriber
1. **event consumer implement**：

Event consumers need to implement the Eventer interface
   ```go
   type Eventer interface {
      ListenName() string
      Concurrence() int64
      Handle(ctx context.Context, data []byte)
   }
   ```
**tips:**

If the rate of production is greater than the rate of consumption, the channel will be full, and production will be blocked. You'd better do adequate testing and balance it out with consumer concurrency.
2. **event consumer register**：

Event consumer registration, that is, declaring the events to which they want to subscribe
   ```go
   event.Register(&handler.Test{})
   ```
## example

In your project start-up process
```go
//init
for eventName, capacity := range events {
    event.InitEvent(envName, capacity)
}

//register
event.Register(&handler.YourEvent1{})
event.Register(&handler.YourEvent2{})
...

//run
event.Run()
```

In your business processes
```go
//trigger
event.Fire("test", []byte("event message!"))
```
