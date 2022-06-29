package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

var sneakerNumVars = map[string]string{}
var sneakerMintVars = map[string]string{}
var sneakerFloorVars = map[string]string{}
var scrollVars = map[string]string{}
var sneakerTypeMintNum = map[int]map[int]map[int][]int{}

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

	for {
		HandleSneakerNum()
		log.Fatalln(1)

		newSneakerPrice = map[int]int{}
		genesShoes = []*Shoe{}
		genesis23w = []*Shoe{}

		sneakerNumVars = map[string]string{}
		sneakerMintVars = map[string]string{}
		sneakerFloorVars = map[string]string{}
		scrollVars = map[string]string{}
		sneakerTypeMintNum = map[int]map[int]map[int][]int{}

		curTime := fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println(curTime)

		var _, total = 0, 0
		var msg = ""
		var price float64 = 0
		var minPrice float64 = 999999999
		var rate = ""

		sneakerPriceContent := GetFileContent("sneaker-price.txt")
		_ = json.Unmarshal([]byte(sneakerPriceContent), &sneakerPrice)

		//msg += fmt.Sprintf(`%s\n`, curTime)

		msg += fmt.Sprintf(`\n`)
		if chain == "104" {
			msg += fmt.Sprintf(`💰 鞋子地板价（bnb）\n`)
		} else {
			msg += fmt.Sprintf(`💰 鞋子地板价（sol）\n`)
		}
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

		rate = CalcRate("shoe-floor.txt", fmt.Sprintf("%f", minPrice))
		Insert("shoe-floor.txt", fmt.Sprintf("%f", minPrice))
		msg += fmt.Sprintf(`全网地板 %.2f ｜增幅 %s\n`, minPrice, rate)
		sneakerMinPrice := minPrice

		// 卷轴
		var scrollTotal = 0

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`📜 卷轴数量（市场挂售）\n`)
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
		rate = CalcRate("scroll-total.txt", fmt.Sprintf("%d", scrollTotal))
		Insert("scroll-total.txt", fmt.Sprintf("%d", scrollTotal))
		msg += fmt.Sprintf(`合计 %d｜增幅 %s｜\n`, scrollTotal, rate)

		// 卷轴地板价

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`💰 卷轴地板价（gmt）\n`)
		msg += fmt.Sprintf(`————————————————\n`)

		minPrice = 999999999
		price = floorPrice(701, 1, 100)
		minPrice = comparePrice(minPrice, price)
		scrollMinPrice := minPrice
		msg += fmt.Sprintf(`灰 %.2f｜`, price)

		price = floorPrice(701, 2, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`绿 %.2f｜`, price)

		price = floorPrice(701, 3, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`蓝 %.2f｜`, price)

		price = floorPrice(701, 4, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`紫 %.2f｜`, price)

		price = floorPrice(701, 5, 100)
		minPrice = comparePrice(minPrice, price)
		msg += fmt.Sprintf(`橙 %.2f｜\n`, price)

		rate = CalcRate("scroll-floor.txt", fmt.Sprintf("%f", minPrice))
		Insert("scroll-floor.txt", fmt.Sprintf("%f", minPrice))

		msg += fmt.Sprintf(`全网地板 %.2f｜增幅 %s\n`, minPrice, rate)
		msg += fmt.Sprintf(`\n`)

		if chain == "104" {
			msg += fmt.Sprintf(`💰 Mint利润（usd）\n`)
			msg += fmt.Sprintf(`————————————————\n`)
			gstPrice, gmtPrice, profit := CalcMintProfitForBSC(sneakerMinPrice, scrollMinPrice)
			msg += fmt.Sprintf(`1GST = %.4fU \n`, gstPrice)
			msg += fmt.Sprintf(`1GMT = %.4fU \n`, gmtPrice)
			msg += fmt.Sprintf(`mint费用 = %.4fU \n`, 360*gstPrice+40*gmtPrice)
			msg += fmt.Sprintf(`卷轴费用 = %.4fU \n`, scrollMinPrice*gmtPrice*2)
			msg += fmt.Sprintf(`升级费用 = %.4fU \n`, 20*gstPrice+10*gmtPrice)
			msg += fmt.Sprintf(`%s\n`, profit)
		} else {
			msg += fmt.Sprintf(`💰 Mint利润（usd）\n`)
			msg += fmt.Sprintf(`————————————————\n`)
			gstPrice, gmtPrice, profit := CalcMintProfitForSol(sneakerMinPrice, scrollMinPrice)
			msg += fmt.Sprintf(`1GST = %.4fU \n`, gstPrice)
			msg += fmt.Sprintf(`1GMT = %.4fU \n`, gmtPrice)
			msg += fmt.Sprintf(`mint费用 = %.4fU \n`, 360*gstPrice+40*gmtPrice)
			msg += fmt.Sprintf(`卷轴费用 = %.4fU \n`, scrollMinPrice*gmtPrice*2)
			msg += fmt.Sprintf(`升级费用 = %.4fU \n`, 20*gstPrice+10*gmtPrice)
			msg += fmt.Sprintf(`%s\n`, profit)
		}

		msg += fmt.Sprintf(`\n`)
		msg += fmt.Sprintf(`PS：数据存在误差，仅供参考，非投资建议 \n`)

		push(msg)

		// 给老的存起来，新的清空
		newSneakerPriceByte, _ := json.Marshal(newSneakerPrice)
		Rewrite("sneaker-price.txt", string(newSneakerPriceByte))
		newSneakerPrice = map[int]int{}
		sneakerPrice = map[int]int{}

		// 推创世
		GenesShoes()

		// 推OG
		if chain == "104" {
			Genesis23wShoes()
		}

		time.Sleep(time.Second * 300)
	}
}

func HandleSneakerNum() {

	var allTotal = 0
	var rate = ""

	if chain == "103" {
		sneakerNumVars["chain_name"] = "SOL"
	} else if chain == "104" {
		sneakerNumVars["chain_name"] = "BSC"
	} else {
		sneakerNumVars["chain_name"] = "ETH"
	}
	sneakerNumVars["time"] = fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))

	allTotal += AutoSetSneakerVar(601, 1, "common_w")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(602, 1, "common_j")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(603, 1, "common_r")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(604, 1, "common_t")
	time.Sleep(time.Second * 5)

	allTotal += AutoSetSneakerVar(601, 2, "uncommon_w")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(602, 2, "uncommon_j")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(603, 2, "uncommon_r")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(604, 2, "uncommon_t")
	time.Sleep(time.Second * 5)

	allTotal += AutoSetSneakerVar(601, 3, "rare_w")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(602, 3, "rare_j")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(603, 3, "rare_r")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(604, 3, "rare_t")
	time.Sleep(time.Second * 5)

	allTotal += AutoSetSneakerVar(601, 4, "epic_w")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(602, 4, "epic_j")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(603, 4, "epic_r")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(604, 4, "epic_t")
	time.Sleep(time.Second * 5)

	allTotal += AutoSetSneakerVar(601, 5, "legendary_w")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(602, 5, "legendary_j")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(603, 5, "legendary_r")
	time.Sleep(time.Second * 5)
	allTotal += AutoSetSneakerVar(604, 5, "legendary_t")
	time.Sleep(time.Second * 5)

	rate = CalcRate("shoe-total.txt", fmt.Sprintf("%d", allTotal))
	Insert("shoe-total.txt", fmt.Sprintf("%d", allTotal))
	sneakerNumVars["rate"] = fmt.Sprintf("%s", rate)

	newNum, oldNum, avgPrice, middlePrice := CalcDiffNumSneakers(sneakerPrice, newSneakerPrice)

	if newNum > oldNum {
		sneakerNumVars["notice"] = "市场通胀中……"
	} else {
		sneakerNumVars["notice"] = "市场通缩中……"
	}
	sneakerNumVars["total"] = fmt.Sprintf("%d", allTotal)
	sneakerNumVars["new_total"] = fmt.Sprintf("%d", newNum)
	sneakerNumVars["consume_total"] = fmt.Sprintf("%d", oldNum)
	sneakerNumVars["new_avg_price"] = fmt.Sprintf("%s", avgPrice)
	sneakerNumVars["new_middle_price"] = fmt.Sprintf("%s", middlePrice)

	if chain == "104" {
		p, num := NumBelowTo(newSneakerPrice)
		sneakerNumVars["below_price_1"] = fmt.Sprintf("%sBNB", p)
		sneakerNumVars["below_price_total_1"] = fmt.Sprintf("%d", num)
		p, num = NumBelowToNext(newSneakerPrice)
		sneakerNumVars["below_price_2"] = fmt.Sprintf("%sBNB", p)
		sneakerNumVars["below_price_total_2"] = fmt.Sprintf("%d", num)
	} else if chain == "103" {
		p, num := NumBelowTo(newSneakerPrice)
		sneakerNumVars["below_price_1"] = fmt.Sprintf("%sSOL", p)
		sneakerNumVars["below_price_total_1"] = fmt.Sprintf("%d", num)
		p, num = NumBelowToNext(newSneakerPrice)
		sneakerNumVars["below_price_2"] = fmt.Sprintf("%sSOL", p)
		sneakerNumVars["below_price_total_2"] = fmt.Sprintf("%d", num)
	} else {

	}

	template := "templates/sneaker-num.html"
	newFile := fmt.Sprintf("o-%s-sneaker-num.html", chain)
	newImage := fmt.Sprintf("%s-sneaker-num.jpg", chain)
	ReplaceVar(template, sneakerNumVars, newFile)
	Html2Image(newFile, newImage)
	PushFile(newImage)
}

func AutoSetSneakerVar(types int, quality int, varName string) int {
	total := sneakerTotal(types, quality)
	rate := CalcRate(fmt.Sprintf("sneaker-%d-%d.txt", types, quality), fmt.Sprintf("%d", total))
	Insert(fmt.Sprintf("sneaker-%d-%d.txt", types, quality), fmt.Sprintf("%d", total))
	if strings.Contains(rate, "-") {
		sneakerNumVars[varName] = fmt.Sprintf(`%s<li style="border-right: none; width: auto;">&nbsp;&nbsp;<label style="font-weight: bold;">总数：%d</label> &nbsp;&nbsp;<label style="color: red;">增幅%s</label></li>`, CalcMintNum(types, quality), total, rate)
	} else {
		sneakerNumVars[varName] = fmt.Sprintf(`%s<li style="border-right: none; width: auto;">&nbsp;&nbsp;<label style="font-weight: bold;">总数：%d</label> &nbsp;&nbsp;<label style="color: green;">增幅%s</label></li>`, CalcMintNum(types, quality), total, rate)
	}

	return total
}

func CalcMintNum(types int, quality int) string {
	output := ""
	for i := 0; i < 4; i++ {
		otds := RemoveDuplicateElement(sneakerTypeMintNum[types][quality][i])
		output += fmt.Sprintf(`<li>M%d:%d</li>`, i, len(otds))
	}
	total4to8 := 0
	for i := 4; i < 8; i++ {
		otds := RemoveDuplicateElement(sneakerTypeMintNum[types][quality][i])
		total4to8 += len(otds)
	}
	output += fmt.Sprintf(`<li style="width: 100px;">M4～7:%d</li>`, total4to8)
	fmt.Println(sneakerTypeMintNum)
	return output
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

			_, ok := sneakerTypeMintNum[types]
			if ok {
				_, ok := sneakerTypeMintNum[types][data.Quantity]
				if ok {
					_, ok := sneakerTypeMintNum[types][data.Quantity][data.Mint]
					if ok {
						sneakerTypeMintNum[types][data.Quantity][data.Mint] = append(sneakerTypeMintNum[types][data.Quantity][data.Mint], data.Otd)
					} else {
						sneakerTypeMintNum[types][data.Quantity][data.Mint] = []int{data.Otd}
					}
				} else {
					sneakerTypeMintNum[types][data.Quantity] = map[int][]int{
						data.Mint: {data.Otd},
					}
				}
			} else {
				sneakerTypeMintNum[types] = map[int]map[int][]int{
					data.Quantity: {
						data.Mint: {data.Otd},
					},
				}
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

			// 重复鞋子出现5次就退出
			_, ok := itemStatic[data.Otd]
			if ok {
				repeatCount++
			} else {
				itemStatic[data.Otd] = 1
			}
			if repeatCount >= 5 {
				break
			}

			_, ok = sneakerTypeMintNum[types]
			if ok {
				_, ok := sneakerTypeMintNum[types][data.Quantity]
				if ok {
					_, ok := sneakerTypeMintNum[types][data.Quantity][data.Mint]
					if ok {
						sneakerTypeMintNum[types][data.Quantity][data.Mint] = append(sneakerTypeMintNum[types][data.Quantity][data.Mint], data.Otd)
					} else {
						sneakerTypeMintNum[types][data.Quantity][data.Mint] = []int{data.Otd}
					}
				} else {
					sneakerTypeMintNum[types][data.Quantity] = map[int][]int{
						data.Mint: {data.Otd},
					}
				}
			} else {
				sneakerTypeMintNum[types] = map[int]map[int][]int{
					data.Quantity: {
						data.Mint: {data.Otd},
					},
				}
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
		fmt.Println("文件打开失败", err)
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

func PushFile(filePath string) {

	url := "https://discord.com/api/webhooks/990263471415386175/R7AWHDh_N-fOKdL0Tt9CuhDDLNA1uu_Mr2CKLwtEQiQ7QqLJcXg_fF5CTqdLIRI1Brhg"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(filePath)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("file1", filepath.Base(filePath))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "apifox/1.0.0 (https://www.apifox.cn)")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func pushToGenes(msg string) {
	pushDcFromConfigKey("genes_webhook", msg)
}

func pushToGenesis23w(msg string) {
	pushDcFromConfigKey("genesis23w_webhook", msg)
}
