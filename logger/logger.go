// logger/logger.go
package logger

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var logFile *os.File
var logger *log.Logger

// Init 初始化日志记录器，接受日志文件路径作为参数
func Init(logFilePath string) {
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

// LogRequest 记录请求日志
func LogRequest(r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	userAgent := r.UserAgent()
	dateTime := time.Now().Format("02/Jan/2006:15:04:05 -0700")

	logEntry := fmt.Sprintf("%s - - [%s] \"%s %s %s\" 200 0 \"-\" \"%s\"",
		ip, dateTime, r.Method, r.RequestURI, r.Proto, userAgent)

	logger.Println(logEntry)
}

// Close 关闭日志文件
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
