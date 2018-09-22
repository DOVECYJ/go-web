package web

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"strings"
)

//编码表
var (
	table = [...]byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3',
		'4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D',
		'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
		'Y', 'Z',
	}
	//TODO:
	//将短网址存储到redis
	shorturlmap = make(map[string]string)
	shortRouter = &router{ //短网址服务器
		addr:       ":8764", //短网址默认服务端口
		mux:        make(map[string]*reqs),
		processerr: defaultErrorHandler,
	}
)

//生成短网址，未使用
func shortenurl(addr string) string {
	sum := md5.Sum([]byte(addr))
	n := rand.Intn(4) * 4
	part := binary.BigEndian.Uint32(sum[n : n+4])
	buf := bytes.NewBuffer([]byte("/"))
	for i := 0; i < 6; i++ {
		index := 0x3D & part
		buf.WriteByte(table[index])
		part >>= 5
	}
	return buf.String()
}

//生成短网址
func ShortenUrl(addr string) string {
	if v, ok := shorturlmap[addr]; ok {
		return v
	}
	sum, surl := md5.Sum([]byte(addr)), ""
	for i := 0; i < 16; i = i + 4 {
		part := binary.BigEndian.Uint32(sum[i : i+4])
		buf := bytes.NewBuffer([]byte("/"))
		for j := 0; j < 6; j++ {
			index := 0x3D & part //0x3D=61
			buf.WriteByte(table[index])
			part >>= 5
		}
		surl = buf.String()
		if !checkRepeat(surl) { //短网址没有重复
			shorturlmap[addr] = surl
			return surl
		}
	}
	//短网址重复：概率极低
	for i := 0; i < 10; i++ { //做10随机次尝试
		n := rand.Intn(12)
		part := binary.BigEndian.Uint32(sum[n : n+4])
		buf := bytes.NewBuffer([]byte("/"))
		for i := 0; i < 6; i++ {
			index := 0x3D & part
			buf.WriteByte(table[index])
			part >>= 5
		}
		surl = buf.String()
		if !checkRepeat(surl) { //短网址没有重复
			shorturlmap[addr] = surl
			return surl
		}
	}
	return surl //返回""表示失败
}

//检查短网址是否重复
func checkRepeat(url string) bool {
	for _, v := range shorturlmap {
		if url == v {
			return true
		}
	}
	return false
}

//返回短网址多路选择器
func ShortServer(addr ...string) *router {
	if len(addr) > 0 && addr[0] != "" {
		shortRouter.addr = addr[0]
	}
	return shortRouter
}

//TODO:
//自动补全网址
//根据url查询短网址
func GetShort(url string) string {
	if url == "" {
		return url
	}
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	if v, ok := shorturlmap[url]; ok {
		return v
	}
	return ""
}

//直接开启短网址服务
//不能与集群服务同时使用
//想只使用短网址服务可以用这个函数快速启动
func ServeShort(addr ...string) {
	if len(addr) > 0 && addr[0] != "" {
		shortRouter.addr = addr[0]
	}
	go shortRouter.Run()
}

//启用短网址服务，需要在集群服务启动前调用
func EnableShortUrlServe() {
	if len(routers) < ROUTERSIZE {
		routers = append(routers, shortRouter)
	} else {
		panic("can't enable short url serve beacuse there are too many routers")
	}
}

//更改短网址服务端口
func SetShortServerPort(port string) {
	shortRouter.addr = port
}
