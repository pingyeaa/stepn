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

		curTime := fmt.Sprintf("\n%s", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println(curTime)
		writeLog(curTime)
		var allTotal, total = 0, 0
		var msg = ""
		var price float64 = 0
		var minPrice float64 = 999999999

		msg += writeLog(fmt.Sprintf("B链鞋价"))
		total = sneakerTotal(601, 1)
		msg += writeLog(fmt.Sprintf("灰Walker数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 1)
		msg += writeLog(fmt.Sprintf("灰Jogger数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 1)
		msg += writeLog(fmt.Sprintf("灰Runner数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 1)
		msg += writeLog(fmt.Sprintf("灰Trainer数量：%d", total))
		msg += writeLog(fmt.Sprintf("----------"))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(601, 2)
		msg += writeLog(fmt.Sprintf("绿Walker数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 2)
		msg += writeLog(fmt.Sprintf("绿Jogger数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 2)
		msg += writeLog(fmt.Sprintf("绿Runner数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 2)
		msg += writeLog(fmt.Sprintf("绿Trainer数量：%d", total))
		msg += writeLog(fmt.Sprintf("----------"))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(601, 3)
		msg += writeLog(fmt.Sprintf("蓝Walker数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 3)
		msg += writeLog(fmt.Sprintf("蓝Jogger数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 3)
		msg += writeLog(fmt.Sprintf("蓝Runner数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 3)
		msg += writeLog(fmt.Sprintf("蓝Trainer数量：%d", total))
		msg += writeLog(fmt.Sprintf("----------"))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(601, 4)
		msg += writeLog(fmt.Sprintf("紫Walker数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(602, 4)
		msg += writeLog(fmt.Sprintf("紫Jogger数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(603, 4)
		msg += writeLog(fmt.Sprintf("紫Runner数量：%d", total))
		allTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(604, 4)
		msg += writeLog(fmt.Sprintf("紫Trainer数量：%d", total))
		msg += writeLog(fmt.Sprintf("----------"))
		allTotal += total

		msg += writeLog(fmt.Sprintf("总鞋数：%d", allTotal))

		var scrollTotal = 0

		msg += writeLog(fmt.Sprintf(""))
		total = sneakerTotal(701, 1)
		msg += writeLog(fmt.Sprintf("灰卷轴数量：%d", total))
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 2)
		msg += writeLog(fmt.Sprintf("绿卷轴数量：%d", total))
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 3)
		msg += writeLog(fmt.Sprintf("蓝卷轴数量：%d", total))
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 4)
		msg += writeLog(fmt.Sprintf("紫卷轴数量：%d", total))
		scrollTotal += total
		time.Sleep(time.Second * 5)

		total = sneakerTotal(701, 5)
		msg += writeLog(fmt.Sprintf("橙卷轴数量：%d", total))
		msg += writeLog(fmt.Sprintf("----------"))
		scrollTotal += total
		time.Sleep(time.Second * 5)

		msg += writeLog(fmt.Sprintf("总卷轴数：%d", scrollTotal))

		msg += writeLog(fmt.Sprintf(""))
		msg += writeLog(fmt.Sprintf("B链地板价"))

		price = floorPrice(601, 1)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("灰Walker地板价：%.2fbnb", price))
		price = floorPrice(602, 1)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("灰Jogger地板价：%.2fbnb", price))
		price = floorPrice(603, 1)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("灰Runner地板价：%.2fbnb", price))
		price = floorPrice(604, 1)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("灰Trainer地板价：%.2fbnb", price))
		msg += writeLog(fmt.Sprintf("----------"))

		price = floorPrice(601, 2)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("绿Walker地板价：%.2fbnb", price))
		price = floorPrice(602, 2)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("绿Jogger地板价：%.2fbnb", price))
		price = floorPrice(603, 2)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("绿Runner地板价：%.2fbnb", price))
		price = floorPrice(604, 2)
		minPrice = comparePrice(minPrice, price)
		msg += writeLog(fmt.Sprintf("绿Trainer地板价：%.2fbnb", price))
		msg += writeLog(fmt.Sprintf("----------"))

		price = floorPrice(601, 3)
		msg += writeLog(fmt.Sprintf("蓝Walker地板价：%.2fbnb", price))
		price = floorPrice(602, 3)
		msg += writeLog(fmt.Sprintf("蓝Jogger地板价：%.2fbnb", price))
		price = floorPrice(603, 3)
		msg += writeLog(fmt.Sprintf("蓝Runner地板价：%.2fbnb", price))
		price = floorPrice(604, 3)
		msg += writeLog(fmt.Sprintf("蓝Trainer地板价：%.2fbnb", price))
		msg += writeLog(fmt.Sprintf("----------"))

		price = floorPrice(601, 4)
		msg += writeLog(fmt.Sprintf("紫Walker地板价：%.2fbnb", price))
		price = floorPrice(602, 4)
		msg += writeLog(fmt.Sprintf("紫Jogger地板价：%.2fbnb", price))
		price = floorPrice(603, 4)
		msg += writeLog(fmt.Sprintf("紫Runner地板价：%.2fbnb", price))
		price = floorPrice(604, 4)
		msg += writeLog(fmt.Sprintf("紫Trainer地板价：%.2fbnb", price))
		msg += writeLog(fmt.Sprintf("----------"))
		msg += writeLog(fmt.Sprintf("全网地板价：%.2fbnb", minPrice))

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

func floorPrice(types int, quantity int) float64 {

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

	price := float64(orderList.Data[0].SellPrice) / float64(1000000)
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
