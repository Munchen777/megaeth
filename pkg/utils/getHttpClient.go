package utils

import (
	"net/url"
	"time"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"

	"main/pkg/global"
)

var clientIndex int

func CreateClient(proxy string) *fasthttp.Client {
	var dial fasthttp.DialFunc

	if proxy != "" {
		proxy, err := url.Parse(proxy)
		if err != nil {
			log.Panicf("Error Unparsing Proxy: %v\n", err)
		}

		switch proxy.Scheme {
		case "http", "https":
			dial = fasthttpproxy.FasthttpHTTPDialer(proxy.String())
		case "socks4", "socks5":
			dial = fasthttpproxy.FasthttpSocksDialer(proxy.String())
		default:
			log.Panicf("Unsupported proxy scheme: %s\n", proxy.Scheme)
		}
	}
	client := &fasthttp.Client{
		Dial: 						   dial,
		MaxIdleConnDuration:           90 * time.Second,
		DisablePathNormalizing:        true,
		ReadTimeout:                   10 * time.Second,
		WriteTimeout:                  10 * time.Second,
		MaxConnWaitTimeout:            15 * time.Second,
		StreamResponseBody:            true,
	}
	return client
}

func GetClient() *fasthttp.Client {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if len(global.Clients) == 0 {
		return nil
	}

	client := global.Clients[clientIndex]
	clientIndex = (clientIndex + 1) % len(global.Clients)

	return client
}
