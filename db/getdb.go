package db

import (
	"encoding/json"
	"fmt"
	"io"
	"ip/config"
	"net/http"
	"os"
	"time"
)

// jsonDBinfo 数据库信息结构体
type jsonDBinfo struct {
	RenewTime string
}

// 检测数据库是否存在
func isDBExists(asnDatabasePath string, countryDatabasePath string) bool {
	// os.Stat() 检测文件是否存在
	var (
		asnDBExists     bool
		countryDBExists bool
	)
	_, err := os.Stat(asnDatabasePath)
	if err == nil {
		asnDBExists = true
	} else if os.IsNotExist(err) {
		asnDBExists = false
	} else {
		logWarning("Failed to check ASN database existence: %s", err)
		asnDBExists = false
	}
	_, err = os.Stat(countryDatabasePath)
	if err == nil {
		countryDBExists = true
	} else if os.IsNotExist(err) {
		countryDBExists = false
	} else {
		logWarning("Failed to check IP database existence: %s", err)
		countryDBExists = false
	}

	if asnDBExists && countryDBExists {
		return true
	} else {
		return false
	}

}

// 检测DB信息是否存在
func isDBinfoExists(cfg *config.Config) bool {
	DBinfoPath := cfg.Mmdb.MmDBPath + "/DBinfo.json"
	_, err := os.Stat(DBinfoPath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		logWarning("Failed to check DBinfo existence: %s", err)
		return false
	}
}

// GetDBinfo 获取数据库信息
func GetDBinfo(cfg *config.Config) (renewTime string, err error) {
	// 检测 DBinfo.json 文件是否存在
	if !isDBinfoExists(cfg) {
		return "", fmt.Errorf("DBinfo.json not found")
	}

	// 读取 DBinfo.json 文件内容
	filePath := cfg.Mmdb.MmDBPath + "/DBinfo.json"
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 使用 JSON 解码器读取数据
	decoder := json.NewDecoder(file)
	var dbinfo jsonDBinfo
	if err := decoder.Decode(&dbinfo); err != nil {
		return "", fmt.Errorf("failed to decode JSON data: %v", err)
	}

	return dbinfo.RenewTime, nil
}

func Is2Update(cfg *config.Config) bool {
	if !isDBinfoExists(cfg) && !isDBExists(ASNDB_Path, CountryDB_Path) {
		return true
	}
	renewTime, err := GetDBinfo(cfg)

	if err != nil {
		logWarning("Failed to get DB info: %s", err)
	}
	// 若当前时间大于renewTime，则更新数据库 Format("2006-01-02 15:04:05")
	if time.Now().Format("2006-01-02 15:04:05") > renewTime {
		return true
	} else {
		return false
	}
}

func GetNewDB(cfg *config.Config) error {
	var err error
	if isDBinfoExists(cfg) {
		// 检测数据库是否存在
		if !isDBExists(ASNDB_Path, CountryDB_Path) {
			logInfo("MmDB not found, downloading...")
			err = pullNewDB(cfg) // 下载最新数据库
			if err != nil {
				logError("Failed to download new database: %v", err)
				return fmt.Errorf("failed to download new database: %v", err)
			}
			logInfo("Successfully downloaded new MmDB")
		}
		if isDBExists(ASNDB_Path, CountryDB_Path) {
			// 检测是否需要更新数据库
			if Is2Update(cfg) {
				logInfo("MmDB need to Update, downloading...")
				err = pullNewDB(cfg) // 下载最新数据库
				if err != nil {
					logError("Failed to download new database: %v", err)
					return fmt.Errorf("failed to download new database: %v", err)
				}
				logInfo("Successfully to Update MmDB")
			}
		}
	} else if !isDBinfoExists(cfg) {
		// 数据库不存在，下载最新数据库
		logInfo("MmDB not found, downloading...")
		err = pullNewDB(cfg) // 下载最新数据库
		if err != nil {
			logError("Failed to download database: %v", err)
			return fmt.Errorf("failed to download database: %v", err)
		}
		logInfo("Successfully downloaded MmDB")
	}
	return nil
}

func pullNewDB(cfg *config.Config) error {
	CloseDB() // 关闭当前数据库连接

	var err error

	// 下载 ASN 数据库文件
	err = DownloadASNDB(cfg.Mmdb.IPinfoKey, ASNDB_Path)
	if err != nil {
		logError("Failed to download ASN database: %v", err)
		return fmt.Errorf("failed to download ASN database: %v", err)
	}
	// 下载 IP 数据库文件
	err = DownloadCountryDB(cfg.Mmdb.IPinfoKey, CountryDB_Path)
	if err != nil {
		logError("Failed to download IP database: %v", err)
		return fmt.Errorf("failed to download IP database: %v", err)
	}

	// 记录数据库信息到 JSON 文件
	err = RecordDBinfo(cfg)
	if err != nil {
		logError("Failed to record DB info: %v", err)
		return fmt.Errorf("failed to record DB info: %v", err)
	}

	// 重载数据库
	err = ReloadDB()
	if err != nil {
		logError("Failed to reload database: %v", err)
		return fmt.Errorf("failed to reload database: %v", err)
	}
	return nil
}

// DownloadASNDB 下载 ASN 数据库文件
func DownloadASNDB(token string, outputPath string) error {
	// 构建下载 URL
	url := fmt.Sprintf("https://ipinfo.io/data/free/asn.mmdb?token=%s", token)

	// 发起 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		logError("Failed to download ASN database: %v", err)
		return fmt.Errorf("failed to download ASN database: %v", err)
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		logError("Failed to download ASN database: received status code %d", resp.StatusCode)
		return fmt.Errorf("failed to download ASN database: received status code %d", resp.StatusCode)
	}

	// 检查同名文件是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		// 如果文件存在，删除旧文件
		err = os.Remove(outputPath)
		if err != nil {
			logError("Failed to remove existing file: %v", err)
			return fmt.Errorf("failed to remove existing file: %v", err)
		}
	}

	// 创建输出文件并设置权限
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644) // 设置文件权限为 0644
	if err != nil {
		logError("Failed to create output file: %v", err)
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close() // 确保在函数结束时关闭文件

	// 将响应体内容写入文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		logError("Failed to write to output file: %v", err)
		return fmt.Errorf("failed to write to output file: %v", err)
	}

	return nil
}

// DownloadCountryDB 下载国家数据库文件
func DownloadCountryDB(token string, outputPath string) error {
	// 构建下载 URL
	url := fmt.Sprintf("https://ipinfo.io/data/free/country.mmdb?token=%s", token)

	// 发起 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		logError("Failed to download IP database: %v", err)
		return fmt.Errorf("failed to download Country database: %v", err)
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		logError("Failed to download IP database: received status code %d", resp.StatusCode)
		return fmt.Errorf("failed to download Country database: received status code %d", resp.StatusCode)
	}

	// 检查同名文件是否已存在
	if _, err := os.Stat(outputPath); err == nil {
		// 如果文件存在，删除旧文件
		err = os.Remove(outputPath)
		if err != nil {
			logError("Failed to remove existing file: %v", err)
			return fmt.Errorf("failed to remove existing file: %v", err)
		}
	}

	// 创建输出文件并设置权限
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644) // 设置文件权限为 0644
	if err != nil {
		logError("Failed to create output file: %v", err)
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close() // 确保在函数结束时关闭文件

	// 将响应体内容写入文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		logError("Failed to write to output file: %v", err)
		return fmt.Errorf("failed to write to output file: %v", err)
	}

	return nil
}

// RecordDBinfo 记录数据库信息到 JSON 文件
func RecordDBinfo(cfg *config.Config) error {
	// 计算更新的时间
	renewTime := time.Now().Add(time.Duration(cfg.Mmdb.UpdateFreq) * time.Hour).Format("2006-01-02 15:04:05")

	// 构建 DBinfo.json 文件内容
	dbinfo := jsonDBinfo{
		RenewTime: renewTime,
	}

	// 写入 DBinfo.json 文件
	filePath := cfg.Mmdb.MmDBPath + "/DBinfo.json"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create or open file: %v", err)
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 使用 JSON 编码器写入数据
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(dbinfo); err != nil {
		return fmt.Errorf("failed to encode JSON data: %v", err)
	}

	return nil
}
