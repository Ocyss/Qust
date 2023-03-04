package main

import (
	"fmt"
	"qust"
)

func main() {
	q := qust.New() //可以指定根网址
	req := q.Post("https://movie.douban.com/").SetUA("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	res, err := req.Do()
	if err != nil {
		fmt.Println(err)
	}
	html, err := res.Html()
	if err != nil {
		panic(err)
	}
	nodes := html.Find(`//*[@id="billboard"]/div[2]/table/tbody/tr`)
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
