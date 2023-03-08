package qust

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"math/rand"
	"net"
	"net/http"
	urlpkg "net/url"
)

var (
	defaultBaseUrl BaseUrl = ""
	defaultClient          = &http.Client{}
	defaultVersion         = 2
)

type (
	BaseUrl string
)

// Engine qust主函数
type Engine struct {
	BaseUrl BaseUrl
	Client  *http.Client
	Version int
	Proxys  *Proxys
}

// New 创建Qust
//
//	args:
//	BaseUrl:基础路径,
//	http.Client:全局客户端,
func New(args ...any) *Engine {
	baseUrl := defaultBaseUrl
	client := defaultClient
	version := defaultVersion
	for _, arg := range args {
		switch arg.(type) {
		case BaseUrl:
			baseUrl = arg.(BaseUrl)
		case *http.Client:
			client = arg.(*http.Client)
		}
	}
	return &Engine{BaseUrl: baseUrl, Client: client, Version: version}
}

func (engine *Engine) SetProxy(protocol string, address string) error {
	httpTransport := &http.Transport{}
	switch protocol {
	case "SOCKS5":
		dialer, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
		if err != nil {
			return err
		}
		httpTransport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial("tcp", address)
		}
	case "HTTPS", "HTTP":
		urlproxy, err := urlpkg.Parse(fmt.Sprintf("%s://%s", protocol, address))
		if err != nil {
			return err
		}
		httpTransport.Proxy = http.ProxyURL(urlproxy)
	default:
		return errors.New("only SOCKS5 is supported HTTPS,HTTP,protocol")
	}
	engine.Client.Transport = httpTransport
	return nil
}

func (engine *Engine) SetProxys(proxys []string, mode int) *Proxys {
	p := &Proxys{proxys, 0, len(proxys), mode}
	engine.Proxys = p
	return p
}

func (engine *Engine) Ask(method string, url string) *Req {
	//url拼接
	if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
		url = string(engine.BaseUrl) + url
	}
	u, _ := urlpkg.Parse(url)
	p := ""
	if engine.Proxys != nil {
		switch engine.Proxys.Mode {
		case 0:
			p = engine.Proxys.Proxy[engine.Proxys.Index]
			engine.Proxys.Index++
			if engine.Proxys.Index == engine.Proxys.Size {
				engine.Proxys.Index = 0
			}
		case 1:
			p = engine.Proxys.Proxy[rand.Intn(engine.Proxys.Size)]
		}
	}
	//生成一个请求体
	req := &Req{u, method, map[string]any{}, engine.Client, nil, make(http.Header), p}
	return req
}

func (engine *Engine) Get(url string) *Req {
	return engine.Ask("GET", url)
}
func (engine *Engine) Post(url string) *Req {
	return engine.Ask("POST", url)
}
func (engine *Engine) Delete(url string) *Req {
	return engine.Ask("DELETE", url)

}
func (engine *Engine) Put(url string) *Req {
	return engine.Ask("PUT", url)
}
func (engine *Engine) Patch(url string) *Req {
	return engine.Ask("PATCH ", url)
}
