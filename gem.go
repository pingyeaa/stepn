package main

import (
	"encoding/json"
	"fmt"
	"strconv"
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
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=%d&chain=%s&page=%d&refresh=true&sessionID=%s", gType, chain, page, sessionId)
		} else {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=%d&chain=%s&page=%d&refresh=false&sessionID=%s", gType, chain, page, sessionId)
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
			gemCount[data.Otd] = data
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
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=501&gType=%d&chain=%s&page=%d&refresh=true&sessionID=%s", gType, chain, page, sessionId)
		} else {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=501&gType=%d&chain=%s&page=%d&refresh=false&sessionID=%s", gType, chain, page, sessionId)
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

	var url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=%d&chain=%s&page=%d&refresh=true&sessionID=%s", gType, chain, 0, sessionId)
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

func GetPriceBelowNextPrice(nfts map[int]*Shoe) (float64, map[int]int) {

	// get min price
	var minPrice = 999999999999
	for _, nft := range nfts {
		if minPrice > nft.SellPrice {
			minPrice = nft.SellPrice
		}
	}

	// calculate next price
	var nextPrice = 0
	fNextPrice := fmt.Sprintf("%.2f", float64(minPrice+10000)/float64(1000000))
	if p, err := strconv.ParseFloat(fNextPrice, 64); err == nil {
		nextPrice = int(p * 1000000)
	}

	// find nft below nextPrice
	var belowNfts = map[int]*Shoe{}
	for _, nft := range nfts {
		if nft.SellPrice < nextPrice {
			belowNfts[nft.Otd] = nft
		}
	}

	// sort by qType
	// qType => number
	var outNfts = map[int]int{}
	for _, belowNft := range belowNfts {
		if _, ok := outNfts[belowNft.GType]; ok {
			outNfts[belowNft.GType]++
		} else {
			outNfts[belowNft.GType] = 1
		}
	}

	// format out var
	for i := 1; i < 5; i++ {
		if _, ok := outNfts[i]; !ok {
			outNfts[i] = 0
		}
	}

	return float64(nextPrice) / float64(1000000), outNfts
}
