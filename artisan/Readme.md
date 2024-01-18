# artisan简介
artisan是一个自定义命令行工具，用于JOB类的快速开发。

使用方式： `artisan $command [$arg1 ...]`

其中$command的格式建议为 `$packageName`:`$commandName`
# 使用方式
1. 通过实现artisan.Commander接口，实现自定义命令，推荐

声明job实现类

```go
    type TestCommand struct {
    }
    
    func (t TestCommand) GetCommand() string {
        return "command:test"
    }
    
    func (t TestCommand) Usage() string {
        return "test command"
    }
    
    func (t TestCommand) Handle(context *artisan.Context) {
        fmt.Println("test command")
    }
```
artisan main，并编译为./artisan
```go
    func main() {
        art := artisan.New()
        art.Add(commands.TestCommand{})
		art.Add(xxxCommand{})
		...
        art.Run()
    }
```
编译后运行指定的job
```shell
    ./artisan command:test
```

2. 通过自定义方法来实现自定义命令
    ```go
    func main() {
        art := artisan.New()
        art.AddFunc("command:test", func(context *artisan.Context) {
            fmt.Println("test command")
        })
        art.Run()
    }
    ```
3. artisan.Run方法的参数
   artisan.Run方法的参数为可变参数，可以通过命令行传入，也可以在代码中传入。
   可通过程序直接调用，用于debug。
    ```go
    func main() {
        art.Run("your command", "arg1", "arg2","...")
    }
    ```

4. 运行
   build后执行，第一个参数为命令，后续参数为命令参数。命令参数可以通过`-`或`--`来指定参数名，也可以用无前缀参数名，支持等号或空格为参数赋值，如果参数无值，则认为是字符串空值。
    ```shell
    ./artisan command:test arg1 -arg2 2 arg7 7 --arg3 3 -arg4=4 --arg5=5 arg6=6
    ```
### 功能类型
artisan分为两大块内容
1. 官方命令
   1. 在infrastructure/artisan中实现，后续可拆为单独组件
   2. 现有功能包：
      1. make：创建资源，包括controller,model,grpc,command等
      2. list：查看资源
2. 自定义命令
   1. 实现artisan.Commander接口的业务内容
   2. 注册命令到artisan

### 使用示例
1. 创建controller
```text
go run cmd/artisan/artisan.go make:controller yourController
```
2. 创建model
```text
go run cmd/artisan/artisan.go make:model yourModel
```

