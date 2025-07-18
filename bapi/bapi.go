package bapi

import (
	"net/http"
	"net/netip"
	"net/url"

	"github.com/infinite-iroha/touka"
)

// 使用req库处理 /bilibili 路由的请求并使用chrome的TLS指纹
func Bilibili(c *touka.Context) {

	// 设置响应头
	c.Header("Access-Control-Allow-Origin", "*") // 允许跨域请求
	c.Header("Content-Type", "application/json") // 内容类型

	// 从请求中获取ip参数
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, touka.H{
			"error": "Missing ip parameter",
		})
		// IP METHOD URL UA PROTOCAL
		c.Warnf("%s %s %s %s %s Missing ip parameter", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto)
		return
	}

	_, err := netip.ParseAddr(ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, touka.H{
			"error": "Invalid ip parameter",
		})
		// IP METHOD URL UA PROTOCAL
		c.Warnf("%s %s %s %s %s Invalid ip parameter", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto)
		return
	}

	// 定义源API的URL并添加ip参数
	apiURL := "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr?ip=" + url.QueryEscape(ip)

	// 使用req库发送请求并使用chrome的TLS指纹
	client := c.GetHTTPC()
	rb := client.NewRequestBuilder("GET", apiURL)
	rb.NoDefaultHeaders()
	rb.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")

	body, err := rb.Bytes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, touka.H{
			"error": "Failed to get response from source API",
		})
		// IP METHOD URL UA PROTOCAL ERROR
		c.Errorf("%s %s %s %s %s Failed to get response from source API: %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto, err)
		return
	}

	if len(body) == 0 {
		c.JSON(http.StatusInternalServerError, touka.H{
			"error": "Failed to get response from source API",
		})
		// IP METHOD URL UA PROTOCAL ERROR
		c.Errorf("%s %s %s %s %s Failed to get response from source API: %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto, err)
		return
	}

	// 返回上游API的响应内容
	c.Raw(http.StatusOK, "application/json", body)
	// IP METHOD URL UA PROTOCAL
	c.Infof("%s %s %s %s %s Success", c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("User-Agent"), c.Request.Proto)
}
