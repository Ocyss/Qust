### Go语言动手写爬虫框架 - Qust第六天 支持Http和Socks代理

- 今日代码总行数：531行 
- 较昨日增加: 31行

***仅记录自己学习标准库的探索过程，并不是教程，前期会埋下很多坑。***

### 一.支持Http和Socks代理

使用官方库，"golang.org/x/net/proxy" 来实现Socket5代理

```go
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
```

### 今日内容展示（注释了错误信息）

```go
func main() {
	q := qust.New() //可以指定根网址
	err := q.SetProxy("SOCKS5", "127.0.0.1:10808")
	if err != nil {
		panic(err)
	}
	req := q.Post("https://www.google.com/")
	req.SetUA("Mozilla/5.0")
	res, err := req.Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Text())
}
```

