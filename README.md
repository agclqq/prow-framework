# 简介
本项目指在提供一个开箱即用的go开发环境。

在基础功能中，提供通用开发方式，使开发者精力专注于业务开发，预留扩展接口，方便功能扩展。

在项目指导中，提供DDD的开发指导，使开发者能够更好的理解DDD，并在实践中逐步掌握DDD的精髓。

## 目录结构
以下是目录结构和DDD四层结构的对应关系（因前后端分离，故不包含user interface layer）

|--application       //对应 application layer，放的是controller <br/>
|--domain            //对应 domain layer，放的是service <br/>
|--infrastructure    //对应 infrastructure layer，放的是基础组件 <br/>

## 实践指导
### controller
1. controller只能依赖下层，即domain,infra，但尽可能只依赖domain
2. 每一个controller的方法，大体只负责三个步骤：
    * 检查参数
    * 调用domain逻辑并组装
    * 返回结果
### domain
1. domain允许同层依赖，但尽可能减少，不能向上依赖
2. 每个domain，职责要单一，通过组合完成复杂功能
### infra
1. 每个infra独立完成一个功能，不依赖上层

## 开发规范
[开发规范](guideline.md)

[安全规范](https://github.com/Tencent/secguide/blob/main/Go%E5%AE%89%E5%85%A8%E6%8C%87%E5%8D%97.md)
## tips
[开发tips](tips.md)

## 工具
### artisan
[artisan 介绍](artisan/Readme.md)

## 组件
### event
[event 介绍](event/readme.md)