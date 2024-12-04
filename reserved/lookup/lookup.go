package lookup

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"

	"github.com/oschwald/maxminddb-golang"
)

var (
	asnDB          *maxminddb.Reader
	countryDB      *maxminddb.Reader
	ASNDB_Path     = "/data/ip/db/asn.mmdb"
	CountryDB_Path = "/data/ip/db/country.mmdb"
)

// ASNRecord 保存ASN数据库的查询结果
type ASNRecord struct {
	ASN    string `maxminddb:"asn"`
	Domain string `maxminddb:"domain"`
	Name   string `maxminddb:"name"`
}

// CountryRecord 保存国家数据库的查询结果
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
		log.Fatalf("Error opening database at %s: %v", path, err)
	}
}

// Init 初始化日志文件和数据库
func Init() {
	openDB(&asnDB, ASNDB_Path)
	openDB(&countryDB, CountryDB_Path)
}

// GetIPHandler 获取IP地址的处理函数
func GetIPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fwdIP := getForwardedIP(r)

	if fwdIP == "" || net.ParseIP(fwdIP) == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		log.Printf("Invalid IP address: %s", fwdIP)
		return
	}

	fmt.Fprintf(w, html.EscapeString(fwdIP))
}

// getForwardedIP 从请求中获取转发的IP地址
func getForwardedIP(r *http.Request) string {
	if fwdIP := r.Header.Get("X-Forwarded-For"); fwdIP != "" {
		return fwdIP
	}
	return r.Header.Get("X-Real-IP")
}

// getIPFromRequest 从请求中获取IP地址
func getIPFromRequest(r *http.Request) string {
	ipStr := r.URL.Query().Get("ip")
	if ipStr == "" {
		ipStr = getForwardedIP(r)
	}
	if ipStr == "" {
		ipStr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ipStr
}

// IPLookupHandler IP查询的处理函数
func IPLookupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userAgent := r.Header.Get("User-Agent")
	ipStr := getIPFromRequest(r)

	ip := net.ParseIP(ipStr)
	if ip == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		return
	}

	responseData, err := getIPInfo(ip, userAgent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func getIPInfo(ip net.IP, userAgent string) (interface{}, error) {
	var asn ASNRecord
	if err := asnDB.Lookup(ip, &asn); err != nil {
		return nil, fmt.Errorf("ASN Lookup failed: %v", err)
	}

	var country CountryRecord
	if err := countryDB.Lookup(ip, &country); err != nil {
		return nil, fmt.Errorf("Country Lookup failed: %v", err)
	}

	return struct {
		IP            string `json:"ip"`
		ASN           string `json:"asn"`
		Domain        string `json:"domain"`
		ISP           string `json:"isp"`
		ContinentCode string `json:"continent_code"`
		ContinentName string `json:"continent_name"`
		CountryCode   string `json:"country_code"`
		CountryName   string `json:"country_name"`
		UserAgent     string `json:"user_agent"`
	}{
		IP:            ip.String(),
		ASN:           asn.ASN,
		Domain:        asn.Domain,
		ISP:           asn.Name,
		ContinentCode: country.Continent,
		ContinentName: country.ContinentName,
		CountryCode:   country.Country,
		CountryName:   country.CountryName,
		UserAgent:     userAgent,
	}, nil
}
