package main

import (
	"fmt"
	"qust"
)

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
		return
	}
	fmt.Println(res.Text())
}
