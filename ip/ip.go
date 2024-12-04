package ip

import (
	"ip/db"
	"ip/logger"
	"net"

	"github.com/gin-gonic/gin"
)

var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	logWarning = logger.LogWarning
	logError   = logger.LogError
)

// 构造响应json
type Response struct {
	IP            string `json:"ip"`
	ASN           string `json:"asn"`
	Domain        string `json:"domain"`
	ISP           string `json:"isp"`
	ContinentCode string `json:"continent_code"`
	ContinentName string `json:"continent_name"`
	CountryCode   string `json:"country_code"`
	CountryName   string `json:"country_name"`
	UserAgent     string `json:"user_agent"`
}

func IPHandler(c *gin.Context) {

	var (
		ip string
		ua string
	)

	// 检测url参数中是否有ip参数 如果IP参数存在，且合法，则使用IP参数
	ipStr := c.Query("ip")
	if ipStr != "" && net.ParseIP(ipStr) != nil {
		ip = ipStr
	} else {
		// 优先获取X-Forwarded-For，其次是X-Real-IP,最后是c.clientIP()
		fwdIP := c.GetHeader("X-Forwarded-For")
		realIP := c.GetHeader("X-Real-IP")
		if fwdIP != "" {
			ip = fwdIP
		} else if realIP != "" {
			ip = realIP
		} else {
			ip = c.ClientIP()
		}
	}

	// 预处理IP地址，转为net.IP类型
	netIP := net.ParseIP(ip)
	if netIP == nil {
		logWarning("Invalid IP address: ", ip)
		c.JSON(400, gin.H{"error": "Invalid IP address"})
		return
	}

	// call DB 获取 GeoIP 数据
	results, err := db.SearchDB(netIP)
	if err != nil {
		logError("SearchDB error: ", err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// 获取User-Agent
	ua = c.GetHeader("User-Agent")

	// 构造响应
	response := Response{
		IP:            ip,
		ASN:           results[0],
		Domain:        results[1],
		ISP:           results[2],
		ContinentCode: results[3],
		ContinentName: results[4],
		CountryCode:   results[5],
		CountryName:   results[6],
		UserAgent:     ua,
	}

	// 响应json
	c.JSON(200, response)
}

func IPPureHandler(c *gin.Context) {
	var ip string
	// 优先获取X-Forwarded-For，其次是X-Real-IP,最后是c.clientIP()
	fwdIP := c.GetHeader("X-Forwarded-For")
	realIP := c.GetHeader("X-Real-IP")
	if fwdIP != "" {
		ip = fwdIP
	} else if realIP != "" {
		ip = realIP
	} else {
		ip = c.ClientIP()
	}

	// 预处理IP地址，转为net.IP类型
	netIP := net.ParseIP(ip)
	if netIP == nil {
		logWarning("Invalid IP address: ", ip)
		c.JSON(400, gin.H{"error": "Invalid IP address"})
		return
	}
	// 纯IP响应,非json格式
	c.String(200, ip)
}
