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
	asnDB     *maxminddb.Reader
	countryDB *maxminddb.Reader
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

// Init 初始化日志文件和数据库
func Init() {
	var err error
	// 打开ASN数据库
	asnDB, err = maxminddb.Open("/data/ipinfo/db/asn.mmdb")
	if err != nil {
		log.Fatal("DB not exist Or error: ", err)
		log.Fatal("Error opening ASN database:", err)
	}

	// 打开国家数据库
	countryDB, err = maxminddb.Open("/data/ipinfo/db/country.mmdb")
	if err != nil {
		log.Fatal("DB not exist Or error: ", err)
		log.Fatal("Error opening country database:", err)
	}
}

// GetIPHandler 获取IP地址的处理函数
func GetIPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 尝试从X-Forwarded-For头部取得IP
	fwdIP := r.Header.Get("X-Forwarded-For")
	if fwdIP == "" {
		// 尝试从X-Real-IP头部取得IP
		fwdIP = r.Header.Get("X-Real-IP")
	}

	// 如果两个头部都没有，则从连接中获取IP
	if fwdIP == "" {
		ip, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
		if err != nil {
			http.Error(w, "Failed to resolve IP address", http.StatusInternalServerError)
			return
		}
		fwdIP = ip.IP.String()
	}

	// 验证 IP 地址格式
	if net.ParseIP(fwdIP) == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		log.Printf("Invalid IP address: %s", fwdIP)
		return
	}
	// 返回IP地址
	fmt.Fprintf(w, html.EscapeString(fwdIP))
}

// IPLookupHandler IP查询的处理函数
func IPLookupHandler(w http.ResponseWriter, r *http.Request) {
	// 允许跨站请求
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从请求中获取User-Agent头部，即浏览器信息
	userAgent := r.Header.Get("User-Agent")

	// 尝试从查询参数获取IP
	ipStr := r.URL.Query().Get("ip")

	// 尝试从X-Forwarded-For头部取得IP
	if ipStr == "" {
		ipStr = r.Header.Get("X-Forwarded-For")
	}

	// 尝试从X-Real-IP头部取得IP
	if ipStr == "" {
		ipStr = r.Header.Get("X-Real-IP")
	}

	// 如果两个头部都没有，则从连接中获取IP
	if ipStr == "" {
		var err error
		ipStr, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Failed to resolve IP address", http.StatusInternalServerError)
			return
		}
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		return
	}

	// 查询ASN记录
	var asn ASNRecord
	if err := asnDB.Lookup(ip, &asn); err != nil {
		http.Error(w, fmt.Sprintf("ASN Lookup failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 查询国家记录
	var country CountryRecord
	if err := countryDB.Lookup(ip, &country); err != nil {
		http.Error(w, fmt.Sprintf("Country Lookup failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 整理响应数据
	responseData := struct {
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
		IP:            ipStr,
		ASN:           asn.ASN,
		Domain:        asn.Domain,
		ISP:           asn.Name,
		ContinentCode: country.Continent,
		ContinentName: country.ContinentName,
		CountryCode:   country.Country,
		CountryName:   country.CountryName,
		UserAgent:     userAgent,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
