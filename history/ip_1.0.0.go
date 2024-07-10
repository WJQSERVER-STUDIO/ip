package main
import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "time"
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
	Continent    string `maxminddb:"continent"`
	Continent_name string `maxminddb:"continent_name"`
	Country   string `maxminddb:"country"`
	Country_name   string `maxminddb:"country_name"`
}
func main() {
    var err error
    // 打开ASN数据库
    asnDB, err = maxminddb.Open("/data/ipinfo/db/asn.mmdb")
    if err != nil {
        log.Fatal("Error opening ASN database:", err)
    }
    defer asnDB.Close()

    // 打开国家数据库
    countryDB, err = maxminddb.Open("/data/ipinfo/db/country.mmdb")
    if err != nil {
        log.Fatal("Error opening country database:", err)
    }
    defer countryDB.Close()

    // 打开或创建日志文件
    logFile, err := os.OpenFile("/data/ipinfo/log/access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }

    
          
            
    

          
          Expand Down
    
    
  
    defer logFile.Close()
    // 创建一个日志记录器
    logger := log.New(logFile, "", 0)
    // 设置HTTP路由处理器并启动服务器
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logRequest(logger, r)
        // 其他处理逻辑...
    })
    http.HandleFunc("/ip-lookup", ipLookupHandler)
    http.HandleFunc("/ip", getIPHandler)
    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
func logRequest(logger *log.Logger, r *http.Request) {
    // 假设 response size 和 status code 固定为 示例值
    // 在实际应用中，您可能需要动态获取这些值
    responseSize := 0
    statusCode := 0
    // 获取请求的IP地址
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    // 获取用户代理
    userAgent := r.UserAgent()
    // 获取日期信息
    dateTime := time.Now().Format("02/Jan/2006:15:04:05 -0700")
    // 格式化日志信息
    logEntry := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d \"-\" \"%s\"",
        ip, dateTime, r.Method, r.RequestURI, r.Proto, statusCode, responseSize, userAgent)
    // 将日志写入文件
    logger.Println(logEntry)
}
func getIPHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // 尝试从X-Forwarded-For头部取得IP
    fwdIP := r.Header.Get("X-Forwarded-For")
    if fwdIP == "" {
        fwdIP = r.Header.Get("X-Real-IP")
    }
    // 如果两个头部都没有，则从连接中获取IP
    if fwdIP == "" {
        ip, _, _ := net.SplitHostPort(r.RemoteAddr)
        fwdIP = ip
    }
    
    // 直接返回IP地址，不使用JSON格式化
    fmt.Fprintf(w, fwdIP)
}
func ipLookupHandler(w http.ResponseWriter, r *http.Request) {
    //允许跨站请求
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    // 从请求中获取User-Agent头部，即浏览器信息
    userAgent := r.Header.Get("User-Agent")   
    
    // 尝试从查询参数获取IP
    ipStr := r.URL.Query().Get("ip")
    if ipStr == "" {
        // 尝试从X-Forwarded-For头部取得IP
        fwdIP := r.Header.Get("X-Forwarded-For")
        if fwdIP == "" {
            fwdIP = r.Header.Get("X-Real-IP")
        }
        // 如果两个头部都没有，则从连接中获取IP
        if fwdIP == "" {
            ip, _, _ := net.SplitHostPort(r.RemoteAddr)
            fwdIP = ip
        }
        ipStr = fwdIP
    }
    ip := net.ParseIP(ipStr)
    if ip == nil {
        http.Error(w, "Invalid IP address", http.StatusBadRequest)
        return
    }
    
    // 查询ASN记录
    var asn ASNRecord
    err := asnDB.Lookup(ip, &asn)
    if err != nil {
        http.Error(w, fmt.Sprintf("ASN Lookup failed: %v", err), http.StatusInternalServerError)
        return
    }
    // 查询国家记录
    var country CountryRecord
    err = countryDB.Lookup(ip, &country)
    if err != nil {
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
        ContinentName: country.Continent_name,
        CountryCode:   country.Country,
        CountryName:   country.Country_name,
        UserAgent:     userAgent,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)
}
