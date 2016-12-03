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
说起中间这个概念，显然要有起始两端，不然中间一词从何说起？  
其中一端是http client发起的请求，另一端是http server对请求的处理，中间件的作用就是在这两端之间的某个环节中搞一搞。在我目前的理解看来，中间件只是一种划分模块的方式，可以让模块处在某个合适的层次中。  
<br/>
echo在对Middleware的处理上挺有意思的，用了golang中的``闭包``特性来实现多个Middlewares和最终处理函数间的``链式调用``，可以在```echo.ServeHTTP```方法中看到具体的相关流程。  
<br/>
为了实现``链式调用``，echo需要一视同仁的对待Middlewares的处理函数和用户自己定义的处理函数，在echo中，这两者都是echo.HandlerFunc类型。  
<br/>
另外echo中还有个Premiddleware的概念，和Middleware的区别是作用的时机有所不同，Premiddlewares作用于Router匹配路由之前，而Middlewares作用于之后。这样一来，就可以在Premiddlewares中做一些可以影响到路由规则的事情，比如处理那个著名的```/user/:username```，哈哈哈。我的 [go-web](https://github.com/SkylakeCoder/go-web "") 项目中也处理过```/user/:username```，我当时的设计：[PatternHandler](https://github.com/SkylakeCoder/go-web/blob/master/web/pattern.go "") 机制也不算太差，但是和Premiddlewares机制相比，就不禁有种膝盖一软的感觉。  
<br/>
最后贴上一段代码，聊表敬意！
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

TODO...
