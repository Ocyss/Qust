### Go语言动手写爬虫框架 - Qust第五天 完善File方法

- 今日代码总行数：500行 
- 较昨日增加: 19行

***仅记录自己学习标准库的探索过程，并不是教程，前期会埋下很多坑。***

### 一.保存文件

简单封装下

```go
type Data struct {
	data []byte
}

func New(data []byte) *Data {
	return &Data{data}
}

func (f *Data) Save(filename string) error {
	err := os.WriteFile(filename, f.data, 0644)
	return err
}
```

### 今日内容展示（注释了错误信息）

调用随机图片api，进行图片保存

```go
func main() {
	q := qust.New() //可以指定根网址
	req := q.Post("https://imgapi.xl0408.top/index.php")
	req.SetUA("Mozilla/5.0")
	res, err := req.Do()
	if err != nil {
		fmt.Println(err)
	}
	err = res.File().Save("./1.png")
	if err != nil {
		fmt.Println(err)
	}
}
```

