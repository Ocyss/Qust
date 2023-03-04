package main

import (
	"fmt"
	"qust"
)

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
