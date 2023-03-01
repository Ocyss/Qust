### Go语言动手写爬虫框架 - Qust第一天 处理简单json

今日代码总行数：191行

### 一.实现qust.New函数

Go原生不支持参数可选，通过...any模拟

```go
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
```

### 二.实现qust.Get接口

> 字符串拼接为性能考虑使用了strings.Builder()，可参考 [极客兔兔大佬的文章](https://geektutu.com/post/hpg-string-concat.html)
>
> 这我也没测试直接+，和Builder哪个快，无所谓小问题

```go
func (engine *Engine) Get(url string) (*Res, error) {
	//var builder strings.Builder
	//builder.WriteString(string(engine.BaseUrl))
	//if len(url) > 8 && url[:7] != "http://" && url[:8] != "https://" {
	//	builder.WriteString(url)
	//}
    //url拼接判断逻辑有问题，现在改成这样day2
	if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
		url = string(engine.BaseUrl) + url
	}
    //进行标准库请求
	html, err := http.Get(url)
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
```

### 三.Json数据处理

代码借鉴了[simplejson](https://github.com/bitly/go-simplejson) 但肯定没人家功能多，不需要删除和设置，爬虫只要能读，问题就不大

```go
func (res *Res) Json() (*qjson.Data, error) {
	return qjson.New(res.Body) //获取json格式数据
}

func New(body []byte) (*Data, error) {
	j := new(Data) 
	err := json.Unmarshal(body, &j.data) //标准库解析json
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (j *Data) Get(key string) *Data {
	m, err := j.Map() //判断是不是map[string]any
	if err == nil {
		if val, ok := m[key]; ok {
			return &Data{val}
		}
	}
	return &Data{nil}
}
```

> 剩下几个获取值的函数就不提了，就反射射射射···完事了

### 今日内容展示（注释了错误信息）

```go
func main() {
	q := qust.New(qust.BaseUrl("http://t.weather.sojson.com/")) //可以指定根网址
	res, _ := q.Get("api/weather/city/101210101")               //不需要全网址，会自动拼接
	data, _ := res.Json()                                       //获取Json格式数据
	fmt.Println(data.Get("cityInfo").GetVal("city"))            //打印城市信息
	for _, v := range data.Get("data").GetArrays("forecast") {  //循环打印两星期内的天气情况
		fmt.Println(v.GetVal("ymd"), v.GetVal("high"), v.GetVal("low"), v.GetVal("notice"))
	}
	/*
		杭州市
		2023-02-27 高温 12℃ 低温 2℃ 愿你拥有比阳光明媚的心情
		2023-02-28 高温 19℃ 低温 7℃ 阴晴之间，谨防紫外线侵扰
		2023-03-01 高温 19℃ 低温 3℃ 不要被阴云遮挡住好心情
		...
	*/
}
```

