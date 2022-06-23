package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"gopkg.in/fatih/set.v0"
)

func Insert(filePath string, value string) {
	file, err := os.OpenFile(fmt.Sprintf("./db/%s", filePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(value + "\n")
	write.Flush()
}

func FindLatest(filePath string) string {
	file, err := os.OpenFile(fmt.Sprintf("./db/%s", filePath), os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lineText string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText = scanner.Text()
	}
	return lineText
}

func CalcRate(filePath string, currentValue string) string {
	value := FindLatest(filePath)
	if value == "" {
		return ""
	}
	if cur, err := strconv.ParseFloat(currentValue, 64); err == nil {
		if s, err := strconv.ParseFloat(value, 64); err == nil {
			diff := cur - s
			return fmt.Sprintf("%.2f%%", diff/s*100)
		}
	}

	return ""
}

func CalcDiffNumSneakers(old map[int]int, new map[int]int) (int, int, string, string) {
	oldSet := set.New(set.ThreadSafe)
	newSet := set.New(set.ThreadSafe)
	for id, _ := range old {
		oldSet.Add(id)
	}
	for id, _ := range new {
		newSet.Add(id)
	}
	consumes := set.Difference(oldSet, newSet).List()
	news := set.Difference(newSet, oldSet).List()

	var totalPrice float64
	for _, id := range news {
		sid := id.(int)
		price := float64(new[sid]) / 1000000
		totalPrice += price
	}
	avgPrice := fmt.Sprintf("%.4f", totalPrice/float64(len(news)))

	var middlePrice string
	var priceList []int
	for _, id := range news {
		priceList = append(priceList, new[id.(int)])
	}
	sort.Ints(priceList)
	for k, price := range priceList {
		if k == len(priceList)/2 {
			middlePrice = fmt.Sprintf("%.4f", float64(price)/1000000)
		}
	}

	return len(news), len(consumes), avgPrice, middlePrice
}

func NumBelowTo(sneakers map[int]int) (string, int) {
	var prices []int
	for _, price := range sneakers {
		prices = append(prices, price)
	}
	sort.Ints(prices)
	if len(prices) == 0 {
		return "", 0
	}
	minPrice := prices[0]
	nextPrice := (minPrice + 0.1*1000000) / 100000 * 100000
	var count = 0
	for _, price := range prices {
		if price < nextPrice {
			count++
		}
	}
	return fmt.Sprintf("%.2f", float64(nextPrice)/1000000), count
}

func GetTokenPrice(address string) float64 {
	url := fmt.Sprintf("https://api.pancakeswap.info/api/v2/tokens/%s", address)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	var data struct {
		Data struct {
			Name     string
			Symbol   string
			Price    string
			PriceBnb string `json:"price_BNB"`
		}
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	price, err := strconv.ParseFloat(data.Data.PriceBnb, 64)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return price
}

func GSTPriceForBSC() float64 {
	return GetTokenPrice("0x4a2c860cec6471b9f5f5a336eb4f38bb21683c98")
}

func GMTPriceForBSC() float64 {
	return GetTokenPrice("0x3019bf2a2ef8040c242c9a4c5c4bd4c81678b2a1")
}

func CalcMintProfitForBSC(sneakerFloor float64, scrollFloor float64) string {
	gstPrice := GSTPriceForBSC()
	gmtPrice := GMTPriceForBSC()
	total := gstPrice*360 + gmtPrice*40 + scrollFloor*2*gmtPrice
	log.Println(gstPrice)
	profit := sneakerFloor*0.94 - total
	return fmt.Sprintf("%.2f*0.94-(%.4f*360+%.4f*40+%.4f*2*%.2f)=%.2f", sneakerFloor, gstPrice, gmtPrice, gmtPrice, scrollFloor, profit)
}
