### Go语言动手写爬虫框架 - Qust第三天 设置数据Body和cookie

- 今日代码总行数：460行 
- 较昨日增加: 68行

***仅记录自己学习标准库的探索过程，并不是教程，前期会埋下很多坑。***

### 一.优化 SetQuery

今天改成下面代码，方便支持更多类型

```go
//例：
req.SetQuery("q1", "s1")
req.SetQuery("q2", "s2", "q3", 666)
req.SetQuery(map[string]float64{"q4": 1.66})
req.SetQuery(map[string]int{"q5": 666})
req.SetQuery(map[string]bool{"q6": false})
```

```go
func (req *Req) SetQuery(k any, args ...any) *Req {
	//标准 2 参数，直接添加放回
	if s, ok := k.(string); ok && len(args) == 1 {
		req.Query[s] = args[0]
		return req
	}
	//多参数，或者结构体，进行遍历
	s1 := ""
	args = append([]any{k}, args...)
	for _, arg := range args {
		if s2, ok := arg.(string); ok {
			if s1 == "" {
				s1 = s2
			} else {
				req.Query[s1] = s2
				s1 = ""
			}
		} else if r := reflect.ValueOf(arg); r.Kind() == reflect.Map { 
            //判断是不是map类型，然后进行添加
			iter := r.MapRange()
			for iter.Next() {
				req.Query[fmt.Sprintf("%v", iter.Key())] = iter.Value()
			}
		}
	}
	return req
}
```

### 二.修改Json的String方法

让现在可以直接打印res.Json()数据格式化输出

```go
func (j *Data) String() string {
   var out bytes.Buffer
   res, err := json.Marshal(j.data)
   if err != nil {
      return ""
   }
   _ = json.Indent(&out, res, "", "\t")
   return out.String()
}
```

### 三.设置请求数据Body

支持了JSON，FROM，XML，MultipartForm 四种格式，也可以自定义格式

```go
type DataType string
const (
	JSON          DataType = "application/json"
	FROM          DataType = "multipart/form-data"
	XML           DataType = "text/xml"
	MultipartForm DataType = "application/x-www-form-urlencoded"
)
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
```

### 四.设置cookie

```go
func (req *Req) SetCookie(data any) *Req {
   switch data.(type) {
   case string:
      req.AddHeader("cookie", data.(string))
   }
   return req
}
```



### 今日内容展示（注释了错误信息）

```go
func main() {
	q := qust.New() //可以指定根网址
	req := q.Post("https://go.apipost.cn/").SetCookie("token=123456")
	err := req.SetData(qust.JSON, map[string]any{"password": "123",
		"username": "root"})
	if err != nil {
		fmt.Println(err)
	}
	res, err := req.Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Json())
}
```

