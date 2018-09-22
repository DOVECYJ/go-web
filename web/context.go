package web

import (
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
)

//http://user:password@localhost:8080/hello/world?name=cyj&location=beijing#fragment

//文件上传大小限制，1G
var maxUploadFileSize int64 = 1 << 30

type Context struct {
	w    http.ResponseWriter
	r    *http.Request
	data map[string]interface{}
}

//TODO:
//remove this function
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:    w,
		r:    r,
		data: nil,
	}
}

//输出字符串
func (c *Context) WriteString(s string) {
	s += "\n"
	c.w.Write([]byte(s))
}
func (c *Context) Write(a ...interface{}) (int, error) {
	return fmt.Fprint(c.w, a...)
}

func (c *Context) Writeln(a ...interface{}) (int, error) {
	return fmt.Fprintln(c.w, a...)
}

func (c *Context) Writef(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(c.w, format, a...)
}

//重定向
func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.w, c.r, url, code)
}

//渲染Html模板
func (c *Context) Html(path string) {
	t := template.Must(template.ParseFiles(path))
	t.Execute(c.w, c.data)
}

//渲染string类型的Html模板
func (c *Context) SHtml(name, html string) {
	t := template.New(name)
	t = template.Must(t.Parse(html))
	t.Execute(c.w, c.data)
}

//TODO:
//彻底封装http.Request，去掉Request()函数

//返回http.Request
func (c *Context) Request() *http.Request {
	return c.r
}

func (c *Context) Param(key string) string {
	if err := c.r.ParseForm(); err != nil {
		panic("parse form error")
	}
	if v, ok := c.r.Form[key]; ok {
		return v[0]
	}
	panic("no value match key")
}

func (c *Context) Params(key string) []string {
	if err := c.r.ParseForm(); err != nil {
		panic("parse form error")
	}
	if v, ok := c.r.Form[key]; ok {
		return v
	}
	panic("no value match key")
}

func (c *Context) File(key string) (multipart.File, *multipart.FileHeader) {
	file, header, err := c.r.FormFile(key)
	if err != nil {
		panic("parse file error")
	}
	return file, header
}

func (c *Context) Files(key string) []*multipart.FileHeader {
	if err := c.r.ParseMultipartForm(maxUploadFileSize); err != nil {
		panic("parse multipart form error")
	}
	files := c.r.MultipartForm.File[key]
	return files
}

func (c *Context) MustParse() {
	if err := c.r.ParseForm(); err != nil {
		panic("parse form error")
	}
	if err := c.r.ParseMultipartForm(maxUploadFileSize); err != nil {
		panic("parse multipart form error")
	}
}

func (c *Context) SetCookie(name, value string) {
	cook := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
	}
	http.SetCookie(c.w, cook)
}

func (c *Context) Cookie(key string) (cook *http.Cookie, err error) {
	cook, err = c.r.Cookie(key)
	return
}

func (c *Context) Cookies() []*http.Cookie {
	return c.r.Cookies()
}

//向上下文中存储数据
func (c *Context) SetData(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if _, ok := c.data[key]; !ok {
		c.data[key] = value
	} else {
		panic("repeated key")
	}
}

//设置文件上传大小限制
func (c *Context) SetUploadSize(n int64) {
	if n > 0 {
		maxUploadFileSize = n
	} else {
		maxUploadFileSize = 1 << 30
	}
}

//for test
const Temp = `
<html>
	<head>
		<title>INDEX</title>
	</head>
	
	<body>
		<h3>你好: </h3>
		<h1>{{.name}}</h1>
	</body>
</html>
`
