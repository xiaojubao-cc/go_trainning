package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	http.HandleFunc("/proxy", func(writer http.ResponseWriter, request *http.Request) {
		// 1. 设置目标 URL
		target, _ := url.Parse("https://www.baidu.com")

		// 2. 创建 Director 并修改请求信息
		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme // 设置 scheme: https
			req.URL.Host = target.Host     // 设置 host: www.baidu.com
			req.Host = target.Host         // 设置 Host header
		}

		// 3. 创建 ReverseProxy 实例
		proxy := &httputil.ReverseProxy{Director: director}

		// 4. 调用 ServeHTTP 开始代理请求
		proxy.ServeHTTP(writer, request)
	})
	http.ListenAndServe(":8080", nil)

}
