package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// 设置目标URL
	target, err := url.Parse("http://192.168.1.18:11434")
	if err != nil {
		log.Fatal("解析目标URL失败:", err)
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("收到响应状态: %d", resp.StatusCode)
		return nil
	}

	// 创建处理函数
	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("收到请求: %s %s", r.Method, r.URL.Path)
		log.Printf("转发请求到: %s%s", target.Host, r.URL.Path)
		proxy.ServeHTTP(w, r)
	}

	// 设置监听地址和处理函数
	http.HandleFunc("/", handler)

	log.Printf("反向代理启动在 http://127.0.0.1:11434")
	log.Printf("将请求转发到 http://192.168.1.18:11434")

	// 启动服务器
	if err := http.ListenAndServe(":11434", nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
