package main

import (
	"encoding/json"
	"fmt"
	"time"
)

var gemHandled = map[int]*Shoe{}

// GemTotal types 501
// gType, 1-Efficiency, 2-Luck, 3-Comfort, 4-Resilience
func GemTotal(gType int) int {

	var gemCount = map[int]*Shoe{}

	var page = 0
	var url = ""

	for {
		if page == 0 {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=%d&chain=%s&page=%d&refresh=true", gType, chain, page)
		} else {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=%d&chain=%s&page=%d&refresh=false", gType, chain, page)
		}
		respByte := Get(url)
		var orderList OrderList
		err = json.Unmarshal(respByte, &orderList)
		if err != nil {
			fmt.Println(string(respByte))
			fmt.Println(err.Error())
			return 0
		}

		if orderList.Data == nil || len(orderList.Data) == 0 {
			break
		}

		for _, data := range orderList.Data {
			data.GType = gType
			gemHandled[data.Otd] = data
		}

		fmt.Print(".")

		page++
		time.Sleep(time.Second * 2)
	}

	gemCount = GemTotalDesc(gType, gemCount)

	var total = len(gemCount)

	return total
}

func GemTotalDesc(gType int, gemCount map[int]*Shoe) map[int]*Shoe {

	var page = 0
	var url = ""

	for {
		if page == 0 {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=501&gType=%d&chain=%s&page=%d&refresh=true", gType, chain, page)
		} else {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=501&gType=%d&chain=%s&page=%d&refresh=false", gType, chain, page)
		}
		respByte := Get(url)
		var orderList OrderList
		err = json.Unmarshal(respByte, &orderList)
		if err != nil {
			fmt.Println(string(respByte))
			fmt.Println(err.Error())
			return nil
		}

		if orderList.Data == nil || len(orderList.Data) == 0 {
			break
		}

		var repeatCount = 0

		for _, data := range orderList.Data {
			_, ok := gemHandled[data.Otd]
			if ok {
				repeatCount++
			} else {
				data.GType = gType
				gemHandled[data.Otd] = data
				gemCount[data.Otd] = data
			}
			if repeatCount >= 5 {
				break
			}
		}

		fmt.Print(".")

		page++
		time.Sleep(time.Second)
	}

	return gemCount
}

func GemFloorPrice(gType int) float64 {

	time.Sleep(time.Second * 1)

	var url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=%d&chain=%s&page=%d&refresh=true", gType, chain, 0)
	respByte := Get(url)

	var orderList OrderList
	err = json.Unmarshal(respByte, &orderList)
	if err != nil {
		fmt.Println(string(respByte))
		fmt.Println(err.Error())
		return 0
	}

	if orderList.Data == nil || len(orderList.Data) == 0 {
		return 0
	}

	price := float64(orderList.Data[0].SellPrice) / float64(1000000)
	fmt.Print(".")

	return price
}
