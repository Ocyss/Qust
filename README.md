<p align="center">
  <a href="https://github.com/qiu-lzsnmb/qust">
    <img src="https://qiu-blog.oss-cn-hangzhou.aliyuncs.com/Article/1677503330021120612.png" alt="Logo" width="180" height="180">
  </a>
  <h1 align="center">Qust</h1>
  <h3 align="center">15天从0实现的爬虫框架</h3>
> 听说Go语言是出了名的造轮子语言，打算新写个项目（造轮子），项目主旨就是简单，轻量，方便的爬虫。
>


下面是标准库进行简单的`JSON`爬取，不管是解析`JSON`，还是获取数据，都异常麻烦，需要写大量的错误判断，和解析代码，面对复杂的`JSON`更是头疼，需要一直反射还可能会遇到错误，这个框架就是为了简化这些操作，让代码简洁并且容错高，错误处理方便。

> 写法会偏向`Python`的`requests`库，毕竟学的第一语言就是`Python`，而且还是直接从`requests`爬虫入门的

```go
func main() {
	url := "http://t.weather.sojson.com/api/weather/city/101210101"
	html, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer html.Body.Close()
	jsonData, _ := io.ReadAll(html.Body)
	var v interface{}
	_ = json.Unmarshal(jsonData, &v)
	data := v.(map[string]interface{})
	fmt.Println(data["cityInfo"].(map[string]interface{})["city"], data["data"].(map[string]interface{})["forecast"])
}
```

### Qust 框架

为了方便了解go的标准库，又不枯燥，所以打算慢慢的写一个项目入手，然后去看go的源码。算是一边学习，一边总结经验和写教程，主要受了[极客兔兔大佬七天系列](https://geektutu.com/post/gee.html)影响，看了`Gee`感觉很有意思，但模仿着写`Gee`又没提升，所以自己探索！

> 也会去参考大佬们的源码，尽量全部使用标准库

#### 这个库的计划 >

- [ ] 格式支持
  - [x] Json
  - [x] Html
  - [x] File
  - [ ] Xml
  - [ ] Csv
- [ ] 便捷性
  - [x] 可以指定基础路径，减少后续的url长度
  - [x] 可以进行分组，不同分组基础路径可以不同
  - [ ] 高并发
  - [x] 模拟表单Form提交
- [ ] 反爬
  - [ ] 代理池
  - [x] 方便的添加headers
  - [ ] 随机生成ua
  - [ ] 可等待js渲染完在爬取html
  - [ ] ck缓存
  - [ ] 支持robots.txt

## 目录

- 第一天：
- 第二天：
- 第三天：
- 第四天：
- 第五天：
- 第六天：
- 第七天：
- 第八天：
- 第九天：
- 第十天：
- 第十一天：
- 第十二天：
- 第十三天：
- 第十四天：
- 第十五天：