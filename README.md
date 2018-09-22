# go-web
a http web framework with golang

特性：

链式路由注册

群组路由注册

短网址服务

中间件支持

兼容go原生net/http包

同时在多个端口开启服务

---

说明：

目前仍在开发中，有些功能还不完善，都在代码中通过TODO标注出来了。

---

## 使用示例

首先需要下载并导入`go-web`包。

```
go get github.com/DOVECYJ/go-web
```

``` go
import github.com/DOVECYJ/go-web/web
```

### 快速开始

处理器是一个函数类型，只要是符合`func(*Context)`类型签名的函数是合法的处理器。

``` go
func hello(c *web.Context) {
	c.WriteString("hello coder")
}
```

通过`NewRouter`函数新建路由选择器，通过`Get`函数可以注册一个`GET`方法的路由，并通过`Run`方法就可快速开启服务。

``` go
func main() {
	r1 := web.NewRouter(":8080")
	r1.Get("/hello", hello)
	r1.Run()
}
```
另一种方式是一起注册很多路由。

``` go
r1.Regist("/hello",  web.RM{
	  GET:  hello,
	  POST:  world, //world是一个处理器，请自行实现
})
```

### 链式路由

路由选择器的每个方法都支持链式路由注册。

``` go
r1.Get("/hello", hello).Post("/world", world)
```

上面的代码会注册两个独立的路由`/hello`和`/world`，而不是`/hello`和`/hello/world`，这是群组路由的功能。对于只需要注册`GET`和`POST`路由的情况，使用链式路由可以不用另起一行。

### 群组路由

在使用上和简单的路由没有太大区别。

``` go
r1.Group("/prefix").Get("/hello", hello)
```

实际注册的是`/prefix/hello`路由。`Group`之下还可以再注册群组路由。

``` go
r1.Group("prefix").Group("sub").Get("hello", hello)
```

实际注册的是`/prefix/sub/hello`。

如果`Group`下有大量路由需要注册，还有一种更快捷的方式。

``` go
r1.Group("/prefix").Regist(web.MRM{
	  "/hello":  web.RM{GET:  hello,  POST:  hello_post},
	  "/world":  web.RM{GET:  world},
	  "/hi":	 web.RM{Post: hi},
})
```

上面的代码将注册三个路由：
`/prefix/hello`
`/prefix/world`
`/prefix/hi`

### 注册到`net/http`

路由选择器本身也是http处理器，可知直接兼容`net/http`包。

``` go
http.ListenAndServe(":8080", r1)
```

在这种情况下，`r1`指定的短端口将失效。

### 开启集群服务

开启一个服务集群的方法很简单，只需要创建多个`router`，注册好路由，然后调用`web.Run`就可以了。

``` go
r1 := web.NewRouter(":8080").Get("/hello", hello)
r2 := web.NewRouter(":8081").Get("/world", world)
r3 := web.NewRouter(":8082").Get("/hi", hi)
web.Run()
```

运行`web.Run`之后，`:8080`、`:8081`和`:8082`端口上将同时开启服务。

### 自定义中间件

你可以将自定义中间件包装成一个处理器，然后像使用处理器那样使用中间件。

``` go
func mid(h func(*web.Context), msg string) func(*web.Context) {
	return func(c *web.Context) {
		c.WriteString(msg)
		h(c)
	}
}

func main() {
	r1 := web.NewRouter(":8080")
	r1.Get("/hello", mid(hello, "我是中间件"))
}
```

或者，你也可以用中间件函数，中间件函数是函数签名符合`func(*web.MContext)`的函数。

``` go
func mid(mc  *web.MContext)  {
	mc.Ctx.WriteString(mc.Data["midware"])
	mc.DoHandler()
}

func main() {
	r1 := web.NewRouter(":8080")
	r1.Middleware(mid, web.M{"midware":"我是中间件"}).Get("/hello", hello)
}
```

### 短网址

启用短网址服务必须使用`web.Run`开启集群服务，而且在这之前需要先调用`web.EnableShortUrlServe()`，因为短网址是一个单独的服务，运行在`:8764`端口，当然你也可以通过`SetShortServerPort`函数更改。

注册短网址需要使用带`Short-`的方法。

``` go
r1.ShortGet("/long/long/hello",  hello)
```

使用如下函数可以查询注册的短网址。

``` go
fmt.Println("短网址：", web.GetShort("localhost:8080/long/long/hello"))
```
在浏览器访问`localhost:8764/qAvARb`就可以看到"hello coder"输出了。

---
