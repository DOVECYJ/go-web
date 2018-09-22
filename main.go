// go-web project main.go
package main

import (
	"fmt"
	"go-web/web"
	"log"
	"net/http"
)

func hello(c *web.Context) {
	c.WriteString("<h1>Hello</h1>")
}

func hi(c *web.Context) {
	c.WriteString("<h1>Hi</h1>")
}

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

func errhandler(err error) error {
	fmt.Println("[err]", err)
	return nil
}

func main() {
	r1 := web.NewRouter()
	r1.Get("/hello", hi).Get("/notfound", nil).Get("/hello/", hello).StaticFS("/static", "./")

	r1.Group("/com").Get("/a", a).Get("/b", b).Get("/c", c)
	r1.Group("/gp").Regist(web.MRM{
		"/d": web.RM{GET: d},
		"/e": web.RM{GET: e},
		"/f": web.RM{GET: f},
	})
	r1.Group("/abc").Group("/123").Get("/g", g)

	r1.WhenError(errhandler)

	if err := http.ListenAndServe(":8080", r1); err != nil {
		log.Fatalln(err)
	}
}
