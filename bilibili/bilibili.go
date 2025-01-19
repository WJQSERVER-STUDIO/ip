package bilibili

import (
	"io"
	"ip/logger"
	"net/http"
	"net/netip"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
)

var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	logWarning = logger.LogWarning
	logError   = logger.LogError
)

// 使用req库处理 /bilibili 路由的请求并使用chrome的TLS指纹
func Bilibili(c *gin.Context) {

	// 设置响应头
	c.Header("Access-Control-Allow-Origin", "*") // 允许跨域请求
	c.Header("Content-Type", "application/json") // 内容类型

	// 从请求中获取ip参数
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing ip parameter",
		})
		// IP METHOD URL UA PROTOCAL
		logWarning("%s %s %s %s %s Missing ip parameter", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto)
		return
	}

	// 验证IP是否正确
	/*
		if net.ParseIP(ip) == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ip parameter",
			})
			// IP METHOD URL UA PROTOCAL
			logWarning("%s %s %s %s %s Invalid ip parameter", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto)
			return
		}
	*/
	_, err := netip.ParseAddr(ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ip parameter",
		})
		// IP METHOD URL UA PROTOCAL
		logWarning("%s %s %s %s %s Invalid ip parameter", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto)
		return
	}

	// 定义源API的URL并添加ip参数
	apiURL := "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr?ip=" + url.QueryEscape(ip)

	// 使用req库发送请求并使用chrome的TLS指纹
	client := req.C().
		//SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36").
		SetTLSFingerprintChrome().
		ImpersonateChrome()

	resp, err := client.R().Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get response from source API",
		})
		// IP METHOD URL UA PROTOCAL ERROR
		logError("%s %s %s %s %s Failed to get response from source API: %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto, err)
		return
	}
	defer resp.Body.Close()

	// 检查源API的响应状态码
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Source API returned non-OK status",
		})
		// IP METHOD URL UA PROTOCAL ERROR
		logError("%s %s %s %s %s Source API returned non-OK status: %d", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto, resp.StatusCode)
		return
	}

	// 读取源API返回的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read response body",
		})
		// IP METHOD URL UA PROTOCAL ERROR
		logError("%s %s %s %s %s Failed to read response body: %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto, err)
		return
	}

	// 返回上游API的响应内容
	c.Data(http.StatusOK, "application/json", body)
	// IP METHOD URL UA PROTOCAL
	logInfo("%s %s %s %s %s Success", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto)
}
