package qust

import (
	"net/http"
)

var (
	defaultBaseUrl BaseUrl      = ""
	defaultClient  *http.Client = nil
)

type (
	BaseUrl string
)

// Engine qust主函数
type Engine struct {
	BaseUrl BaseUrl
	Client  *http.Client
}

// New 创建Qust
//
//	args:
//	BaseUrl:基础路径,
//	http.Client:全局客户端,
func New(args ...any) *Engine {
	baseUrl := defaultBaseUrl
	client := defaultClient
	for _, arg := range args {
		switch arg.(type) {
		case BaseUrl:
			baseUrl = arg.(BaseUrl)
		case *http.Client:
			client = arg.(*http.Client)
		}
	}
	return &Engine{BaseUrl: baseUrl, Client: client}
}

func (engine *Engine) Ask(method string, url string) *Req {
	//url拼接
	if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
		url = string(engine.BaseUrl) + url
	}
	//生成一个请求体
	req := &Req{Url: url, Method: method, Client: engine.Client, Query: map[string]any{}}
	//支持其他协议
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	req.Request = request
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
