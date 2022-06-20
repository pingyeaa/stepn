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

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var conn *gorm.DB
var err error
var cookie = "SESSIONIDD2=344285707892823049:1655692521308:14338088; Path=/"

func main() {

	// https://apilb.stepn.com/run/login?type=2&account=173224989&password=NH5PB87Pgm5PbFxPuKbPBR5l3rvPBhWPuIWhBrAPvPbhbX7FsIvhBo5Qbmvh8y7Q2D5PBrbQ68bhbX7FkPbQ387QBo5lVF7FXDWQsfAP&deviceInfo=model%3AiPhone%23systemVersion%3A15.5%23systemName%3AiOS%23physical%3Atrue%23buildNumber%3A702%23os%3AIOS
	// https://apilb.stepn.com/run/login?type=2&account=173224989&password=sH5PB87Pgm5PbFdPuKbPBR5l38vPBhWPuIWhBrAPvPbhbV7FsIvhBo5QbFvh8y7Q2D5PBrbQv8bhbX7FvPbQ3f7QBo5QVF7FXDWQsfAP&deviceInfo=model%3AiPhone%23systemVersion%3A15.5%23systemName%3AiOS%23physical%3Atrue%23buildNumber%3A702%23os%3AIOS
	// https://apilb.stepn.com/run/login?type=2&account=173224989&password=vH5PB87Pgm5PbFdPuKbPBR5l3rvPBhWPuIWhBrAPvPbhbD7FsIvhBo5QbFvh8y7Q2D5PBrbQu8bhbX7FxPbQ3P7QBoxQVF7FXDWQsfAP&deviceInfo=model%3AiPhone%23systemVersion%3A15.5%23systemName%3AiOS%23physical%3Atrue%23buildNumber%3A702%23os%3AIOS

	for {

		curTime := fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println(curTime)
		writeLog(curTime)
		var allTotal, total = 0, 0
		var msg = ""
		var price float64 = 0
		var minPrice float64 = 999999999

		//msg += fmt.Sprintf(`%s\n`, curTime)
		msg += fmt.Sprintf(`👟 鞋子数量\n`)
		msg += fmt.Sprintf(`————————————————\n`)
		msg += fmt.Sprintf(`灰｜`)

		total = sneakerTotal(601, 1)
		msg += fmt.Sprintf(`W %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 1)
		msg += fmt.Sprintf(`J %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 1)
		msg += fmt.Sprintf(`R %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 1)
		msg += fmt.Sprintf(`T %d｜ \n`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		msg += fmt.Sprintf(`绿｜`)
		total = sneakerTotal(601, 2)
		msg += fmt.Sprintf(`W %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 2)
		msg += fmt.Sprintf(`J %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 2)
		msg += fmt.Sprintf(`R %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 2)
		msg += fmt.Sprintf(`T %d｜ \n`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		msg += fmt.Sprintf(`蓝｜`)
		total = sneakerTotal(601, 3)
		msg += fmt.Sprintf(`W %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 3)
		msg += fmt.Sprintf(`J %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 3)
		msg += fmt.Sprintf(`R %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 3)
		msg += fmt.Sprintf(`T %d｜ \n`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		msg += fmt.Sprintf(`紫｜`)
		total = sneakerTotal(601, 4)
		msg += fmt.Sprintf(`W %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 4)
		msg += fmt.Sprintf(`J %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 4)
		msg += fmt.Sprintf(`R %d｜`, total)
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 4)
		msg += fmt.Sprintf(`T %d｜ \n`, total)
		allTotal += total
		msg += fmt.Sprintf(`总鞋数：%d\n`, allTotal)

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`💰 鞋子地板价（bnb）\n`)
		msg += fmt.Sprintf(`————————————————\n`)
		msg += fmt.Sprintf(`灰｜`)

		price = floorPrice(601, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`W %.2f｜`, price)
		price = floorPrice(602, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`J %.2f｜`, price)
		price = floorPrice(603, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`R %.2f｜`, price)
		price = floorPrice(604, 1, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`T %.2f｜ \n`, price)

		msg += fmt.Sprintf(`绿｜`)
		price = floorPrice(601, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`W %.2f｜`, price)
		price = floorPrice(602, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`J %.2f｜`, price)
		price = floorPrice(603, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`R %.2f｜`, price)
		price = floorPrice(604, 2, 1000000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`T %.2f｜\n`, price)

		msg += fmt.Sprintf(`蓝｜`)
		price = floorPrice(601, 3, 1000000)
		msg += fmt.Sprintf(`W %.2f｜`, price)
		price = floorPrice(602, 3, 1000000)
		msg += fmt.Sprintf(`J %.2f｜`, price)
		price = floorPrice(603, 3, 1000000)
		msg += fmt.Sprintf(`R %.2f｜`, price)
		price = floorPrice(604, 3, 1000000)
		msg += fmt.Sprintf(`T %.2f｜\n`, price)

		msg += fmt.Sprintf(`紫｜`)
		price = floorPrice(601, 4, 1000000)
		msg += fmt.Sprintf(`W %.2f｜`, price)
		price = floorPrice(602, 4, 1000000)
		msg += fmt.Sprintf(`J %.2f｜`, price)
		price = floorPrice(603, 4, 1000000)
		msg += fmt.Sprintf(`R %.2f｜`, price)
		price = floorPrice(604, 4, 1000000)
		msg += fmt.Sprintf(`T %.2f｜\n`, price)
		msg += fmt.Sprintf(`全网地板：%.2f\n`, minPrice)

		// 卷轴
		var scrollTotal = 0

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`📜 卷轴数量\n`)
		msg += fmt.Sprintf(`————————————————\n`)
		total = sneakerTotal(701, 1)
		msg += fmt.Sprintf(`灰 %d｜`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 2)
		msg += fmt.Sprintf(`绿 %d｜`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 3)
		msg += fmt.Sprintf(`蓝 %d｜`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 4)
		msg += fmt.Sprintf(`紫 %d｜`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 5)
		msg += fmt.Sprintf(`橙 %d｜\n`, total)
		scrollTotal += total
		time.Sleep(time.Second * 5)
		msg += fmt.Sprintf(`合计：%d\n`, scrollTotal)

		// 卷轴地板价

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`💰 卷轴地板价（gmt）\n`)
		msg += fmt.Sprintf(`————————————————\n`)

		minPrice = 999999999
		price = floorPrice(701, 1, 10000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`灰 %.2f｜`, price)

		price = floorPrice(701, 2, 10000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`绿 %.2f｜`, price)

		price = floorPrice(701, 3, 10000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`蓝 %.2f｜`, price)

		price = floorPrice(701, 4, 10000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`紫 %.2f｜`, price)

		price = floorPrice(701, 5, 10000)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`橙 %.2f｜\n`, price)

		msg += fmt.Sprintf(`全网地板：%.2f\n`, minPrice)

		go push(msg)

		time.Sleep(time.Second * 600)
	}
}

func sneakerTotal(types int, quantity int) int {

	var page, total = 0, 0

	for {
		var url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&quality=%d&chain=104&page=%d&refresh=false", types, quantity, page)
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
			log.Fatalln(err.Error())
		}

		if orderList.Data == nil || len(orderList.Data) == 0 {
			break
		}

		total += len(orderList.Data)

		fmt.Print(".")

		page++
		time.Sleep(time.Second)
	}

	return total
}

func floorPrice(types int, quantity int, zeroNum int) float64 {

	time.Sleep(time.Second * 1)

	var url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&quality=%d&chain=104&page=%d&refresh=true", types, quantity, 0)
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
		log.Fatalln(err.Error())
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
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(content + "\n")
	write.Flush()
	return content + `\n`
}

func push(msg string) {
	content := []byte(fmt.Sprintf(`{"content":"%s"}`, msg))
	fmt.Println(string(content))
	var url = fmt.Sprintf("https://discord.com/api/webhooks/987903832946262137/EG10I7wB5rCWxB7--auYlcnxRQtxdyIF7Z3Q3OQQfNdqv3qYyt4RQQA0tnqurQ92iSWE")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
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
