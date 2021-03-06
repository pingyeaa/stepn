package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/ini.v1"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var conn *gorm.DB
var err error
var cookie = ""
var cfg *ini.File
var sneakerPrice = map[int]int{}
var newSneakerPrice = map[int]int{}
var itemStatic = map[int]int{}
var chain = "104"
var genesShoes []*Shoe
var genesis23w []*Shoe

func main() {

	cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err.Error())
	}

	key, err := cfg.Section("stepn").GetKey("cookie")
	if err != nil {
		log.Fatalln(err.Error())
	}
	cookie = key.String()

	key, err = cfg.Section("stepn").GetKey("chain")
	if err != nil {
		log.Fatalln(err.Error())
	}
	chain = key.String()

	// https://apilb.stepn.com/run/login?type=2&account=173224989&password=NH5PB87Pgm5PbFxPuKbPBR5l3rvPBhWPuIWhBrAPvPbhbX7FsIvhBo5Qbmvh8y7Q2D5PBrbQ68bhbX7FkPbQ387QBo5lVF7FXDWQsfAP&deviceInfo=model%3AiPhone%23systemVersion%3A15.5%23systemName%3AiOS%23physical%3Atrue%23buildNumber%3A702%23os%3AIOS
	// https://apilb.stepn.com/run/login?type=2&account=173224989&password=sH5PB87Pgm5PbFdPuKbPBR5l38vPBhWPuIWhBrAPvPbhbV7FsIvhBo5QbFvh8y7Q2D5PBrbQv8bhbX7FvPbQ3f7QBo5QVF7FXDWQsfAP&deviceInfo=model%3AiPhone%23systemVersion%3A15.5%23systemName%3AiOS%23physical%3Atrue%23buildNumber%3A702%23os%3AIOS
	// https://apilb.stepn.com/run/login?type=2&account=173224989&password=vH5PB87Pgm5PbFdPuKbPBR5l3rvPBhWPuIWhBrAPvPbhbD7FsIvhBo5QbFvh8y7Q2D5PBrbQu8bhbX7FxPbQ3P7QBoxQVF7FXDWQsfAP&deviceInfo=model%3AiPhone%23systemVersion%3A15.5%23systemName%3AiOS%23physical%3Atrue%23buildNumber%3A702%23os%3AIOS
	// https://apilb.stepn.com/run/login?type=2&account=173224989&password=kHAP2FvQgDbhBHAQ8VbFBKAFshAPxPdPkIxhTKvQN87Q687PkfvQbqvQko7QbDdQTRAP3K7FTHdFsPxQTRxhTIAlWm5Q3hAP68bQXDWh&deviceInfo=model%3AiPhone%23systemVersion%3A15.5%23systemName%3AiOS%23physical%3Atrue%23buildNumber%3A703%23os%3AIOS

	for {

		newSneakerPrice = map[int]int{}
		genesShoes = []*Shoe{}
		genesis23w = []*Shoe{}

		curTime := fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println(curTime)
		writeLog(curTime)
		var allTotal, total = 0, 0
		var msg = ""
		var price float64 = 0
		var minPrice float64 = 999999999
		var rate = ""

		sneakerPriceContent := GetFileContent("sneaker-price.txt")
		_ = json.Unmarshal([]byte(sneakerPriceContent), &sneakerPrice)

		//msg += fmt.Sprintf(`%s\n`, curTime)
		msg += fmt.Sprintf(`???? ??????????????????????????????\n`)
		msg += fmt.Sprintf(`????????????????????????????????????????????????\n`)
		msg += fmt.Sprintf(`??????`)

		total = sneakerTotal(601, 1)
		msg += fmt.Sprintf(`W %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 1)
		msg += fmt.Sprintf(`J %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 1)
		msg += fmt.Sprintf(`R %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 1)
		msg += fmt.Sprintf(`T %d??? \n`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		msg += fmt.Sprintf(`??????`)
		total = sneakerTotal(601, 2)
		msg += fmt.Sprintf(`W %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 2)
		msg += fmt.Sprintf(`J %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 2)
		msg += fmt.Sprintf(`R %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 2)
		msg += fmt.Sprintf(`T %d??? \n`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		msg += fmt.Sprintf(`??????`)
		total = sneakerTotal(601, 3)
		msg += fmt.Sprintf(`W %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 3)
		msg += fmt.Sprintf(`J %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 3)
		msg += fmt.Sprintf(`R %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 3)
		msg += fmt.Sprintf(`T %d??? \n`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		msg += fmt.Sprintf(`??????`)
		total = sneakerTotal(601, 4)
		msg += fmt.Sprintf(`W %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 4)
		msg += fmt.Sprintf(`J %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 4)
		msg += fmt.Sprintf(`R %d???`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 4)
		msg += fmt.Sprintf(`T %d??? \n`, total)
		allTotal += total
		rate = CalcRate("shoe-total.txt", fmt.Sprintf("%d", allTotal))
		Insert("shoe-total.txt", fmt.Sprintf("%d", allTotal))

		newNum, oldNum, avgPrice, middlePrice := CalcDiffNumSneakers(sneakerPrice, newSneakerPrice)

		msg += fmt.Sprintf(`????????? %d????????? %s \n`, allTotal, rate)
		msg += fmt.Sprintf(`???????????? %d????????? %d \n`, newNum, oldNum)
		msg += fmt.Sprintf(`???????????? %s?????????????????? %s \n`, avgPrice, middlePrice)

		if chain == "104" {
			p, num := NumBelowTo(newSneakerPrice)
			msg += fmt.Sprintf(`%sbnb???????????? %d \n`, p, num)
			p, num = NumBelowToNext(newSneakerPrice)
			msg += fmt.Sprintf(`%sbnb???????????? %d \n`, p, num)
		} else {
			p, num := NumBelowTo(newSneakerPrice)
			msg += fmt.Sprintf(`%ssol???????????? %d \n`, p, num)
			p, num = NumBelowToNext(newSneakerPrice)
			msg += fmt.Sprintf(`%ssol???????????? %d \n`, p, num)
		}

		msg += fmt.Sprintf(`\n`)
		if chain == "104" {
			msg += fmt.Sprintf(`???? ??????????????????bnb???\n`)
		} else {
			msg += fmt.Sprintf(`???? ??????????????????sol???\n`)
		}
		msg += fmt.Sprintf(`????????????????????????????????????????????????\n`)
		msg += fmt.Sprintf(`??????`)

		price = floorPrice(601, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`W %.2f???`, price)
		price = floorPrice(602, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`J %.2f???`, price)
		price = floorPrice(603, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`R %.2f???`, price)
		price = floorPrice(604, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`T %.2f??? \n`, price)

		msg += fmt.Sprintf(`??????`)
		price = floorPrice(601, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`W %.2f???`, price)
		price = floorPrice(602, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`J %.2f???`, price)
		price = floorPrice(603, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`R %.2f???`, price)
		price = floorPrice(604, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`T %.2f???\n`, price)

		msg += fmt.Sprintf(`??????`)
		price = floorPrice(601, 3, 1000000)
		msg += fmt.Sprintf(`W %.2f???`, price)
		price = floorPrice(602, 3, 1000000)
		msg += fmt.Sprintf(`J %.2f???`, price)
		price = floorPrice(603, 3, 1000000)
		msg += fmt.Sprintf(`R %.2f???`, price)
		price = floorPrice(604, 3, 1000000)
		msg += fmt.Sprintf(`T %.2f???\n`, price)

		msg += fmt.Sprintf(`??????`)
		price = floorPrice(601, 4, 1000000)
		msg += fmt.Sprintf(`W %.2f???`, price)
		price = floorPrice(602, 4, 1000000)
		msg += fmt.Sprintf(`J %.2f???`, price)
		price = floorPrice(603, 4, 1000000)
		msg += fmt.Sprintf(`R %.2f???`, price)
		price = floorPrice(604, 4, 1000000)
		msg += fmt.Sprintf(`T %.2f???\n`, price)

		rate = CalcRate("shoe-floor.txt", fmt.Sprintf("%f", minPrice))
		Insert("shoe-floor.txt", fmt.Sprintf("%f", minPrice))
		msg += fmt.Sprintf(`???????????? %.2f ????????? %s\n`, minPrice, rate)
		sneakerMinPrice := minPrice

		// ??????
		var scrollTotal = 0

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`???? ??????????????????????????????\n`)
		msg += fmt.Sprintf(`????????????????????????????????????????????????\n`)
		total = sneakerTotal(701, 1)
		msg += fmt.Sprintf(`??? %d???`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 2)
		msg += fmt.Sprintf(`??? %d???`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 3)
		msg += fmt.Sprintf(`??? %d???`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 4)
		msg += fmt.Sprintf(`??? %d???`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 5)
		msg += fmt.Sprintf(`??? %d???\n`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)
		rate = CalcRate("scroll-total.txt", fmt.Sprintf("%d", scrollTotal))
		Insert("scroll-total.txt", fmt.Sprintf("%d", scrollTotal))
		msg += fmt.Sprintf(`?????? %d????????? %s???\n`, scrollTotal, rate)

		// ???????????????

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`???? ??????????????????gmt???\n`)
		msg += fmt.Sprintf(`????????????????????????????????????????????????\n`)

		minPrice = 999999999
		price = floorPrice(701, 1, 100)
		minPrice = comparePrice(minPrice, price)
		scrollMinPrice := minPrice
		msg += fmt.Sprintf(`??? %.2f???`, price)

		price = floorPrice(701, 2, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`??? %.2f???`, price)

		price = floorPrice(701, 3, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`??? %.2f???`, price)

		price = floorPrice(701, 4, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`??? %.2f???`, price)

		price = floorPrice(701, 5, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`??? %.2f???\n`, price)

		rate = CalcRate("scroll-floor.txt", fmt.Sprintf("%f", minPrice))
		Insert("scroll-floor.txt", fmt.Sprintf("%f", minPrice))

		msg += fmt.Sprintf(`???????????? %.2f????????? %s\n`, minPrice, rate)
		msg += fmt.Sprintf(`\n`)

		if chain == "104" {
			msg += fmt.Sprintf(`???? Mint?????????usd???\n`)
			msg += fmt.Sprintf(`????????????????????????????????????????????????\n`)
			gstPrice, gmtPrice, profit := CalcMintProfitForBSC(sneakerMinPrice, scrollMinPrice)
			msg += fmt.Sprintf(`1GST = %.4fU \n`, gstPrice)
			msg += fmt.Sprintf(`1GMT = %.4fU \n`, gmtPrice)
			msg += fmt.Sprintf(`mint?????? = %.4fU \n`, 360*gstPrice+40*gmtPrice)
			msg += fmt.Sprintf(`???????????? = %.4fU \n`, scrollMinPrice*gmtPrice*2)
			msg += fmt.Sprintf(`???????????? = %.4fU \n`, 20*gstPrice+10*gmtPrice)
			msg += fmt.Sprintf(`%s\n`, profit)
		} else {
			msg += fmt.Sprintf(`???? Mint?????????usd???\n`)
			msg += fmt.Sprintf(`????????????????????????????????????????????????\n`)
			gstPrice, gmtPrice, profit := CalcMintProfitForSol(sneakerMinPrice, scrollMinPrice)
			msg += fmt.Sprintf(`1GST = %.4fU \n`, gstPrice)
			msg += fmt.Sprintf(`1GMT = %.4fU \n`, gmtPrice)
			msg += fmt.Sprintf(`mint?????? = %.4fU \n`, 360*gstPrice+40*gmtPrice)
			msg += fmt.Sprintf(`???????????? = %.4fU \n`, scrollMinPrice*gmtPrice*2)
			msg += fmt.Sprintf(`???????????? = %.4fU \n`, 20*gstPrice+10*gmtPrice)
			msg += fmt.Sprintf(`%s\n`, profit)
		}

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`PS?????????????????????????????????????????????????????? \n`)

		push(msg)

		// ?????????????????????????????????
		newSneakerPriceByte, _ := json.Marshal(newSneakerPrice)
		Rewrite("sneaker-price.txt", string(newSneakerPriceByte))
		newSneakerPrice = map[int]int{}
		sneakerPrice = map[int]int{}

		// ?????????
		GenesShoes()

		// ???OG
		if chain == "104" {
			Genesis23wShoes()
		}

		time.Sleep(time.Second * 300)
	}
}

func sneakerTotal(types int, quantity int) int {

	var page = 0
	var url = ""

	for {
		if page == 0 {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&quality=%d&chain=%s&page=%d&refresh=true", types, quantity, chain, page)
		} else {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&quality=%d&chain=%s&page=%d&refresh=false", types, quantity, chain, page)
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
			return 0
		}
		req.Header.Set("cookie", cookie)
		req.Header.Set("accept", "application/json")
		req.Header.Set("accept-language", "zh-CN")
		req.Header.Set("host", "apilb.stepn.com")
		req.Header.Set("group", "173224989")
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer resp.Body.Close()
		respByte, _ := ioutil.ReadAll(resp.Body)
		if string(respByte) == `{"code":102001,"msg":"Player hasnt logged in yet"}` {
			log.Fatalln(`{"code":102001,"msg":"Player hasnt logged in yet"}`)
		}

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

			itemStatic[data.Otd] = 1

			if types != 701 {
				newSneakerPrice[data.Otd] = data.SellPrice
				if data.Otd == 9999 {
					fmt.Println("find you", data.Otd, data.SellPrice)
				}
				if chain == "103" && data.Otd < 10000 {
					data.TypeID = types
					genesShoes = append(genesShoes, data)
				}
				if chain == "104" && data.Otd < 20000 {
					data.TypeID = types
					genesShoes = append(genesShoes, data)
				}
				if chain == "104" && data.Otd < 30000 && data.Otd > 20000 {
					data.TypeID = types
					genesis23w = append(genesis23w, data)
				}
			}
		}

		fmt.Print(".")

		page++
		time.Sleep(time.Second * 2)
	}

	sneakerTotalDesc(types, quantity)

	var total = len(itemStatic)
	itemStatic = map[int]int{}

	return total
}

func sneakerTotalDesc(types int, quantity int) {

	var page = 0
	var url = ""

	for {
		if page == 0 {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=%d&quality=%d&chain=%s&page=%d&refresh=true", types, quantity, chain, page)
		} else {
			url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=%d&quality=%d&chain=%s&page=%d&refresh=false", types, quantity, chain, page)
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		req.Header.Set("cookie", cookie)
		req.Header.Set("accept", "application/json")
		req.Header.Set("accept-language", "zh-CN")
		req.Header.Set("host", "apilb.stepn.com")
		req.Header.Set("group", "173224989")
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer resp.Body.Close()
		respByte, _ := ioutil.ReadAll(resp.Body)
		if string(respByte) == `{"code":102001,"msg":"Player hasnt logged in yet"}` {
			log.Fatalln(`{"code":102001,"msg":"Player hasnt logged in yet"}`)
		}

		var orderList OrderList
		err = json.Unmarshal(respByte, &orderList)
		if err != nil {
			fmt.Println(string(respByte))
			fmt.Println(err.Error())
			return
		}

		if orderList.Data == nil || len(orderList.Data) == 0 {
			break
		}

		var repeatCount = 0

		for _, data := range orderList.Data {

			// ??????????????????5????????????
			_, ok := itemStatic[data.Otd]
			if ok {
				repeatCount++
			} else {
				itemStatic[data.Otd] = 1
			}
			if repeatCount >= 5 {
				break
			}

			if types != 701 {
				newSneakerPrice[data.Otd] = data.SellPrice
				if data.Otd == 9999 {
					fmt.Println("find you", data.Otd, data.SellPrice)
				}
				if chain == "103" && data.Otd < 10000 {
					data.TypeID = types
					genesShoes = append(genesShoes, data)
				}
				if chain == "104" && data.Otd < 20000 {
					data.TypeID = types
					genesShoes = append(genesShoes, data)
				}
				if chain == "104" && data.Otd < 30000 && data.Otd > 20000 {
					data.TypeID = types
					genesis23w = append(genesis23w, data)
				}
			}
		}

		fmt.Print(".")

		page++
		time.Sleep(time.Second)
	}

	return
}

func floorPrice(types int, quantity int, zeroNum int) float64 {

	time.Sleep(time.Second * 1)

	var url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&quality=%d&chain=%s&page=%d&refresh=true", types, quantity, chain, 0)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	req.Header.Set("cookie", cookie)
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("host", "apilb.stepn.com")
	req.Header.Set("group", "173224989")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	if string(respByte) == `{"code":102001,"msg":"Player hasnt logged in yet"}` {
		log.Fatalln(`{"code":102001,"msg":"Player hasnt logged in yet"}`)
	}

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

	price := float64(orderList.Data[0].SellPrice) / float64(zeroNum)
	fmt.Print(".")

	return price
}

func comparePrice(price1 float64, price2 float64) float64 {
	if fmt.Sprintf("%.2f", price2) == "0.00" {
		return price1
	}
	if price2 == 0 {
		return price1
	}
	if price1 > price2 {
		return price2
	}
	return price1
}

func writeLog(content string) string {
	filePath := "./logs/" + time.Now().Format("2006-01-02") + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("??????????????????", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(content + "\n")
	write.Flush()
	return content + `\n`
}

func push(msg string) {

	var webhook *ini.Key
	if chain == "104" {
		webhook, err = cfg.Section("discord").GetKey("webhook")
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else {
		webhook, err = cfg.Section("discord").GetKey("sol_webhook")
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	content := []byte(fmt.Sprintf(`{"content":"%s"}`, msg))
	fmt.Println(string(content))
	req, err := http.NewRequest("POST", webhook.String(), bytes.NewBuffer(content))
	if err != nil {
		log.Fatalln(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respByte))
}

func pushDcFromConfigKey(configKey string, msg string) {
	webhook, err := cfg.Section("discord").GetKey(configKey)
	if err != nil {
		log.Fatalln(err.Error())
	}
	content := []byte(fmt.Sprintf(`{"content":"%s"}`, msg))
	fmt.Println(string(content))
	req, err := http.NewRequest("POST", webhook.String(), bytes.NewBuffer(content))
	if err != nil {
		log.Fatalln(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respByte))
}

func pushToGenes(msg string) {
	pushDcFromConfigKey("genes_webhook", msg)
}

func pushToGenesis23w(msg string) {
	pushDcFromConfigKey("genesis23w_webhook", msg)
}
