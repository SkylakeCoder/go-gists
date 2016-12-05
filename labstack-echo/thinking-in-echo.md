# Thinking in Echo

## 关于Context

从使用者的角度来看，echo.Context 会是最先接触到的概念，echo.Context接口的设计意图是将http的Request和Response封装在一起，形成一个Context，这么做的意图是显而易见的，因为Response的行为依赖于Request，将两者组合起来可以很好的表达出这种依赖性的概念。

```go
func helloWorld(c echo.Context) error {
    // echo.Context...
}
echo.GET("/", helloWorld)
```

### echo.Context 和 beego.Context

在beego中也有类似的Context概念，两者在具体的实现上有很大的不同，感觉beego.Context在结构上比echo.Context更清晰一些：

 - echo.Context是个巨型的接口，将Request和Response的方法集定义全部堆在一起，然后用一个context结构体来具体的实现出来，这一点感觉有点乱。
 - beego.Context是个结构体，其内部组合了几个职责更明确的子结构体来处理具体事宜，结构上更清晰一点。
 - 有利有弊，echo.Context概念单一，对于使用者来说，只需要知道一个Context足以，而beego.Context则需要知道更多一些的概念，从这个角度来说，echo.Context的设计有点类似Facade的感觉。  

### node.js中的Context？
node.js中并没有类似上述所说的Context概念:
```js
function helloWorld(req, res) {
    // req, res...
}
app.get("/", helloWorld)
```
从上述代码片段可以直观感受到，就算没有Context, 使用起来一样很简洁，我在 [go-web](https://github.com/SkylakeCoder/go-web "") 项目中的做法就是模仿的node.js。

---
## echo中的Middleware
说起中间这个概念，显然要有起始两端，不然中间一词从何说起呢？  
其中一端是http client发起的请求，另一端是http server对请求的处理，中间件的作用就是在这两端之间的某个环节中搞一搞。在我目前的理解看来，中间件只是一种划分模块的方式，可以让模块处在某个合适的层次中。  
<br/>
echo在对Middleware的处理上挺有意思的，用了golang中的``闭包``特性来实现多个Middlewares和最终处理函数间的``链式调用``，可以在```echo.ServeHTTP```方法中看到相关的处理。  
为了实现``链式调用``，echo需要一视同仁的对待Middlewares中的处理函数和用户自己定义的处理函数，在echo中，这两者都是echo.HandlerFunc类型。  
```go
MiddlewareFunc func(HandlerFunc) HandlerFunc
```
另外echo中还有个Premiddleware的概念，和Middleware的区别是作用的时机有所不同，Premiddlewares作用于Router匹配路由之前，Middlewares作用于之后。这样一来，就可以在Premiddlewares中做一些可以影响到路由规则的事情，比如处理HTTP请求的重定向。  
<br/>
最后贴一段代码，聊表敬意！闭包实现的调用链实在是清楚干脆。
```go
    // Middleware
    h := func(c Context) error {
        method := r.Method
        path := r.URL.Path
        e.router.Find(method, path, c)
        h := c.Handler()
        for i := len(e.middleware) - 1; i >= 0; i-- {
            h = e.middleware[i](h)
        }
        return h(c)
    }
    // Premiddleware
    for i := len(e.premiddleware) - 1; i >= 0; i-- {
        h = e.premiddleware[i](h)
    }
    // Execute chain
    if err := h(c); err != nil {
        e.HTTPErrorHandler(err, c)
    }
```

---
## echo中的路由处理
echo在router.go模块中实现路由功能，一个特色是echo用了```前缀树```结构来实现路由规则的注册和查询。

### echo选择前缀树的原因猜测？
之所以选择前缀树而不是直接用map(golang中的是哈希型map)，感觉上可能有如下两点原因：
- 前缀树比哈希型map更省空间？  
  <br/>
  在我现有的见识看来，web开发中的路由规则再多又能多到哪里去呢？一万条路由规则？好吧，就算如此，也根本不用去操心内存占用的问题。  
  <br/>
- 最坏情况下，前缀树的查找性能优于哈希型map？  
  <br/>
  理论上如此，但既然echo是使用前缀树来处理路由，那咱就考虑http路由的情景，在这种情景下，前缀树的查找性能真的比哈希型map要强？  
  实际测试一下吧，就测试前缀树 和 ```map[string]bool``` 这两者到底谁的查找速度更快，具体测试代码见[这里](https://github.com/SkylakeCoder/go-gists/tree/master/radixtree "")，先说结论，在我的测试中，前缀树的查找性能``远远不如``golang中的map。。。  
  另外在写测试时，我本想直接精简echo router.go中的代码，摘出一段干净能用的前缀树来，但在具体的实现中，操作前缀树的代码与处理路由的代码交织在一起，还有一些```goto```...最后只成功摘取了insert方法，测试中前缀树的Find方法我是自行实现的，但是不用担心效率问题，不太可能会比echo原有的实现差。

上述原因如果在http路由场景下都不适用的话，那是否可以认为echo在路由时选择前缀树并不是合理的选择？并且带来了明显的坏处，就是router.go的代码可读性变的很差。

### echo选择前缀树的真正原因？
我问了一朋友对这个问题的看法，他认为：使用前缀树可以更方便的处理url中的Pattern，在pattern复杂时，map就不一定比前缀树快了。  
应该是如此，echo的确是在树的遍历过程中同时处理了一些pattern，比如那个著名的 ``/user/:username`` pattern. 我在[go-web](https://github.com/SkylakeCoder/go-web)项目中也处理过这个pattern，我的处理方式没有echo快，但在结构上比echo要好，见这里：[PatternHandler](https://github.com/SkylakeCoder/go-web/blob/master/web/pattern.go)

### echo的路由总结
- 路由的注册  
echo在服务启动之前，会完成所有路由的注册工作。  
echo模块中有一系列的路由注册方法(GET POST ...), 实际的路由注册工作由router模块来完成，围绕着路由path来维护一棵字符串前缀树，这棵树中的节点类型是如下的node结构体：
    ```go
    node struct {
        kind          kind
        label         byte
        prefix        string
        parent        *node
        children      children
        ppath         string
        pnames        []string
        methodHandler *methodHandler
    }
    ```
    节点中的parent属性，router并没有去真正的使用，并且维护parent属性的代码中还有Bug，不过不会影响到什么。  
    节点中的methodHandler属性，只有叶子节点中的才是有效值，这个属性的作用是将method+path和具体的HandlerFunc关联起来。至此路由的注册工作就完成了。

- 路由的查找  
    查找是根据method+path,从前缀树中找到特定的node，然后用node中的数据来填充Context，具体的实现中遍历树和处理特殊路由规则的逻辑混杂在一起，读起来比较头疼，就不详细展开了。。。

---
## echo中使用的graceful.Server
echo用了graceful.Server, 没有直接去用原生的http.Server. 那graceful.Server比原生的好在哪呢？尚未调查过，[graceful](https://github.com/tylerb/graceful) 自己的说法是：graceful可以优雅的关掉服务。在echo中设置ShutdownTimeout然后通过graceful来起作用，另外从示例来看，graceful还可以监测连接的状态。

---
## echo中的Binder
TODO...
