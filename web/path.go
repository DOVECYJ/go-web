package web

//TODO:
//实现更多的路径过滤以及合法性检查功能

func cleanURL(url string) string {
	if url == "" {
		return "/"
	}
	if url[0] != '/' {
		url = "/" + url
	}
	if url[len(url)-1] == '/' {
		url = url[0 : len(url)-1]
	}
	return url
}
