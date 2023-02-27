package main

import (
	"fmt"
	"qust"
)

func main() {
	q := qust.New(qust.BaseUrl("http://t.weather.sojson.com/")) //可以指定根网址
	res, _ := q.Get("api/weather/city/101210101")               //不需要全网址，会自动拼接
	data, _ := res.Json()                                       //获取Json格式数据
	fmt.Println(data.Get("cityInfo").GetVal("city"))            //打印城市信息
	for _, v := range data.Get("data").GetArrays("forecast") {  //循环打印两星期内的天气情况
		fmt.Println(v.GetVal("ymd"), v.GetVal("high"), v.GetVal("low"), v.GetVal("notice"))
	}
	/*
		杭州市
		2023-02-27 高温 12℃ 低温 2℃ 愿你拥有比阳光明媚的心情
		2023-02-28 高温 19℃ 低温 7℃ 阴晴之间，谨防紫外线侵扰
		2023-03-01 高温 19℃ 低温 3℃ 不要被阴云遮挡住好心情
		...
	*/
}
