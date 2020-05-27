package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	//resp, err := http.Get("https://www.baidu.com/s?wd=broadcom")
	//resp, err := http.Get("https://www.baidu.com/s?wd=股票历史数据")
	//resp, err := http.Get("http://www.baidu.com/link?url=omjAUk9ydFjsgmzPM1lftoIljI6kGkJyfSqhJ1Wzk5j6LMepxYnurWFQwl9xabw1v1zhyee0Rispyotpi_kV7gONoeXw59niHQqoYSbVYCu")
	//resp, err := http.Get("https://blog.csdn.net/luanpeng825485697/article/details/78442062")
	//resp, err := http.Get("http://data.gtimg.cn/flashdata/hushen/latest/daily/sz000002.js?maxage=10000000000&amp;visitDstTime=1")
	resp, err := http.Get("http://data.gtimg.cn/flashdata/hushen/latest/daily/sz300002.js?maxage=43201&amp;visitDstTime=1")

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(body))
}
