package main

import (
	"fmt"
	"go-web/web"
)

func a(c *web.Context) {
	c.WriteString("<h1>Aa</h1>")
}

func b(c *web.Context) {
	c.WriteString("<h1>Bb</h1>")
}

func c(c *web.Context) {
	c.WriteString("<h1>Cc</h1>")
}

func d(c *web.Context) {
	c.WriteString("<h1>Dd</h1>")
}

func e(c *web.Context) {
	c.WriteString("<h1>Ee</h1>")
}

func f(c *web.Context) {
	c.WriteString("<h1>Ff</h1>")
}

func g(c *web.Context) {
	c.WriteString("<h1>Gg</h1>")
}

func h(c *web.Context) {
	c.WriteString("<h1>Hh</h1>")
}

func i(c *web.Context) {
	c.WriteString("<h1>Ii</h1>")
}

func j(c *web.Context) {
	c.WriteString("<h1>Jj</h1>")
}
func k(c *web.Context) {
	c.WriteString("<h1>Kk</h1>")
}

func l(c *web.Context) {
	c.WriteString("<h1>Ll</h1>")
}

func midw(mc *web.MContext) {
	mc.Ctx.WriteString("中间件")
	mc.DoHandler()
}

func main() {
	//基础路由
	r1 := web.NewRouter(":8080")
	r1.Get("/a", a)
	r1.Post("/b", b)
	r1.Regist("/c", web.RM{
		GET:  c,
		POST: c,
	})
	//链式路由
	r2 := web.NewRouter(":8081")
	r2.Get("/d", d).Post("/e", e)
	//群组路由
	r3 := web.NewRouter(":8082")
	r3.Group("/com").Regist(web.MRM{
		"/f": web.RM{GET: f, POST: f},
		"/i": web.RM{GET: i},
	})
	//短网址
	r4 := web.NewRouter(":8083")
	r4.ShortGet("/longlongurlj", j)
	fmt.Println("短网址：", web.GetShort("localhost:8083/longlongurlj"))
	//中间件
	r5 := web.NewRouter(":8084")
	r5.Get("/k", web.Logg(k))
	r5.Middleware(midw, nil).Get("/l", l)

	web.EnableShortUrlServe()
	web.Run()
}
