package main

import (
	"fmt"
	"qust"
	"sync"
	"time"
)

var wg sync.WaitGroup

func get(q *qust.Engine, name int) {
	defer wg.Done()
	req := q.Post("https://imgapi.xl0408.top/index.php")
	req.SetUA("Mozilla/5.0")
	res, err := req.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		err = res.File().Save(fmt.Sprintf("./%d.png", name))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	q := qust.New() //可以指定根网址
	q.SetProxys(
		[]string{
			"http://127.0.0.1:10809",
		}, 1)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		time.Sleep(time.Second)
		go get(q, i)
	}
	wg.Wait()
}
