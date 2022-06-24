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
	file, err := os.OpenFile(fmt.Sprintf("./db/%s-%s", chain, filePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(value + "\n")
	write.Flush()
}

func FindLatest(filePath string) string {
	file, err := os.OpenFile(fmt.Sprintf("./db/%s-%s", chain, filePath), os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
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

func Rewrite(filePath string, value string) {
	name := fmt.Sprintf("./db/%s-%s", chain, filePath)
	_ = os.Remove(name)
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(value + "\n")
	write.Flush()
}

func GetFileContent(filePath string) string {
	content, err := os.ReadFile(fmt.Sprintf("./db/%s-%s", chain, filePath))
	//file, err := os.OpenFile(fmt.Sprintf("./db/%s-%s", chain, filePath), os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return ""
	}
	//defer file.Close()
	//r := bufio.NewReader(file)
	//content, err := r.ReadByte()
	return string(content)
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

func NumBelowToNext(sneakers map[int]int) (string, int) {
	var prices []int
	for _, price := range sneakers {
		prices = append(prices, price)
	}
	sort.Ints(prices)
	if len(prices) == 0 {
		return "", 0
	}
	minPrice := prices[0]
	nextPrice := (minPrice + 0.2*1000000) / 100000 * 100000
	var count = 0
	for _, price := range prices {
		if price < nextPrice {
			count++
		}
	}
	return fmt.Sprintf("%.2f", float64(nextPrice)/1000000), count
}

func GetTokenPriceForBSC(address string) (float64, float64) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/token_price/binance-smart-chain?contract_addresses=%s&vs_currencies=bnb%%2Cusd", address)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return 0, 0
	}
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	if err != nil {
		log.Println(err.Error())
		return 0, 0
	}
	var data = map[string]map[string]float64{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err.Error())
		return 0, 0
	}
	price := data[address]["usd"]
	priceBnb := data[address]["bnb"]
	return price, priceBnb
}

func GetTokenPriceForSol(address string) float64 {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/token_price/solana?contract_addresses=%s&vs_currencies=usd", address)
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
	var data = map[string]map[string]float64{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	price := data[address]["usd"]
	return price
}

func GetCoinPrice(name string) float64 {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", name)
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
	var data = map[string]map[string]float64{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	price := data[name]["usd"]
	return price
}

func GSTPriceForBSC() (float64, float64) {
	return GetTokenPriceForBSC("0x4a2c860cec6471b9f5f5a336eb4f38bb21683c98")
}

func GMTPriceForBSC() (float64, float64) {
	return GetTokenPriceForBSC("0x3019bf2a2ef8040c242c9a4c5c4bd4c81678b2a1")
}

func GSTPriceForSol() float64 {
	return GetTokenPriceForSol("AFbX8oGjGpmVFywbVouvhQSRmiW2aR1mohfahi4Y2AdB")
}

func GMTPriceForSol() float64 {
	return GetTokenPriceForSol("7i5KKsX2weiTkry7jA4ZwSuXGhs5eJBEjY8vVxR4pfRx")
}

func BnbPrice() float64 {
	return GetCoinPrice("binancecoin")
}

func SolPrice() float64 {
	return GetCoinPrice("solana")
}

func CalcMintProfitForBSC(sneakerFloor float64, scrollFloor float64) (float64, float64, string) {
	gstPrice, _ := GSTPriceForBSC()
	gmtPrice, _ := GMTPriceForBSC()
	bnbPrice := BnbPrice()
	total := gstPrice*360 + gmtPrice*40 + scrollFloor*2*gmtPrice + 20*gstPrice + 10*gmtPrice
	profit := sneakerFloor*bnbPrice*0.94 - total
	return gstPrice, gmtPrice, fmt.Sprintf("%.2fx%.2fx0.94-(%.4fx360+%.4fx40+%.4fx2x%.2f)-(20x%.4f+10x%.4f)=%.2fusd", sneakerFloor, bnbPrice, gstPrice, gmtPrice, gmtPrice, scrollFloor, gstPrice, gmtPrice, profit)
}

func CalcMintProfitForSol(sneakerFloor float64, scrollFloor float64) (float64, float64, string) {
	gstPrice := GSTPriceForSol()
	gmtPrice := GMTPriceForSol()
	solPrice := SolPrice()
	total := gstPrice*360 + gmtPrice*40 + scrollFloor*2*gmtPrice + 20*gstPrice + 10*gmtPrice
	profit := sneakerFloor*solPrice*0.94 - total
	return gstPrice, gmtPrice, fmt.Sprintf("%.2fx%.2fx0.94-(%.4fx360+%.4fx40+%.4fx2x%.2f)-(20x%.4f+10x%.4f)=%.2fusd", sneakerFloor, solPrice, gstPrice, gmtPrice, gmtPrice, scrollFloor, gstPrice, gmtPrice, profit)
}

func GenesShoes() string {
	msg := `\n`
	if chain == "104" {
		msg += `👑 创世数据（BSC）\n`
		msg += `————————————————\n`
	}
	if chain == "103" {
		msg += `️👑 创世数据（Sol）\n`
		msg += `————————————————\n`
	}

	var genesOtd []int
	for _, shoe := range genesShoes {
		genesOtd = append(genesOtd, shoe.Otd)
	}
	sort.Ints(genesOtd)

	var minPrice = 999999999999

	unitName := ""
	if chain == "104" {
		unitName = "BNB"
	}
	if chain == "103" {
		unitName = "Sol"
	}

	for _, otd := range genesOtd {
		for _, shoe := range genesShoes {
			if otd == shoe.Otd {
				color := ""
				if shoe.Quantity == 1 {
					color = "灰"
				}
				if shoe.Quantity == 2 {
					color = "绿"
				}
				if shoe.Quantity == 3 {
					color = "蓝"
				}
				if shoe.Quantity == 4 {
					color = "紫"
				}
				if shoe.Quantity == 5 {
					color = "橙"
				}

				typeName := ""
				if shoe.TypeID == 601 {
					typeName = "W"
				}
				if shoe.TypeID == 602 {
					typeName = "J"
				}
				if shoe.TypeID == 603 {
					typeName = "R"
				}
				if shoe.TypeID == 604 {
					typeName = "T"
				}

				if minPrice > shoe.SellPrice {
					minPrice = shoe.SellPrice
				}

				msg += fmt.Sprintf(`#%d：%s%s，Lv%d，Mint%d，%.2f%s\n`, shoe.Otd, color, typeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000, unitName)
			}
		}
	}

	if len(genesShoes) == 0 {
		msg += `暂无数据\n`
	}
	msg += `————————————————\n`
	msg += fmt.Sprintf(`挂售总数：%d\n`, len(genesOtd))

	prevTotalValue := FindLatest("genes-total.txt")
	if prevTotal, err := strconv.ParseFloat(prevTotalValue, 64); err == nil {
		rate := CalcRate("genes-total.txt", fmt.Sprintf("%d", len(genesShoes)))
		Insert("genes-total.txt", fmt.Sprintf("%d", len(genesShoes)))
		msg += fmt.Sprintf(`新增：%.f｜增幅：%s\n`, float64(len(genesShoes))-prevTotal, rate)
	}

	if len(genesShoes) == 0 {
		msg += fmt.Sprintf(`地板价：0%s`, unitName)
	} else {
		msg += fmt.Sprintf(`地板价：%.2f%s`, float64(minPrice)/1000000, unitName)
	}

	return msg
}
