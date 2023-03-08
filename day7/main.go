package main

import (
	"fmt"
	"qust"
)

func main() {
	q := qust.New() //可以指定根网址
	err := q.SetProxy("HTTP", "127.0.0.1:10809")
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
