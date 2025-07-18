package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"ip/api"
	"ip/config"
	"ip/db"

	"github.com/fenthope/reco"
	"github.com/infinite-iroha/touka"
)

//go:embed pages
var pagesFS embed.FS

var (
	cfg        *config.Config
	configfile = "/data/ip/config/config.toml"
	router     *touka.Engine
)

func ReadFlag() {
	cfgfile := flag.String("c", configfile, "config file path")
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

func updateDB(cfg *config.Config) {
	go db.LoopForUpdate(cfg)
}

func setupDB(cfg *config.Config) {
	db.DBinit(cfg)
}

func init() {
	ReadFlag()
	loadConfig()

	router = touka.Default()
	api.InitHandleRouter(cfg, router)
}

func setLogger(cfg *config.Config) *reco.Logger {
	lcfg := reco.Config{
		Level:           reco.LevelInfo,            // 最低记录级别为 DEBUG
		Mode:            reco.ModeText,             // 输出 JSON 格式
		FilePath:        cfg.Log.LogFilePath,       // 日志文件路径
		EnableRotation:  true,                      // 启用文件轮转
		MaxFileSizeMB:   int64(cfg.Log.MaxLogSize), // 单个文件最大 5MB
		MaxBackups:      3,                         // 保留 3 个旧备份
		CompressBackups: true,                      // 压缩旧备份
		Async:           true,                      // 启用异步写入 (默认)
		BufferSize:      4096,                      // 异步缓冲区大小
	}

	// 创建 Logger 实例
	logger, err := reco.New(lcfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	return logger
}

func main() {
	logger := setLogger(cfg)
	db.LoggerInit(logger)
	updateDB(cfg)
	setupDB(cfg)
	router.SetLogger(logger)
	pFS, err := fs.Sub(pagesFS, "pages")
	if err != nil {
		logger.Errorf("Failed to load embedded pages: %v\n", err)
		return
	}
	router.SetUnMatchFS(http.FS(pFS))
	err = router.RunShutdown(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		logger.Errorf("Failed to start server: %v\n", err)
	}
	defer logger.Close() // 确保在退出时关闭日志文件
}
