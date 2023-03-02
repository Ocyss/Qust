package qust

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Req struct {
	Url     string
	Method  string
	Headers map[string]string
	Query   map[string]any
	Request *http.Request
	Client  *http.Client
	Data    any
}

func (req *Req) Do() (*Res, error) {
	//判断有没有指定Client，没有就使用默认
	client := http.DefaultClient
	if req.Client != nil {
		client = req.Client
	}
	//进行param拼接
	params := url.Values{}
	for k, v := range req.Query {
		params.Add(k, fmt.Sprintf("%v", v))
	}
	if req.Request.URL.RawQuery != "" {
		req.Request.URL.RawQuery += "&"
	}
	req.Request.URL.RawQuery += params.Encode()
	//发送请求
	response, err := client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", response.Request.URL, response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	res := &Res{StatusCode: response.StatusCode, Proto: response.Proto, Header: response.Header, Body: body}
	return res, nil
}

func (req *Req) dataEncoding(data interface{}) (io.Reader, error) {
	var err error
	jsonData := make([]byte, 0, 0)
	if t, ok := data.(string); ok {
		jsonData = []byte(t)
	} else {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	return bytes.NewBuffer(jsonData), nil
}

func (req *Req) AddHeader(k any, args ...any) *Req {
	//判断是不是标准的(str,str)是就直接添加放回
	if s, ok := k.(string); ok && len(args) == 1 {
		req.Request.Header.Add(s, args[0].(string))
		return req
	}
	//针对以下几种情况
	//AddHeader("head2", "value2", "head3", "value3")
	//AddHeader(map[string]string{"head4": "value4"})
	s1 := ""
	args = append([]any{k}, args...)
	for _, arg := range args {
		switch arg.(type) {
		case string:
			if s1 == "" {
				s1 = arg.(string)
			} else {
				req.Request.Header.Add(s1, arg.(string))
				s1 = ""
			}
		case map[string]string:
			v := k.(map[string]string)
			for k, v := range v {
				req.Request.Header.Add(k, v)
			}
		}
	}

	return req
}

func (req *Req) SetQuery(k any, args ...any) *Req {
	if s, ok := k.(string); ok && len(args) == 1 {
		req.Query[s] = args[0]
		return req
	}
	s1 := ""
	args = append([]any{k}, args...)
	for _, arg := range args {
		switch arg.(type) {
		case string:
			if s1 == "" {
				s1 = arg.(string)
			} else {
				req.Query[s1] = arg.(string)
				s1 = ""
			}
		case map[string]any:
			for k, v := range arg.(map[string]any) {
				req.Query[k] = v
			}
		}
	}
	return req
}
