# Thinking in Echo

## 关于Context

从使用者的角度来说，echo.Context会是最先接触到的一个概念，故先从这个下手，echo.Context接口的设计意图是将http的Request和Response封装组合在一起，形成一个Context，这么做的意图是显而易见的，因为Response的行为依赖与Request，将两者组合起来可以很好的表达出这种依赖性的概念。

```go
func helloWorld(c echo.Context) error {
    // echo.Context...
}
echo.GET("/", helloWorld)
```

### echo.Context 和 beego.Context

在beego中也有类似的Context概念，两者在具体的实现上有很大的不同，感觉beego.Context在结构上比echo.Context更清晰一些：

 - echo.Context是个巨型的接口，将Request和Response的方法集定义全部堆在一起，然后用一个context结构体来具体的实现出来，这一点感觉有点乱。
 - beego.Context是个结构体，其内部组合了几个职责更明确的子结构体来处理具体事宜，结构上更清晰。
 - 有利有弊，echo.Context概念单一，对于使用者来说，只需要知道一个Context足以，而beego.Context则需要知道更多一些的概念，从这个角度来说，echo.Context的设计有点类似Facade的感觉。  

### node.js中的Context？
node.js中并没有类似上述所说的Context概念:
```js
function helloWorld(req, res) {
    // req, res...
}
app.get("/", helloWorld)
```
从上述代码片段可以直观感受到，就算没有Context, 使用起来一样很简洁，[go-web](https://github.com/SkylakeCoder/go-web "") 项目中的做法就是模仿的node.js。

---
## 关于接口的用户体验
