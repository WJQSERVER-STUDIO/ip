package db

import (
	"fmt"
	"ip/config"
	"ip/logger"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	logWarning = logger.LogWarning
	logError   = logger.LogError
)

var (
	asnDB          *maxminddb.Reader
	countryDB      *maxminddb.Reader
	ASNDB_Path     = "/data/ip/db/asn.mmdb"     // ASN 数据库路径
	CountryDB_Path = "/data/ip/db/country.mmdb" // 国家数据库路径
)

// ASNDB 结构体定义
type ASNRecord struct {
	ASN    string `maxminddb:"asn"`    // 自治系统编号
	Domain string `maxminddb:"domain"` // 域名
	Name   string `maxminddb:"name"`   // 名称
}

// CountryDB 结构体定义
type CountryRecord struct {
	Continent     string `maxminddb:"continent"`      // 大洲代码
	ContinentName string `maxminddb:"continent_name"` // 大洲名称
	Country       string `maxminddb:"country"`        // 国家代码
	CountryName   string `maxminddb:"country_name"`   // 国家名称
}

// openDB 打开数据库
func openDB(db **maxminddb.Reader, path string) {
	var err error
	*db, err = maxminddb.Open(path)
	if err != nil {
		logError("Error opening database at %s: %v", path, err)
		return
	}
}

// DBinit 初始化日志文件和数据库
func DBinit(cfg *config.Config) {
	ASNDB_Path = cfg.Mmdb.ASNDBPath
	CountryDB_Path = cfg.Mmdb.CountryDBPath
	openDB(&asnDB, ASNDB_Path)         // 打开 ASN 数据库
	openDB(&countryDB, CountryDB_Path) // 打开国家数据库
}

// ReloadDB 重新加载数据库
func ReloadDB() error {
	// 关闭当前数据库连接
	if asnDB != nil {
		asnDB.Close()
	}
	if countryDB != nil {
		countryDB.Close()
	}

	// 重新打开数据库
	openDB(&asnDB, ASNDB_Path)
	openDB(&countryDB, CountryDB_Path)

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if asnDB != nil {
		asnDB.Close() // 关闭 ASN 数据库
		logInfo("ASN database closed successfully.")
	}
	if countryDB != nil {
		countryDB.Close() // 关闭国家数据库
		logInfo("Country database closed successfully.")
	}
}

// SearchDB 根据 IP 地址查询 ASN 和国家信息
func SearchDB(ip net.IP) (results []string, err error) {
	var (
		asn     ASNRecord
		country CountryRecord
	)

	// 查询 ASN 信息
	err = asnDB.Lookup(ip, &asn)
	if err != nil {
		logError("ASN Lookup failed: %v", err)
		return nil, fmt.Errorf("ASN Lookup failed: %v", err)
	}

	// 查询国家信息
	err = countryDB.Lookup(ip, &country)
	if err != nil {
		logError("Country Lookup failed: %v", err)
		return nil, fmt.Errorf("Country Lookup failed: %v", err)
	}

	// 聚合结果到切片中
	results = []string{
		asn.ASN,
		asn.Domain,
		asn.Name,
		country.Continent,
		country.ContinentName,
		country.Country,
		country.CountryName,
	}

	return results, nil
}
