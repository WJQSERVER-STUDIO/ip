package db

import (
	"ip/config"
	"time"
)

func LoopForUpdate(cfg *config.Config) {
	var err error
	for {
		if Is2Update(cfg) {
			logInfo("Updating database...")
			err = GetNewDB(cfg)
			if err != nil {
				logWarning("Failed to update database: %s", err)
			}
		}

		// 睡眠6小时
		time.Sleep(6 * time.Hour)
	}
}
