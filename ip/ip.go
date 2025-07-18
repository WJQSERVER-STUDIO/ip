package ip

import (
	"ip/db"
	"net"
	"net/netip"

	"github.com/infinite-iroha/touka"
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

func IPHandler(c *touka.Context) {

	var (
		ip string
		ua string
	)

	// 检测url参数中是否有ip参数 如果IP参数存在，且合法，则使用IP参数
	ipStr := c.Query("ip")
	if ipStr != "" && net.ParseIP(ipStr) != nil {
		ip = ipStr
	} else {
		ip = c.ClientIP()
	}

	// 预处理IP地址，转为net.IP类型
	netIP, err := netip.ParseAddr(ip)
	if err != nil {
		c.Warnf("Invalid IP address: %s", ip)
		c.JSON(400, touka.H{"error": "Invalid IP address"})
		return
	}

	// call DB 获取 GeoIP 数据
	results, err := db.SearchDB(netIP)
	if err != nil {
		c.Errorf("SearchDB error: %s", err)
		c.JSON(500, touka.H{"error": "Internal Server Error"})
		return
	}

	// 获取User-Agent
	ua = c.UserAgent()

	c.Header("Access-Control-Allow-Origin", "*") // 允许跨域请求
	c.Header("Content-Type", "application/json") // 内容类型

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

func IPPureHandler(c *touka.Context) {
	ip := c.ClientIP()

	// 预处理IP地址，转为net.IP类型
	netIP := net.ParseIP(ip)
	if netIP == nil {
		c.Warnf("Invalid IP address: %s", ip)
		c.JSON(400, touka.H{"error": "Invalid IP address"})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*") // 允许跨域请求

	// 纯IP响应,非json格式
	c.Text(200, ip)
}
