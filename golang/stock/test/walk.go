package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//func Walk(root string, walkFn WalkFunc) error
//type WalkFunc func(path string, info os.FileInfo, err error) error
func main() {
	filepath.Walk("asset", func(path string, info os.FileInfo, err error) error {
		file := filepath.Base(path)
		if strings.ToUpper(file) != file {
			return nil
		}

		fmt.Printf("%+v\n", GetStockHistoryFromLocal(file))
		return nil
	})
}

type Stock struct {
	Date           string  `csv:"日期"`
	Code           string  `csv:"股票代码"`
	Name           string  `csv:"名称"`
	TClose         float64 `csv:"收盘价"`
	High           float64 `csv:"最高价"`
	Low            float64 `csv:"最低价"`
	TOpen          float64 `csv:"开盘价"`
	PClose         float64 `csv:"前收盘"`
	Increase       float64 `csv:"涨跌额"`
	IncreaseRate   float64 `csv:"涨跌幅"`
	TurnOverRate   float64 `csv:"换手率"`
	TurnOverVolume float64 `csv:"成交量"`
	TurnOverValue  float64 `csv:"成交金额"`
	TotalValue     float64 `csv:"总市值"`
	InMarketValue  float64 `csv:"流通市值"`
}

func (s *Stock) String() string {
	return fmt.Sprintf("%s %s %f %f %f %f", s.Code, s.Date, s.TOpen, s.TClose, s.High, s.Low)
}

func GetStockHistoryFromLocal(code string) []*Stock {
	name := "asset/" + code

	c, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(c), "\n")

	sts := make([]*Stock, 0, len(lines))

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if i == 0 {
			continue
		}

		st := ParseStock(line)

		sts = append(sts, st)
	}

	return sts
}

func ParseStock(str string) *Stock {
	var st Stock
	str = strings.TrimSpace(str)
	fields := strings.Split(str, ",")

	for j, f := range fields {
		var err error
		switch j {

		case 0:
			st.Date = f
		case 1:
			st.Code = strings.Trim(f, "'")
		case 2:
		case 3:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.TClose, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 4:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.High, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 5:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.Low, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 6:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.TOpen, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 7:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.PClose, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 8:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.Increase, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 9:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.IncreaseRate, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 10:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.TurnOverRate, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 11:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.TurnOverVolume, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 12:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.TurnOverValue, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 13:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.TotalValue, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		case 14:
			f = strings.TrimSpace(f)
			if f == "None" {
				f = "0.0"
			}
			st.InMarketValue, err = strconv.ParseFloat(f, 64)
			if err != nil {
				panic(err)
			}
		}
	}

	return &st
}
