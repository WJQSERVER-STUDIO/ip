package main

import (
	"log"
	"net/http"

	"github.com/WJQSERVER-STUDIO/ip/logger"
	"github.com/WJQSERVER-STUDIO/ip/lookup"
	"github.com/WJQSERVER-STUDIO/ip/proxy"
)

func init() {
	// 初始化日志记录器，传入日志文件路径
	logger.Init("/data/ipinfo/log/access.log")

	// 初始化数据库
	lookup.Init()
}

func main() {
	defer logger.Close() // 确保在程序结束时关闭日志文件
	// 设置HTTP路由处理器并启动服务器
	http.HandleFunc("/", LogRequestWrapper(func(w http.ResponseWriter, r *http.Request) {
		// 其他处理逻辑...
	}))
	http.HandleFunc("/ip-lookup", LogRequestWrapper(lookup.IPLookupHandler))
	http.HandleFunc("/ip", LogRequestWrapper(lookup.GetIPHandler))
	http.HandleFunc("/bilibili", LogRequestWrapper(proxy.BilibiliHandlerWithChromeTLS))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// LogRequestWrapper 包装日志记录功能
func LogRequestWrapper(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.LogRequest(r)
		handler(w, r)
	}
}
