package web

import (
	"log"
	"net/http"
	"sync"
)

//TODO:
//同时在多个端口开启不同服务（√）
//自动生成短网址（√）
//静态文件服务（√）
//群组路由（√）
//将302重定向改为301
//要不要定义接口

//常量，用作下标
const (
	GET        = iota    // 0
	HEAD                 // 1
	POST                 // 2
	PUT                  // 3
	DELETE               // 4
	TRACE                // 5
	OPTIONS              // 6
	CONNECT              // 7
	PATCH                // 8
	ROUTERSIZE = 1 << 10 // 1024
	DEBUG      = true
)

//TODO:
//取消全局路由选择器个数限制
//全局路由选择器
var routers = make([]*router, 0, ROUTERSIZE)

//Method表
type reqs [9]handler

//Request Method，每个字段代表一种Method
type RM struct {
	GET     handler
	HEAD    handler
	POST    handler
	PUT     handler
	DELETE  handler
	TRACE   handler
	OPTIONS handler
	CONNECT handler
	PATCH   handler
}

//Map of RM，在群组路由中用来集体注册
type MRM map[string]RM

//处理器函数
type handler func(*Context)

func (f handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(&Context{w, r, nil})
}

//将http包的处理器函数转化为handler类型处理器
func FromHttpFunc(h func(http.ResponseWriter, *http.Request)) handler {
	return func(c *Context) {
		h(c.w, c.r)
	}
}

//将http包的处理器转化为handler类型处理器
func FromHttpHandler(h http.Handler) handler {
	return func(c *Context) {
		h.ServeHTTP(c.w, c.r)
	}
}

//将handler转化为http包的处理器函数
func ToHttpFunc(h handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(&Context{w, r, nil})
	}
}

//将handler转化为http包的处理器
func ToHttpHandler(h handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(&Context{w, r, nil})
	})
}

//重定向处理器
func RedirectHandler(url string, code int) handler {
	return func(c *Context) {
		http.Redirect(c.w, c.r, url, code)
	}
}

//404处理器
func NotFoundHandler(c *Context) {
	http.NotFound(c.w, c.r)
}

//501处理器
func NotImplementHandler(c *Context) {
	http.Error(c.w, "Not implement", http.StatusNotImplemented)
}

//处理error
type errorHandler func(error) error

//默认错误处理函数
func defaultErrorHandler(err error) error {
	log.Fatalln(err)
	return nil
}

//用于传递数据
type M map[string]interface{}

//TODO:
//加上读写锁保护
//路由选择器
type router struct {
	lock       sync.RWMutex
	mux        map[string]*reqs
	addr       string
	processerr errorHandler
}

//路由注册函数
func (r *router) Regist(url string, rm RM) {
	if rm.GET != nil {
		r.Get(url, rm.GET)
	}
	if rm.HEAD != nil {
		r.Head(url, rm.HEAD)
	}
	if rm.POST != nil {
		r.Post(url, rm.POST)
	}
	if rm.PUT != nil {
		r.Put(url, rm.PUT)
	}
	if rm.DELETE != nil {
		r.Delete(url, rm.DELETE)
	}
	if rm.TRACE != nil {
		r.Trace(url, rm.TRACE)
	}
	if rm.OPTIONS != nil {
		r.Options(url, rm.OPTIONS)
	}
	if rm.CONNECT != nil {
		r.Connect(url, rm.CONNECT)
	}
	if rm.PATCH != nil {
		r.Patch(url, rm.PATCH)
	}
}

//注册路由Get方法
func (r *router) Get(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}

	if v, ok := r.mux[url]; ok {
		(*v)[GET] = h
	} else {
		r.mux[url] = &reqs{GET: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{GET: RedirectHandler(url, 302)}
		} else if v[GET] == nil {
			(*v)[GET] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Head方法
func (r *router) Head(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[HEAD] = h
	} else {
		r.mux[url] = &reqs{HEAD: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{HEAD: RedirectHandler(url, 302)}
		} else if v[HEAD] == nil {
			(*v)[HEAD] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Post方法
func (r *router) Post(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[POST] = h
	} else {
		r.mux[url] = &reqs{POST: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{POST: RedirectHandler(url, 302)}
		} else if v[POST] == nil {
			(*v)[POST] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Put方法
func (r *router) Put(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[PUT] = h
	} else {
		r.mux[url] = &reqs{PUT: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{PUT: RedirectHandler(url, 302)}
		} else if v[PUT] == nil {
			(*v)[PUT] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Delete方法
func (r *router) Delete(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[DELETE] = h
	} else {
		r.mux[url] = &reqs{DELETE: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{DELETE: RedirectHandler(url, 302)}
		} else if v[DELETE] == nil {
			(*v)[DELETE] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Trace方法
func (r *router) Trace(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[TRACE] = h
	} else {
		r.mux[url] = &reqs{TRACE: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{TRACE: RedirectHandler(url, 302)}
		} else if v[TRACE] == nil {
			(*v)[TRACE] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Option方法
func (r *router) Options(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[OPTIONS] = h
	} else {
		r.mux[url] = &reqs{OPTIONS: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{OPTIONS: RedirectHandler(url, 302)}
		} else if v[OPTIONS] == nil {
			(*v)[OPTIONS] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Connect方法
func (r *router) Connect(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[CONNECT] = h
	} else {
		r.mux[url] = &reqs{CONNECT: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{CONNECT: RedirectHandler(url, 302)}
		} else if v[CONNECT] == nil {
			(*v)[CONNECT] = RedirectHandler(url, 302)
		}
	}
	return r
}

//注册路由Patch方法
func (r *router) Patch(url string, h handler) *router {
	if url == "" {
		panic("url or handler can't be empty.")
	}
	if h == nil {
		h = NotImplementHandler
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if r.mux == nil {
		r.mux = make(map[string]*reqs)
	}
	if v, ok := r.mux[url]; ok {
		(*v)[PATCH] = h
	} else {
		r.mux[url] = &reqs{PATCH: h}
	}

	if url[len(url)-1] == '/' {
		relateUrl := url[0 : len(url)-1]
		if v, ok := r.mux[relateUrl]; !ok {
			r.mux[relateUrl] = &reqs{PATCH: RedirectHandler(url, 302)}
		} else if v[PATCH] == nil {
			(*v)[PATCH] = RedirectHandler(url, 302)
		}
	}
	return r
}

/*-+--+--+--+--+-短网址路由-+--+--+--+--+-*/

//注册路由Get方法,同时注册短网址
func (r *router) ShortGet(url string, h handler) *router {
	r.Get(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Get(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Head方法,同时注册短网址
func (r *router) ShortHead(url string, h handler) *router {
	r.Head(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Head(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Post方法,同时注册短网址
func (r *router) ShortPost(url string, h handler) *router {
	r.Post(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Post(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Put方法,同时注册短网址
func (r *router) ShortPut(url string, h handler) *router {
	r.Put(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Put(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Delete方法,同时注册短网址
func (r *router) ShortDelete(url string, h handler) *router {
	r.Delete(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Delete(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Trace方法,同时注册短网址
func (r *router) ShortTrace(url string, h handler) *router {
	r.Trace(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Trace(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Option方法,同时注册短网址
func (r *router) ShortOptions(url string, h handler) *router {
	r.Options(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Options(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Connect方法,同时注册短网址
func (r *router) ShortConnect(url string, h handler) *router {
	r.Connect(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Connect(surl, RedirectHandler(addr, 302))
	return r
}

//注册路由Patch方法,同时注册短网址
func (r *router) ShortPatch(url string, h handler) *router {
	r.Patch(url, h)
	addr := "http://localhost" + r.addr + url
	surl := ShortenUrl(addr)
	shortRouter.Patch(surl, RedirectHandler(addr, 302))
	return r
}

//群组路由
func (r *router) Group(root string) *group {
	if root == "" {
		root = "/"
	} else if root[0] != '/' {
		root = "/" + root
	}
	n := len(root)
	if n > 1 && root[n-1] == '/' {
		root = root[0 : n-1]
	}
	return &group{root, r}
}

//中间件
func (r *router) Middleware(mw middlewareFunc, data M) *middleware {
	return &middleware{mw, data, r}
}

//静态文件服务
func (r *router) StaticFS(url, path string) {
	n, relativePath := len(url), url
	if url[n-1] == '/' {
		relativePath = url[0 : len(url)-1]
	} else {
		url += "/"
	}

	h := http.StripPrefix(relativePath, http.FileServer(http.Dir(path)))

	r.Get(url, FromHttpHandler(h))
	//r.Head(url, FromHttpHandler(h))
	r.Get(relativePath, RedirectHandler(url, 302))
	//r.Head(relativePath, FromHttpHandler(http.RedirectHandler(url, 302)))
}

//注册error处理函数
func (r *router) WhenError(h errorHandler) {
	r.processerr = h
}

//根据请求url和方法获取处理器
func (r *router) getHandler(rq *http.Request) (h handler) {
	url := rq.URL.Path
	index := method(rq.Method)
	if hs, ok := r.mux[url]; ok {
		if h = hs[index]; h == nil {
			h = NotFoundHandler
		}
		return
	}
	//路由匹配，层级下降
	n := 0
	h = NotFoundHandler
	for k, v := range r.mux {
		if match(url, k) && (h == nil || len(k) > n) {
			h = v[index]
			n = len(k)
		}
	}
	return
}

//TODO:
//将Connect方法区别对待
func (r *router) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	//log.Printf("[%-7s] - %s\n", rq.Method, rq.URL)
	if rq.RequestURI == "*" {
		if rq.ProtoAtLeast(1, 1) {
			rw.Header().Set("Connection", "close")
		}
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if DEBUG {
		log.Printf("[%s] - [%s%s] %s (%d)\n", rq.Method, rq.Host, rq.URL, rq.Proto, rq.ContentLength)
	}
	h := r.getHandler(rq)
	h(&Context{rw, rq, nil})
}

//启动Http服务
func (r *router) Run(addr ...string) error {
	if len(addr) > 0 && addr[0] != "" {
		r.addr = addr[0]
	}
	server := &http.Server{
		Addr:    r.addr,
		Handler: r,
	}
	return server.ListenAndServe()
}

//由路由选择器自己处理error
func (r *router) RunFree(addr ...string) {
	if len(addr) > 0 {
		r.addr = addr[0]
	}
	server := http.Server{
		Addr:    r.addr,
		Handler: r,
	}
	r.processerr(server.ListenAndServe())
}

//创建路由选择器
func NewRouter(addr ...string) *router {
	defaddr := ":8080"
	if len(addr) > 0 && addr[0] != "" {
		defaddr = addr[0]
	}
	if len(defaddr) > 0 && defaddr[0] != ':' {
		defaddr = ":" + defaddr
	}

	r := &router{
		mux:        nil,
		addr:       defaddr,
		processerr: defaultErrorHandler,
	}
	if len(routers) < ROUTERSIZE {
		routers = append(routers, r)
	} else {
		panic("you have created too many routers, should <= 1024")
	}
	return r
}

//由Method得到处理器地址
func method(m string) int {
	switch m {
	case "GET":
		return 0
	case "HEAD":
		return 1
	case "POST":
		return 2
	case "PUT":
		return 3
	case "DELETE":
		return 4
	case "TRACE":
		return 5
	case "OPTIONS":
		return 6
	case "CONNECT":
		return 7
	case "PATCH":
		return 8
	default:
		return -1
	}
}

//路由匹配
func match(url, path string) bool {
	n := len(path)
	if path[n-1] != '/' {
		return url == path
	}
	return len(url) > n && url[0:n] == path
}

//启动Http集群服务
func Run() {
	if len(routers) == 0 {
		log.Println("no router to run")
	} else if len(routers) > 1024 {
		panic("maybe too many routers(should ≤ 1024)")
	} else {
		errchan := make(chan error)
		for i := range routers {
			go func(r *router, ch chan<- error) {
				ch <- r.processerr(r.Run())
			}(routers[i], errchan)
		}
		for i := 0; i < len(routers); i++ {
			<-errchan
		}
	}
}

/*-+--+--+--+--+--+--+--+--+--+-群组路由-+--+--+--+--+--+--+--+--+--+-*/

//群组路由注册器
type group struct {
	root string
	r    *router
}

//群组理由下还可以注册群组路由
func (g *group) Group(url string) *group {
	if url == "/" {
		url = ""
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if url[len(url)-1] == '/' {
		url = url[0 : len(url)-1]
	}
	return &group{g.root + url, g.r}
}

//群组路由注册
func (g *group) Regist(mrm MRM) *group {
	for k, v := range mrm {
		fullurl := g.root
		if k[0] != '/' {
			fullurl = fullurl + "/" + k
		} else {
			fullurl += k
		}
		g.r.Regist(fullurl, v)
	}
	return g
}

//在group路由下注册子路由的Get方法
func (g *group) Get(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Get(fullurl, h)
	return g
}

//在group路由下注册子路由的Head方法
func (g *group) Head(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Head(fullurl, h)
	return g
}

//在group路由下注册子路由的Post方法
func (g *group) Post(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Post(fullurl, h)
	return g
}

//在group路由下注册子路由的Put方法
func (g *group) Put(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Put(fullurl, h)
	return g
}

//在group路由下注册子路由的Delete方法
func (g *group) Delete(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Delete(fullurl, h)
	return g
}

//在group路由下注册子路由的Trace方法
func (g *group) Trace(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Trace(fullurl, h)
	return g
}

//在group路由下注册子路由的Option方法
func (g *group) Options(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Options(fullurl, h)
	return g
}

//在group路由下注册子路由的Connect方法
func (g *group) Connect(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Connect(fullurl, h)
	return g
}

//在group路由下注册子路由的Patch方法
func (g *group) Patch(url string, h handler) *group {
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	g.r.Patch(fullurl, h)
	return g
}

/*-+--+--+--+--+-短网址路由-+--+--+--+--+-*/

//注册群组路由Get方法,同时注册短网址
func (g *group) ShortGet(url string, h handler) *group {
	g.Get(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Get(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Head方法,同时注册短网址
func (g *group) ShortHead(url string, h handler) *group {
	g.Head(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Head(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Post方法,同时注册短网址
func (g *group) ShortPost(url string, h handler) *group {
	g.Post(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Post(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Put方法,同时注册短网址
func (g *group) ShortPut(url string, h handler) *group {
	g.Put(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Put(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Delete方法,同时注册短网址
func (g *group) ShortDelete(url string, h handler) *group {
	g.Delete(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Delete(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Trace方法,同时注册短网址
func (g *group) ShortTrace(url string, h handler) *group {
	g.Trace(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Trace(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Options方法,同时注册短网址
func (g *group) ShortOptions(url string, h handler) *group {
	g.Options(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Options(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Connect方法,同时注册短网址
func (g *group) ShortConnect(url string, h handler) *group {
	g.Connect(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Connect(surl, RedirectHandler(addr, 302))
	return g
}

//注册群组路由Patch方法,同时注册短网址
func (g *group) ShortPatch(url string, h handler) *group {
	g.Patch(url, h)
	if url[0] != '/' {
		url = "/" + url
	}
	fullurl := g.root + url
	addr := "http://localhost" + g.r.addr + fullurl
	surl := ShortenUrl(addr)
	shortRouter.Patch(surl, RedirectHandler(addr, 302))
	return g
}
