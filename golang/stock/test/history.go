package main

import (
	"fmt"
	"history"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//http://quotes.money.163.com/service/chddata.html?code=代码&start=开始时间&end=结束时间&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP
//http://flashquote.stock.hexun.com/Quotejs/MA/1_000001_MA.html?
//http://quote.stock.hexun.com/hqzx/industryquotestd.aspx?sorttype=4&page=1&count=50&time=00080

const (
	SHENZHEN = iota
	SHANGHAI
)

var StockCodeR = regexp.MustCompile("Code:[[:space:]]+(?P<code>[[:alnum:]]{6})[[:space:]]+Name")
var HexunDailyK = "http://flashquote.stock.hexun.com/Quotejs/DA/1_%s_DA.html?"
var HexunMarketSummary = "http://quote.stock.hexun.com/hqzx/industryquotestd.aspx?sorttype=4&page=1&count=50&time=00080"

//[19970704,0.00,6.66,7.45,6.66,6.85,53921321,375882768]

var DailyR = regexp.MustCompile(`\[(?P<date>[[:digit:]]{8}),(?P<p1>[[:digit:].]+),(?P<p2>[[:digit:].]+),(?P<p3>[[:digit:].]+),(?P<p4>[[:digit:].]+),(?P<p5>[[:digit:].]+),(?P<v1>[[:digit:]]+),(?P<v2>[[:digit:]]+)\]`)

func main() {
	sh, err := GetStockListByMarket(SHANGHAI)
	if err != nil {
		panic(err)
	}

	/*
		for _, code := range sh {
			<-time.Tick(time.Duration(time.Second * 1))
			GetHistoryDataFromHexun(code)
		}
	*/

	for _, code := range sh {
		name := "asset/" + code
		if !IsHistoryFileValid(name) {
			os.Remove(name)
			continue
		}

		hi, err := GetHistoryFromFile(code)
		if err != nil {
			panic(err)
		}
		hie, err := hi.GetHighest()
		if err != nil {
			panic(err)
		}

		loe, err := hi.GetLowest()
		if err != nil {
			panic(err)
		}

		fmt.Println(hie.Date, hie.High, loe.Date, loe.Low)

	}
	GetMarketSummaryFromHexun()
}

func GetHistoryFromFile(code string) (*history.History, error) {
	content, err := ioutil.ReadFile("asset/" + code)
	if err != nil {
		return nil, err
	}

	days := DailyR.FindAllStringSubmatch(string(content), -1)
	fmt.Printf("%d\n", len(days))

	dis := make([]*history.Daily, 0, len(days))
	for _, day := range days {
		di, err := GetDailyInfoFromString(day)
		if err != nil {
			panic(err)
		}
		dis = append(dis, di)
	}

	return &history.History{
		Info: dis,
	}, nil
}

func GetStockListByMarket(market int) ([]string, error) {
	var file string
	if market == SHENZHEN {
		file = "asset/sz.txt"
	} else if market == SHANGHAI {
		file = "asset/sh.txt"
	} else {
		return nil, fmt.Errorf("Inavlid market: %d", market)
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	stocks := make([]string, 0, 5000)
	fields := strings.Split(string(content), "\n")

	for _, field := range fields {
		matches := StockCodeR.FindStringSubmatch(field)
		if len(matches) == 2 {
			stocks = append(stocks, matches[1])
		}
	}

	return stocks, nil
}

func GetHistoryDataFromHexun(code string) ([]byte, error) {
	url := fmt.Sprintf(HexunDailyK, code)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println(string(body))

	AppendToFile("asset/"+code, body)

	return body, nil
}

func GetMarketSummaryFromHexun() ([]byte, error) {
	resp, err := http.Get(HexunMarketSummary)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println(string(body))

	return body, nil

}

func AppendToFile(name string, data []byte) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(data)
}

func IsHistoryFileValid(name string) bool {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return false
	}

	if strings.Contains(string(content), "DOCTYPE") {
		return false
	}

	return true
}

func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func IsDirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
