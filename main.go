package main

import (
	"flag"
	"fmt"
	"log"

	"ip/api"
	"ip/config"
	"ip/db"
	"ip/logger"

	"github.com/gin-gonic/gin"
)

var (
	cfg        *config.Config
	configfile = "/data/go/config/config.toml"
	router     *gin.Engine
)

// 日志模块
var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	logWarning = logger.LogWarning
	logError   = logger.LogError
)

func ReadFlag() {
	cfgfile := flag.String("cfg", configfile, "config file path")
	flag.Parse()
	configfile = *cfgfile
}

func loadConfig() {
	var err error
	// 初始化配置
	cfg, err = config.LoadConfig(configfile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Loaded config: %v\n", cfg)
}

func setupLogger() {
	// 初始化日志模块
	var err error
	err = logger.Init(cfg.Log.LogFilePath, cfg.Log.MaxLogSize) // 传递日志文件路径
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logw("Logger initialized")
	logw("Init Completed")
}

func updateDB() {
	go db.LoopForUpdate(cfg)
}

func setupDB() {
	db.DBinit(cfg)
}

func setupApi(cfg *config.Config, router *gin.Engine) {
	api.InitHandleRouter(cfg, router)
}

func init() {
	ReadFlag()
	loadConfig()
	setupLogger()
	updateDB()
	setupDB()

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.UseH2C = false
	setupApi(cfg, router)
}

func main() {
	err := router.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		logError("Failed to start server: %v\n", err)
	}
	defer logger.Close() // 确保在退出时关闭日志文件
}
