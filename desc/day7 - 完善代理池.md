### Go语言动手写爬虫框架 - Qust第七天 完善代理池

- 今日代码总行数：559行 
- 较昨日增加: 28行

***仅记录自己学习标准库的探索过程，并不是教程，前期会埋下很多坑。***

### 一.代理池结构体

只实现了Http代理，分为两种模式，一种迭代，一种随机

```go
type Proxys struct {
	Proxy []string
	Index int
	Size  int
	Mode  int
}

func (engine *Engine) SetProxys(proxys []string, mode int) *Proxys {
	p := &Proxys{proxys, 0, len(proxys), mode}
	engine.Proxys = p
	return p
}
```

### 二.随机代理

每次生成请求的时候，随机一个代理字符串，传进去，在发起请求的时候解析到客户端，进行代理请求

```go
func (engine *Engine) Ask(method string, url string) *Req {
   //url拼接
   if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
      url = string(engine.BaseUrl) + url
   }
   u, _ := urlpkg.Parse(url)
   ************************************
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
   *************************************
   //生成一个请求体
   req := &Req{u, method, map[string]any{}, engine.Client, nil, make(http.Header), p}
   return req
}
```

```go
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
   **********************************************
   if req.Proxy != "" {
      urlproxy, err := url.Parse(req.Proxy)
      if err != nil {
         return nil, err
      }
      client.Transport = &http.Transport{Proxy: http.ProxyURL(urlproxy)}
   }
   **********************************************
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
```

### 今日内容展示（注释了错误信息）

```go
var wg sync.WaitGroup

func get(q *qust.Engine, name int) {
	defer wg.Done()
	req := q.Post("https://imgapi.xl0408.top/index.php")
	req.SetUA("Mozilla/5.0")
	res, err := req.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		err = res.File().Save(fmt.Sprintf("./%d.png", name))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	q := qust.New() //可以指定根网址
	q.SetProxys(
		[]string{
			"http://127.0.0.1:10809",
		}, 1)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		time.Sleep(time.Second)
		go get(q, i)
	}
	wg.Wait()
}
```

