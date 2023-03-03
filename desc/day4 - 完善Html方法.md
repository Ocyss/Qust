### Go语言动手写爬虫框架 - Qust第四天 完善Html方法

- 今日代码总行数：481行 
- 较昨日增加: 21行

***仅记录自己学习标准库的探索过程，并不是教程，前期会埋下很多坑。***

### 一.解析

这里考虑到用原生库写html解析会很麻烦，要写的比较多，所以进行封装

使用的第三方库 github.com/antchfx/htmlquery

> 之前Python也经常用xpath，也很方便，浏览器F12就可以直接复制使用

```go
type Data struct {
   *html.Node
}

func New(body *html.Node) *Data {
   return &Data{Node: body}
}
```

这里选择用继承的方式，然后扩展一些方法

### 二.根据xpath获取Node

因为是继承，所以这里也需要麻烦点

```go
// Find 查找全部
func (h *Data) Find(expr string) (res []*Data) {
	list := htmlquery.Find(h.Node, expr)
	res = make([]*Data, len(list))
	for i, v := range list {
		res[i] = &Data{Node: v}
	}
	return
}

// FindOne 查找单个
func (h *Data) FindOne(expr string) *Data {
	node := htmlquery.FindOne(h.Node, expr)
	return &Data{Node: node}
}
```

### 三.获取内容和属性

```go
// Text 获取文本数据
func (h *Data) Text() string {
	return htmlquery.InnerText(h.Node)
}

// Attr 获取属性
func (h *Data) Attr(name string) string {
	return htmlquery.SelectAttr(h.Node, name)
}
```

### 今日内容展示（注释了错误信息）

爬取豆瓣一周口碑榜

```go
func main() {
	q := qust.New()
	req := q.Post("https://movie.douban.com/").SetUA("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	res, _ := req.Do()//发送请求
	html, _ := res.Html()//获取Html节点
	nodes := html.Find(`//*[@id="billboard"]/div[2]/table/tbody/tr`) //直接右键获取的xpath
	for _, node := range nodes {
		fmt.Println(node.Text())
	}
	/*
	   1鲸
	   2千寻小姐
	   3惠子，凝视
	   4正义回廊
	 */
}
```

