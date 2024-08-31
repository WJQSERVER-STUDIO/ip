package proxy

import (
	"io"
	"net/http"
	"net/url"
)

// BilibiliHandler 处理 /bilibili 路由的请求
func BilibiliHandler(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头部允许跨站请求
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//设置UA
	w.Header().Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	// 从请求中获取ip参数
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}
	ip := params.Get("ip")

	// 定义源API的URL
	apiURL := "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr"

	// 如果ip参数存在，将其添加到API URL
	if ip != "" {
		apiURL += "?ip=" + url.QueryEscape(ip)
	}

	// 向源API发送请求
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to get response from source API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 检查源API的响应状态码
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Source API returned non-OK status", http.StatusInternalServerError)
		return
	}

	// 读取源API返回的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// 设置响应头Content-Type为application/json
	w.Header().Set("Content-Type", "application/json")
	// 将源API的响应内容写入响应体
	w.Write(body)
}
