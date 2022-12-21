package connection

import (
	"net"
	"net/http"
	"time"
)

var (
	Client *http.Client
)

func InitClient() {
	var (
		kaTimeout = 600 * time.Second
		timeout   = 10 * time.Second
	)
	defaultTransport := &http.Transport{
		Dial:                (&net.Dialer{KeepAlive: kaTimeout}).Dial,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	Client = &http.Client{Transport: defaultTransport, Timeout: timeout}
}
