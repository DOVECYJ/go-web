package web

import (
	//"encoding/base64"
	"io"
	"log"
	"net/http/httptest"
	"time"
)

//自定义中间件上下文
type MContext struct {
	h    handler
	Ctx  *Context
	Data M
}

//执行响应
func (m *MContext) DoHandler() {
	m.h(m.Ctx)
}

type middlewareFunc func(mc *MContext)

type middleware struct {
	ware middlewareFunc
	data M
	r    *router
}

//注册Get方法处理器的中间件
func (m *middleware) Get(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Get(url, fn)
	return m
}

//注册Head方法处理器的中间件
func (m *middleware) Head(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Head(url, fn)
	return m
}

//注册Post方法处理器的中间件
func (m *middleware) Post(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Post(url, fn)
	return m
}

//注册Put方法处理器的中间件
func (m *middleware) Put(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Put(url, fn)
	return m
}

//注册Delete方法处理器的中间件
func (m *middleware) Delete(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Delete(url, fn)
	return m
}

//注册Trace方法处理器的中间件
func (m *middleware) Trace(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Trace(url, fn)
	return m
}

//注册Options方法处理器的中间件
func (m *middleware) Options(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Options(url, fn)
	return m
}

//注册Connect方法处理器的中间件
func (m *middleware) Connect(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Connect(url, fn)
	return m
}

//注册Patch方法处理器的中间件
func (m *middleware) Patch(url string, h handler) *middleware {
	fn := func(c *Context) {
		m.ware(&MContext{h, c, m.data})
	}
	m.r.Patch(url, fn)
	return m
}

//TODO:
//实现默认的基础中间件
//以下为普通中间件，直接注册

//响应特定主机
func ServeSpecialHost(h handler, host string) handler {
	return func(c *Context) {
		if c.r.Host == host {
			h(c)
		} else {
			c.w.WriteHeader(403)
		}
	}
}

//TODO:
//将数据缓存移到数据库
var logfmt = "[%s] - [%s] (%d) [%v] 访问次数(%d)\n"
var rcount = make(map[string]uint64)

//log 中间件
func Logg(h handler, w ...io.Writer) handler {
	if len(w) > 0 && w[0] != nil {
		log.SetOutput(w[0])
	}
	log.SetPrefix("[go-web]")
	return func(c *Context) {
		rq := c.r
		recod := httptest.NewRecorder()
		tstart := time.Now()
		h(&Context{recod, c.r, nil})
		for k, v := range recod.Header() {
			c.w.Header()[k] = v
		}
		c.w.Write(recod.Body.Bytes())
		tend := time.Now()

		addr := rq.Host + rq.URL.String()
		rcount[addr]++
		log.Printf(logfmt, rq.Method, addr, recod.Code, tend.Sub(tstart), rcount[addr])
	}
}

//追加自定义响应内容
func AppendResponse(h handler, msg string) handler {
	return func(c *Context) {
		h(c)
		c.WriteString(msg)
	}
}

//基础认证
func BasicAuth(h handler, user, password string) handler {
	return func(c *Context) {
		u, p, ok := c.r.BasicAuth()
		if !ok {
			//TODO:
			//弹出密码框
			c.w.WriteHeader(401)
			return
		}
		if !checkauth(user, password, u, p) {
			c.w.WriteHeader(403)
			return
		}
		h(c)
	}
}

func checkauth(user, password, inputuser, inputpassword string) bool {
	if user != inputuser && password == inputpassword {
		println("cheak user")
		return false
	}
	//p := base64.URLEncoding.EncodeToString([]byte(password))
	if password != inputpassword {
		return false
	}
	return true
}
