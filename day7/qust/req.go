package qust

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Req struct {
	Url    *url.URL
	Method string
	Query  map[string]any
	Client *http.Client
	Body   io.Reader
	Header http.Header
	Proxy  string
}
type DataType string

const (
	JSON          DataType = "application/json"
	FROM          DataType = "multipart/form-data"
	XML           DataType = "text/xml"
	MultipartForm DataType = "application/x-www-form-urlencoded"
)

func dataEncoding(args []any) (res map[string]string) {
	res = make(map[string]string)
	s1 := ""
	for _, arg := range args {
		if s2, ok := arg.(string); ok {
			if s1 == "" {
				s1 = s2
			} else {
				res[s1] = s2
				s1 = ""
			}
		} else if r := reflect.ValueOf(arg); r.Kind() == reflect.Map { //判断是不是map类型，然后进行添加
			iter := r.MapRange()
			for iter.Next() {
				res[fmt.Sprintf("%v", iter.Key())] = fmt.Sprintf("%v", iter.Value())
			}
		}
	}
	return
}

func (req *Req) Do() (*Res, error) {
	client := req.Client
	//进行param拼接
	params := url.Values{}
	for k, v := range req.Query {
		params.Add(k, fmt.Sprintf("%v", v))
	}
	if req.Url.RawQuery != "" {
		req.Url.RawQuery += "&"
	}
	req.Url.RawQuery += params.Encode()
	//发送请求
	Req, err := http.NewRequest(req.Method, req.Url.String(), req.Body)
	Req.Header = req.Header
	if req.Proxy != "" {
		urlproxy, err := url.Parse(req.Proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(urlproxy)}
	}
	response, err := client.Do(Req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	res := &Res{StatusCode: response.StatusCode, Proto: response.Proto, Header: response.Header, Body: body}
	if response.StatusCode > 399 {
		return res, fmt.Errorf("http get error : uri=%v , statusCode=%v", response.Request.URL, response.StatusCode)
	}
	return res, nil
}

func (req *Req) AddHeader(k any, args ...any) *Req {
	//判断是不是标准的(str,str)是就直接添加放回
	if s, ok := k.(string); ok && len(args) == 1 {
		req.Header.Add(s, args[0].(string))
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
				req.Header.Add(s1, arg.(string))
				s1 = ""
			}
		case map[string]string:
			v := k.(map[string]string)
			for k, v := range v {
				req.Header.Add(k, v)
			}
		}
	}
	return req
}

func (req *Req) SetQuery(k any, args ...any) *Req {
	//标准 2 参数，直接添加放回
	if s, ok := k.(string); ok && len(args) == 1 {
		req.Query[s] = args[0]
		return req
	}
	//多参数，或者结构体，进行遍历
	args = append([]any{k}, args...)
	for k, v := range dataEncoding(args) {
		req.Query[k] = v
	}
	return req
}

func (req *Req) SetData(dtype DataType, datas ...any) (err error) {
	var reader io.Reader
	if len(datas) == 0 {
		return errors.New("not enough parameters")
	} else if req.Method == "GET" {
		return errors.New("get requests cannot have a body")
	}
	switch dtype {
	case JSON:
		if v, ok := datas[0].([]byte); ok {
			reader = bytes.NewBuffer(v)
		} else {
			jsonData := make([]byte, 0, 0)
			if t, ok := datas[0].(string); ok {
				jsonData = []byte(t)
			} else {
				jsonData, err = json.Marshal(datas[0])
				if err != nil {
					return
				}
			}
			reader = bytes.NewBuffer(jsonData)
		}
	case FROM, MultipartForm:
		if v, ok := datas[0].(url.Values); ok {
			reader = strings.NewReader(v.Encode())
		} else {
			fromurl := make(url.Values)
			for k, v := range dataEncoding(datas) {
				fromurl.Add(k, v)
			}
			reader = strings.NewReader(fromurl.Encode())
		}
	case XML:
		xmlData := make([]byte, 0, 0)
		if t, ok := datas[0].(string); ok {
			xmlData = []byte(t)
		} else {
			xmlData, err = xml.Marshal(datas[0])
			if err != nil {
				return
			}
		}
		reader = bytes.NewBuffer(xmlData)
	default:
		if v, ok := datas[0].(io.Reader); ok {
			reader = v
		} else {
			return errors.New("use a custom contentType, the second parameter must be of type io.Reader")
		}
	}
	rc, ok := reader.(io.ReadCloser)
	if !ok && reader != nil {
		rc = io.NopCloser(reader)
	}
	req.Body = rc
	req.Header.Set("Content-Type", string(dtype))
	return nil
}

func (req *Req) SetCookie(data any) *Req {
	switch data.(type) {
	case string:
		req.AddHeader("cookie", data.(string))
	}
	return req
}

func (req *Req) SetUA(ua string) *Req {
	req.Header.Set("User-Agent", ua)
	return req
}
