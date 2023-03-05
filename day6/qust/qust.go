package qust

import (
	"net/http"
	urlpkg "net/url"
)

var (
	defaultBaseUrl BaseUrl      = ""
	defaultClient  *http.Client = nil
	defaultVersion              = 2
)

type (
	BaseUrl string
)

// Engine qust主函数
type Engine struct {
	BaseUrl BaseUrl
	Client  *http.Client
	Version int
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

func (engine *Engine) Ask(method string, url string) *Req {
	//url拼接
	if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
		url = string(engine.BaseUrl) + url
	}
	u, _ := urlpkg.Parse(url)
	//生成一个请求体
	req := &Req{u, method, map[string]any{}, engine.Client, nil, make(http.Header)}
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
