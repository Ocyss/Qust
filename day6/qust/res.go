package qust

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	qfile "qust/convert/file"
	qhtml "qust/convert/html"
	qjson "qust/convert/json"
	"strings"
)

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
func (res *Res) Html() (*qhtml.Data, error) {
	doc, err := html.Parse(bytes.NewReader(res.Body))
	if err != nil {
		return nil, err
	}
	return qhtml.New(doc), nil
}

// File 获取File数据
func (res *Res) File() *qfile.Data {
	return qfile.New(res.Body)
}

// Encoding 进行编码
func (res *Res) Encoding(args ...any) error {
	var e encoding.Encoding    //编码格式
	for _, arg := range args { //参数判断，可以传字符串，或者直接编码格式
		switch arg.(type) {
		case string:
			//直接支持常用编码，方便使用
			if v := arg.(string); strings.EqualFold(v, "gbk") {
				e = simplifiedchinese.GBK
			} else if strings.EqualFold(v, "GB18030") {
				e = simplifiedchinese.GB18030
			} else if strings.EqualFold(v, "HZGB2312") {
				e = simplifiedchinese.HZGB2312
			}
		case encoding.Encoding:
			e = arg.(encoding.Encoding)
		}
	}
	if e.NewEncoder == nil {
		return errors.New("you did not pass the encoding format")
	} else {
		reader := transform.NewReader(bytes.NewReader(res.Body), e.NewDecoder())
		d, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		res.Body = d
		return nil
	}
}

// Text 获取文本
func (res *Res) Text() string {
	return string(res.Body)
}
