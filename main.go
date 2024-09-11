package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/WJQSERVER-STUDIO/ip/lookup"
	"github.com/WJQSERVER-STUDIO/ip/proxy"
)

var logger *log.Logger

func main() {
	// 初始化日志文件
	logFile, err := os.OpenFile("/data/ipinfo/log/access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	logger = log.New(logFile, "", log.LstdFlags)

	// 初始化数据库
	lookup.Init()

	// 设置HTTP路由处理器并启动服务器
	http.HandleFunc("/", LogRequestWrapper(func(w http.ResponseWriter, r *http.Request) {
		// 其他处理逻辑...
	}))
	http.HandleFunc("/ip-lookup", LogRequestWrapper(lookup.IPLookupHandler))
	http.HandleFunc("/ip", LogRequestWrapper(lookup.GetIPHandler))
	http.HandleFunc("/bilibili", LogRequestWrapper(proxy.BilibiliHandler))

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// LogRequestWrapper 包装日志记录功能
func LogRequestWrapper(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		handler(w, r)
	}
}

// 记录请求日志
func logRequest(r *http.Request) {
	// 获取请求的IP地址
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// 获取用户代理
	userAgent := r.UserAgent()

	// 获取日期信息
	dateTime := time.Now().Format("02/Jan/2006:15:04:05 -0700")

	// 格式化日志信息
	logEntry := fmt.Sprintf("%s - - [%s] \"%s %s %s\" 200 0 \"-\" \"%s\"",
		ip, dateTime, r.Method, r.RequestURI, r.Proto, userAgent)

	// 将日志写入文件
	logger.Println(logEntry)
}
