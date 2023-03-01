### Go语言动手写爬虫框架 - Qust第二天 完善其他方法和Header,Query添加

- 今日代码总行数：392行 
- 较昨日增加: 202行

### 一.修复个小Bug

之前的url拼接判断逻辑有问题，现在改成这样

```go
	if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
		url = string(engine.BaseUrl) + url
	}
```

### 二.编码转换

> 写爬虫最头疼的就是编码问题了，这里使用官方golang.org/x/text包来进行转换

```go
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
```

### 三.完善POST等其他请求方式

使用 `http.NewRequest(method, url, nil)` 方法生成请求体，方便后续，添加Header和Data，并且支持常用的五大协议，`Get`,`Post`,`Delete`,`Put`,`Patch`

```go
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
```

### 四.添加Header

参考了下`gorm`

```go
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
```

### 五.设置Query

代码同上，主要是将值对应到req.Query map[string]any

### 六.请求函数Do

```go
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
```

### 今日内容展示（注释了错误信息）

```go
func main() {
	q := qust.New() //可以指定根网址
	req := q.Post("https://go.apipost.cn/")
	req.AddHeader("head1", "value1")
	req.AddHeader("head2", "value2", "head3", "value3")
	req.AddHeader(map[string]string{"head4": "value4"})
	req.SetQuery("q1", "s1")
	req.SetQuery("q3", 666)
	req.SetQuery(map[string]any{"q4": 1.66})
	res, err := req.Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Text())
}
```

