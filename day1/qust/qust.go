package qust

import (
	"io"
	"net/http"
	qjson "qust/convert/json"
	"strings"
)

// Engine qust主函数
type Engine struct {
	BaseUrl BaseUrl
}

// Res 返回值类型
type Res struct {
	StatusCode int
	Proto      string
	Header     http.Header
	Body       []byte
}

// Json 获取Json数据
func (res *Res) Json() (*qjson.Data, error) {
	return qjson.New(res.Body)
}

// Xml 获取Xml数据
func (res *Res) Xml() {

}

// Html 获取Html数据
func (res *Res) Html() {

}

// Text 获取文本
func (res *Res) Text() {

}

var (
	defaultBaseUrl = BaseUrl("")
)

type BaseUrl string

// New 创建Qust
//
//	args:
//	BaseUrl:基础路径
func New(args ...any) *Engine {
	baseUrl := defaultBaseUrl
	for _, arg := range args {
		switch arg.(type) {
		case BaseUrl:
			baseUrl = arg.(BaseUrl)
		}
	}
	return &Engine{BaseUrl: baseUrl}
}

func (engine *Engine) Get(url string) (*Res, error) {
	var builder strings.Builder
	builder.WriteString(string(engine.BaseUrl))
	if len(url) > 8 && url[:7] != "http://" && url[:8] != "https://" {
		builder.WriteString(url)
	}

	html, err := http.Get(builder.String())
	if err != nil {
		return nil, err
	}
	defer html.Body.Close()
	body, err := io.ReadAll(html.Body)
	if err != nil {
		return nil, err
	}
	res := &Res{StatusCode: html.StatusCode, Proto: html.Proto, Header: html.Header, Body: body}
	return res, nil
}

func (engine *Engine) Post(url string) {

}
func (engine *Engine) Delete(url string) {

}
func (engine *Engine) Put(url string) {

}
