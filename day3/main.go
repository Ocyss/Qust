package main

import (
	"fmt"
	"qust"
)

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
