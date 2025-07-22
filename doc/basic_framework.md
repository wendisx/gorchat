## Basic Framework

- [database](#database)
- [repository](#repository)
- [usecase](#usecase)
- [handler](#handler)
- [route](#route)

> <a id="database">database</a>

为了能够使得尽可能多地接入不同数据库，完全有理由为不同数据库实现自己独立的
接口，混用接口往往会导致意外的麻烦。数据库层面的初始化核心始终只是简单地尝试连接给定connUrl的连接，当然如果一定要排除假连接的情况最好在创建连接池后便尝
试测试连接。

```go
// 使用golang大致的db层
package database

import (
    ...
)

type [T->DB] interface {
}

func New[T]DB() [T->DB] {
   ... 
   return [T->DB]{}
}
```

> <a id="repository">repository</a>

repository层相当于一个数据库数据加工层，数据库的数据总是庞大的，在一定程度上
是混乱的，逻辑操作层需要一定格式的数据来处理，repo层作为了一个数据缓冲的地方
。为了组织数据，repo层显然需要通过database层的语言和接口尝试拿到一些数据并自
行组织，显然repo层一个较为深的层级，考虑当repo层发生预警错误时，其本身的错误
较为单一，并不涉及到逻辑层次的问题。此时的预警本身类似一种信号，这被认为是re
po与database交流时出现的问题，需要被详细记录，但是不需要被上层处理，因为这只
是表明下层出现了错误，该错误也许是一个严重的错误，但是没有任何理由使其需要返
回到客户端

```go
// 使用golang大致的repo层
package repository

import (
    ...
)

type [T->Repository] intereface {}

type [t->Repository] struct { 
    db [T->DB]
    log [T->Logger]
}

func New[T->Repository](db [T->DB],log [T->Logger]) [T->Repository] { 
    return &[t->Repository]{ 
        db: db,
        log: log,
    }
}

func (r *[r->Repository]) <_>() [T->Logger] { 
    return r.log
}

// unix哲学 -- 只做一件事，做好一件事
// motion --> insert | delete | update | select
// pointer --> 谁调用，谁提供容器
func (r *[r->Repository]) <[motion][object]>(ctx context.Context,container...) error {
    [motion->Sql] := `
        ...
    `
    result,err := (exec $[motion->Sql] $container)
    with err!=nil (solve err #log #return)
    with result (parse result #return)
}

// 这里的错误大致为
type [T->RError] struct {
    Code int
    Level int
}
```

> <a id="usecase">usecase</a>

usecase就是核心的logic处理层，想象有一个来自上层解析出来的一定格式的数据，这
是另一种混乱格式，usecase完全有理由拒绝处理这种格式的数据，因此需要经过初步
的格式化后传递到usecase层，uscase有自己的处理风格，由于repo层的调用接口的参
数通常会存在较大差异，假象一只口渴的鹦鹉拿着一个杯子向一只聪明的乌鸦寻求一点
帮助，而乌鸦只能用桶子来尝试从一口水井中得到更多的水，显然这可以填满很多的杯
子，但是一杯就足以。显然鹦鹉就是更上一层次的调用，而对此需要提出自己需要什么
并尝试提供容器，而乌鸦则是usecase层，显然乌鸦无法确定这个杯子需要装什么，因
此为了尽可能地尝试让鹦鹉获得自己需要的，乌鸦显然更多时候需要不同的容器来从深
一层的地方获取，最终将这些东西组合。

```go
// 使用golang大致的usecase
package usecase

import (
    ...
)

type [T->Usecase] interface {}

type [t->Usecase] struct {
    repo [T->Repository]
    log [T->Logger]
    c context.Context
    t time.Duration
}

func New[T->Usecase](repo [T->Repository]) [T->Usecase] {
    return &[t->Usecase]{
        repo: repo,
        log: repo.<_>(),
        c: context.Background(),
        t: 3*time.Second,
    }
}

func (u *[t->Usecase]) <_>() [T->Logger] {
    return u.log
}

func (u *[t->Usecase]) <[motion]>(container...) error { 
    ctx,cancle := context.<_>(...)
    defer cancle()
    err := (u.repo.<_>() (format $container))
    with err (format err #return)
    with container (format $container)
}

// usecase层错误大致
type UError struct { 
    Code int 
    Message string
}
```
> <a id="handler">handler</a>

handler本身作为一个预处理，将来自client校验后的参数解析格式为request，在尝试
调用usecase后拿到格式后的数据转换为response，一个合理的设计是req和res作为
usecase层的接口规范，这避免了不必要的数据格式设计。handler仍然需要将所有可能
的数据格式为usecase可操作的，同时，handler层的返回一个合理的操作是通过最后的
统一错误接口返回，这显得很完美。

```go
// 使用golang大致的handler
package handler

import (
    ...
) 

type [T->Handler] interface {}

type [t->Handler] struct {
    ucase [T->Usecase]
    log [T->Logger]
    res [T->Response]
}

func New[T->Handler](ucase [T->Usecase],res [T->Response]) [T->Handler] {
    return &[t->Handler]{
        ucase: ucase,
        log: ucase.<_>(),
        res: res,
    } 
}

func (h *[t->Handler]) <[motion]>(c ^.Context) error {
    [motion->Req],ok := c.<_>($key).(*[motion->Req])
    with !ok (format err #return)
    [motion->Res] := *[motion->Res]{}
    err := (h.ucase.<_>() $[motion->Req] $[motion->Res])
    with err!=nil ($err #return)
}

// handler层错误大致
type HError struct { 
    Code int
    Message string
}
```

> <a id="route">route</a>

route层用于将端点和handler入口绑定，导向请求到正确的处理流程。

```go
// 使用golang大致实现handler
package api

const (
    [prefix] = ...
)

func RouteSetup() {
    for (register[endpoint->Route]() #log)
}

func register[endpoint->Route](dep...) {
    g := (group $[prefix])

    repo := (New[T->Repository]() $[dep.db] $[dep.log])
    ucase := (New[T->Usecase]() $repo)
    handler := (New[T->Handler]() $ucase) 
    for (
        with public g.<_>($[dep.middleware])
        with private g.<_>($[dep.middleware])
        g.<method>() $[endpoint] handler.<_>()
        )
}
```
