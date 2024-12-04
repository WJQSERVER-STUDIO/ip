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
	ASNDB_Path     = "/data/ip/db/asn.mmdb"
	CountryDB_Path = "/data/ip/db/country.mmdb"
)

// ASNDB
type ASNRecord struct {
	ASN    string `maxminddb:"asn"`
	Domain string `maxminddb:"domain"`
	Name   string `maxminddb:"name"`
}

// CountryDB
type CountryRecord struct {
	Continent     string `maxminddb:"continent"`
	ContinentName string `maxminddb:"continent_name"`
	Country       string `maxminddb:"country"`
	CountryName   string `maxminddb:"country_name"`
}

func openDB(db **maxminddb.Reader, path string) {
	var err error
	*db, err = maxminddb.Open(path)
	if err != nil {
		logError("Error opening database at %s: %v", path, err)
		return
	}
}

// 初始化日志文件和数据库
func DBinit(cfg *config.Config) {
	openDB(&asnDB, ASNDB_Path)
	openDB(&countryDB, CountryDB_Path)
}

func SearchDB(ip net.IP) (results []string, err error) {

	var (
		asn     ASNRecord
		country CountryRecord
	)

	err = asnDB.Lookup(ip, &asn)
	if err != nil {
		logError("ASN Lookup failed: %v", err)
		return nil, fmt.Errorf("ASN Lookup failed: %v", err)
	}

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
