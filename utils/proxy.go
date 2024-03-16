package utils

import (
	"math/rand"

	"github.com/DataHenHQ/useragent"
)

func GetUserAgent() string {
	ua, _ := useragent.Desktop()
	if ua == "" {
		ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"
	}
	ua += " baiduboxapp"
	return ua
}

func GetUserAgentByMobile() string {
	ua, _ := useragent.Mobile()
	if ua == "" {
		ua = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/112.0.0.0"
	}
	ua += " baiduboxapp"
	return ua
}

func GetXT5Auth(n int) string {
	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func GetProxyUrl() string {
	return "cloudnproxy.baidu.com:443"
}
