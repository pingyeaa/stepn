package main

import (
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
	"sort"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	_ "github.com/go-sql-driver/mysql"
)

var err error
var cookie = ""
var cfg *ini.File
var sneakerPrice = map[int]int{}
var newSneakerPrice = map[int]int{}
var itemStatic = map[int]int{}
var chain = "104"
var genesShoes []*Shoe
var genesis23w []*Shoe
var sneakerMinPrice float64 = 0
var scrollMinPrice float64 = 0

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

		newSneakerPrice = map[int]int{}
		genesShoes = []*Shoe{}
		genesis23w = []*Shoe{}

		sneakerNumVars = map[string]string{}
		sneakerMintVars = map[string]string{}
		sneakerFloorVars = map[string]string{}
		scrollVars = map[string]string{}
		sneakerTypeMintNum = map[int]map[int]map[int][]int{}

		sneakerPriceContent := GetFileContent("sneaker-price.txt")
		_ = json.Unmarshal([]byte(sneakerPriceContent), &sneakerPrice)

		HandleSneakerNum()
		HandleSneakerFloor()
		HandleScroll()
		HandleMint()
		HandleGem()

		// 给老的存起来，新的清空
		newSneakerPriceByte, _ := json.Marshal(newSneakerPrice)
		Rewrite("sneaker-price.txt", string(newSneakerPriceByte))
		newSneakerPrice = map[int]int{}
		sneakerPrice = map[int]int{}

		// 推创世
		HandleGenesis()

		// 推OG
		if chain == "104" {
			HandleOG()
		}

		time.Sleep(time.Second * 300)
	}
}

func HandleGenesis() {

	var ogVars = map[string]string{}
	var genesOtd []int
	var handled = map[int]*Shoe{}
	var uniqueCommon = map[int]*Shoe{}
	var uniqueUncommon = map[int]*Shoe{}
	var uniqueRare = map[int]*Shoe{}
	var uniqueEpic = map[int]*Shoe{}
	var uniqueLegendary = map[int]*Shoe{}
	var minPrice = 9999999999
	var count = 0
	var unitName = ""

	if chain == "104" {
		ogVars["chain_name"] = "BSC"
		unitName = "BNB"
	} else if chain == "103" {
		ogVars["chain_name"] = "SOL"
		unitName = "SOL"
	} else {
		ogVars["chain_name"] = "ETH"
		unitName = "ETH"
	}
	ogVars["time"] = fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))

	// 创世鞋排序
	for _, shoe := range genesShoes {
		genesOtd = append(genesOtd, shoe.Otd)
	}
	genesOtd = RemoveDuplicateElement(genesOtd)
	sort.Ints(genesOtd)

	// OG去重复
	for _, otd := range genesOtd {
		for _, shoe := range genesShoes {
			if otd == shoe.Otd {
				_, ok := handled[shoe.Otd]
				if ok {
					continue
				}
				handled[shoe.Otd] = shoe
				qualityMapper := map[int]string{1: "灰", 2: "绿", 3: "蓝", 4: "紫", 5: "橙"}
				typesMapper := map[int]string{601: "W", 602: "J", 603: "R", 604: "T"}
				shoe.Color = qualityMapper[shoe.Quantity]
				shoe.TypeName = typesMapper[shoe.TypeID]
				if minPrice > shoe.SellPrice {
					minPrice = shoe.SellPrice
				}
				if shoe.Quantity == 1 {
					uniqueCommon[shoe.Otd] = shoe
				}
				if shoe.Quantity == 2 {
					uniqueUncommon[shoe.Otd] = shoe
				}
				if shoe.Quantity == 3 {
					uniqueRare[shoe.Otd] = shoe
				}
				if shoe.Quantity == 4 {
					uniqueEpic[shoe.Otd] = shoe
				}
				if shoe.Quantity == 5 {
					uniqueLegendary[shoe.Otd] = shoe
				}
			}
		}
	}

	// common html
	var commonHtml = ``
	if len(uniqueCommon) == 0 {
		commonHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueCommon {
		style := `border-top: none;`
		if count == len(uniqueCommon)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		commonHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2f%s</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000, unitName)
	}
	ogVars["common"] = commonHtml

	// uncommon html
	var uncommonHtml = ``
	if len(uniqueUncommon) == 0 {
		uncommonHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueUncommon {
		style := `border-top: none;`
		if count == len(uniqueUncommon)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		uncommonHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2f%s</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000, unitName)
	}
	ogVars["uncommon"] = uncommonHtml

	// rare html
	var rareHtml = ``
	if len(uniqueRare) == 0 {
		rareHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueRare {
		style := `border-top: none;`
		if count == len(uniqueRare)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		rareHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2f%s</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000, unitName)
	}
	ogVars["rare"] = rareHtml

	// epic html
	var epicHtml = ``
	if len(uniqueEpic) == 0 {
		epicHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueEpic {
		style := `border-top: none;`
		if count == len(uniqueEpic)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		epicHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2f%s</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000, unitName)
	}
	ogVars["epic"] = epicHtml

	// legendary html
	var legendaryHtml = ``
	if len(uniqueLegendary) == 0 {
		legendaryHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueLegendary {
		style := `border-top: none;`
		if count == len(uniqueLegendary)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		legendaryHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2f%s</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000, unitName)
	}
	ogVars["legendary"] = legendaryHtml

	ogVars["total"] = fmt.Sprintf("%d", len(genesOtd))
	ogVars["floor"] = fmt.Sprintf("%.2f%s", float64(minPrice)/1000000, unitName)

	// export image & push to discord
	template := "templates/sneaker-genesis.html"
	newFile := fmt.Sprintf("o-%s-sneaker-genesis.html", chain)
	newImage := chain + "-sneaker_genesis.jpg"
	ReplaceVar(template, ogVars, newFile)
	Html2Image(newFile, newImage)
	webhook, err := cfg.Section("discord").GetKey("genes_webhook")
	if err != nil {
		log.Fatalln(err.Error())
	}
	PushFile(newImage, webhook.String())
}

func HandleOG() {

	var ogVars = map[string]string{}
	var genesOtd []int
	var handled = map[int]*Shoe{}
	var uniqueCommon = map[int]*Shoe{}
	var uniqueUncommon = map[int]*Shoe{}
	var uniqueRare = map[int]*Shoe{}
	var uniqueEpic = map[int]*Shoe{}
	var uniqueLegendary = map[int]*Shoe{}
	var minPrice = 9999999999
	var count = 0

	if chain != "104" {
		return
	}
	ogVars["chain_name"] = "BSC"
	ogVars["time"] = fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))

	// OG鞋排序
	for _, shoe := range genesis23w {
		genesOtd = append(genesOtd, shoe.Otd)
	}
	genesOtd = RemoveDuplicateElement(genesOtd)
	sort.Ints(genesOtd)

	// OG去重复
	for _, otd := range genesOtd {
		for _, shoe := range genesis23w {
			if otd == shoe.Otd {
				_, ok := handled[shoe.Otd]
				if ok {
					continue
				}
				handled[shoe.Otd] = shoe
				qualityMapper := map[int]string{1: "灰", 2: "绿", 3: "蓝", 4: "紫", 5: "橙"}
				typesMapper := map[int]string{601: "W", 602: "J", 603: "R", 604: "T"}
				shoe.Color = qualityMapper[shoe.Quantity]
				shoe.TypeName = typesMapper[shoe.TypeID]
				if minPrice > shoe.SellPrice {
					minPrice = shoe.SellPrice
				}
				if shoe.Quantity == 1 {
					uniqueCommon[shoe.Otd] = shoe
				}
				if shoe.Quantity == 2 {
					uniqueUncommon[shoe.Otd] = shoe
				}
				if shoe.Quantity == 3 {
					uniqueRare[shoe.Otd] = shoe
				}
				if shoe.Quantity == 4 {
					uniqueEpic[shoe.Otd] = shoe
				}
				if shoe.Quantity == 5 {
					uniqueLegendary[shoe.Otd] = shoe
				}
			}
		}
	}

	// common html
	var commonHtml = ``
	if len(uniqueCommon) == 0 {
		commonHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueCommon {
		style := `border-top: none;`
		if count == len(uniqueCommon)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		commonHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2fBNB</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000)
	}
	ogVars["common"] = commonHtml

	// uncommon html
	var uncommonHtml = ``
	if len(uniqueUncommon) == 0 {
		uncommonHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueUncommon {
		style := `border-top: none;`
		if count == len(uniqueUncommon)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		uncommonHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2fBNB</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000)
	}
	ogVars["uncommon"] = uncommonHtml

	// rare html
	var rareHtml = ``
	if len(uniqueRare) == 0 {
		rareHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueRare {
		style := `border-top: none;`
		if count == len(uniqueRare)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		rareHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2fBNB</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000)
	}
	ogVars["rare"] = rareHtml

	// epic html
	var epicHtml = ``
	if len(uniqueEpic) == 0 {
		epicHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueEpic {
		style := `border-top: none;`
		if count == len(uniqueEpic)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		epicHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2fBNB</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000)
	}
	ogVars["epic"] = epicHtml

	// legendary html
	var legendaryHtml = ``
	if len(uniqueLegendary) == 0 {
		legendaryHtml = "   暂无数据"
	}
	count = 0
	for _, shoe := range uniqueLegendary {
		style := `border-top: none;`
		if count == len(uniqueLegendary)-1 {
			style = `border-top: none; border-bottom: none;`
		}
		count++
		legendaryHtml += fmt.Sprintf(`
		<div class="fangkuai" style="%s">
			<ul>
				<li>#%d</li>
				<li>%s-%s</li>
				<li>Lv%d</li>
				<li>Mint%d</li>
				<li style="border-right: none; width: 150px;">  %.2fBNB</li>
			</ul>
		</div>`, style, shoe.Otd, shoe.Color, shoe.TypeName, shoe.Level, shoe.Mint, float64(shoe.SellPrice)/1000000)
	}
	ogVars["legendary"] = legendaryHtml

	ogVars["total"] = fmt.Sprintf("%d", len(genesOtd))
	ogVars["floor"] = fmt.Sprintf("%.2fBNB", float64(minPrice)/1000000)

	// export image & push to discord
	template := "templates/sneaker-og.html"
	newFile := fmt.Sprintf("o-%s-sneaker-og.html", chain)
	newImage := chain + "-sneaker_og.jpg"
	ReplaceVar(template, ogVars, newFile)
	Html2Image(newFile, newImage)
	webhook, err := cfg.Section("discord").GetKey("genesis23w_webhook")
	if err != nil {
		log.Fatalln(err.Error())
	}
	PushFile(newImage, webhook.String())
}

func HandleMint() {

	if chain == "103" {
		sneakerMintVars["chain_name"] = "SOL"
	} else if chain == "104" {
		sneakerMintVars["chain_name"] = "BSC"
	} else {
		sneakerMintVars["chain_name"] = "ETH"
	}
	sneakerMintVars["time"] = fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))

	if chain == "104" {
		gstPrice, gmtPrice, profit := CalcMintProfitForBSC(sneakerMinPrice, scrollMinPrice)
		sneakerMintVars["gst_price_u"] = fmt.Sprintf("%.4fU", gstPrice)
		sneakerMintVars["gmt_price_u"] = fmt.Sprintf("%.4fU", gmtPrice)
		sneakerMintVars["scroll_price_u"] = fmt.Sprintf("%.4fU", scrollMinPrice*gmtPrice*2)
		sneakerMintVars["mint_price_u"] = fmt.Sprintf("%.4fU", 360*gstPrice+40*gmtPrice)
		sneakerMintVars["upgrade_price_u"] = fmt.Sprintf("%.4fU", 20*gstPrice+10*gmtPrice)
		sneakerMintVars["sneaker_floor_price_u"] = fmt.Sprintf("%.4fBNB", sneakerMinPrice)
		sneakerMintVars["formula"] = profit
	} else {
		gstPrice, gmtPrice, profit := CalcMintProfitForSol(sneakerMinPrice, scrollMinPrice)
		sneakerMintVars["gst_price_u"] = fmt.Sprintf("%.4fU", gstPrice)
		sneakerMintVars["gmt_price_u"] = fmt.Sprintf("%.4fU", gmtPrice)
		sneakerMintVars["scroll_price_u"] = fmt.Sprintf("%.4fU", scrollMinPrice*gmtPrice*2)
		sneakerMintVars["mint_price_u"] = fmt.Sprintf("%.4fU", 360*gstPrice+40*gmtPrice)
		sneakerMintVars["upgrade_price_u"] = fmt.Sprintf("%.4fU", 20*gstPrice+10*gmtPrice)
		sneakerMintVars["sneaker_floor_price_u"] = fmt.Sprintf("%.4fSol", sneakerMinPrice)
		sneakerMintVars["formula"] = profit
	}

	template := "templates/sneaker-mint.html"
	newFile := fmt.Sprintf("o-%s-sneaker-mint.html", chain)
	newImage := chain + "-sneaker_mint.jpg"
	ReplaceVar(template, sneakerMintVars, newFile)
	Html2Image(newFile, newImage)
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
	PushFile(newImage, webhook.String())
}

func HandleGem() {

	var minPrice float64 = 999999999
	var price float64 = 0
	var allTotal = 0
	var total = 0
	var vars = map[string]string{}
	var unitName = ""
	gemHandled = map[int]*Shoe{}

	if chain == "103" {
		vars["chain_name"] = "SOL"
		unitName = "SOL"
	} else if chain == "104" {
		vars["chain_name"] = "BSC"
		unitName = "BNB"
	} else {
		vars["chain_name"] = "ETH"
		unitName = "ETH"
	}
	vars["time"] = fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))

	total = GemTotal(1)
	vars["e_total"] = fmt.Sprintf("%d", total)
	allTotal += total
	time.Sleep(time.Second * 5)

	total = GemTotal(2)
	vars["l_total"] = fmt.Sprintf("%d", total)
	allTotal += total
	time.Sleep(time.Second * 5)

	total = GemTotal(3)
	vars["c_total"] = fmt.Sprintf("%d", total)
	allTotal += total
	time.Sleep(time.Second * 4)

	total = GemTotal(2)
	vars["r_total"] = fmt.Sprintf("%d", total)
	allTotal += total
	time.Sleep(time.Second * 5)

	rate := CalcRate("gem-total.txt", fmt.Sprintf("%d", allTotal))
	Insert("gem-total.txt", fmt.Sprintf("%d", allTotal))
	vars["total"] = fmt.Sprintf("%d", allTotal)
	vars["rate"] = SetFloatColor(rate)

	// 宝石地板价
	minPrice = 999999999
	price = GemFloorPrice(1)
	vars["e_floor"] = fmt.Sprintf("%.2f%s", price, unitName)
	minPrice = comparePrice(minPrice, price)
	scrollMinPrice = minPrice

	price = GemFloorPrice(2)
	vars["l_floor"] = fmt.Sprintf("%.2f%s", price, unitName)
	minPrice = comparePrice(minPrice, price)

	price = GemFloorPrice(3)
	vars["c_floor"] = fmt.Sprintf("%.2f%s", price, unitName)
	minPrice = comparePrice(minPrice, price)

	price = GemFloorPrice(4)
	vars["r_floor"] = fmt.Sprintf("%.2f%s", price, unitName)
	minPrice = comparePrice(minPrice, price)

	rate = CalcRate("gem-floor.txt", fmt.Sprintf("%f", minPrice))
	Insert("gem-floor.txt", fmt.Sprintf("%f", minPrice))
	vars["floor_price"] = fmt.Sprintf("%.2f%s", minPrice, unitName)
	vars["floor_rate"] = SetFloatColor(rate)

	var gemGTypeLevelMapper = map[int]map[int]int{}
	for _, gem := range gemHandled {
		_, ok := gemGTypeLevelMapper[gem.GType]
		if ok {
			_, ok := gemGTypeLevelMapper[gem.GType][gem.Quantity]
			if ok {
				gemGTypeLevelMapper[gem.GType][gem.Quantity]++
			} else {
				gemGTypeLevelMapper[gem.GType][gem.Quantity] = 1
			}
		} else {
			gemGTypeLevelMapper[gem.GType] = map[int]int{
				gem.Quantity: 1,
			}
		}
	}

	var eHtml = ``
	for i := 1; i < 10; i++ {
		number := 0
		_, ok := gemGTypeLevelMapper[1]
		if ok {
			_, ok := gemGTypeLevelMapper[1][i]
			if ok {
				number = gemGTypeLevelMapper[1][i]
				eHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, number)
			} else {
				eHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
			}
		} else {
			eHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
		}
	}
	vars["efficiency"] = eHtml

	var lHtml = ``
	for i := 1; i < 10; i++ {
		number := 0
		_, ok := gemGTypeLevelMapper[2]
		if ok {
			_, ok := gemGTypeLevelMapper[2][i]
			if ok {
				number = gemGTypeLevelMapper[2][i]
				lHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, number)
			} else {
				lHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
			}
		} else {
			lHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
		}
	}
	vars["luck"] = lHtml

	var cHtml = ``
	for i := 1; i < 10; i++ {
		number := 0
		_, ok := gemGTypeLevelMapper[3]
		if ok {
			_, ok := gemGTypeLevelMapper[3][i]
			if ok {
				number = gemGTypeLevelMapper[3][i]
				cHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, number)
			} else {
				cHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
			}
		} else {
			cHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
		}
	}
	vars["comfort"] = cHtml

	var rHtml = ``
	for i := 1; i < 10; i++ {
		number := 0
		_, ok := gemGTypeLevelMapper[4]
		if ok {
			_, ok := gemGTypeLevelMapper[4][i]
			if ok {
				number = gemGTypeLevelMapper[4][i]
				rHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, number)
			} else {
				rHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
			}
		} else {
			rHtml += fmt.Sprintf(`<li>L%d:%d</li>`, i, 0)
		}
	}
	vars["resilience"] = rHtml

	belowPrice, belowTotals := GetPriceBelowNextPrice(gemHandled)
	vars["below_price"] = fmt.Sprintf("%.2f%s", belowPrice, unitName)
	vars["below_total_e"] = fmt.Sprintf("%d", belowTotals[1])
	vars["below_total_l"] = fmt.Sprintf("%d", belowTotals[2])
	vars["below_total_c"] = fmt.Sprintf("%d", belowTotals[3])
	vars["below_total_r"] = fmt.Sprintf("%d", belowTotals[4])

	template := "templates/gem.html"
	newFile := fmt.Sprintf("o-%s-gem.html", chain)
	newImage := chain + "-gem.jpg"
	ReplaceVar(template, vars, newFile)
	Html2Image(newFile, newImage)
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
	PushFile(newImage, webhook.String())
}

func HandleScroll() {

	var minPrice float64 = 999999999
	var price float64 = 0
	var scrollTotal = 0
	var total = 0

	if chain == "103" {
		scrollVars["chain_name"] = "SOL"
	} else if chain == "104" {
		scrollVars["chain_name"] = "BSC"
	} else {
		scrollVars["chain_name"] = "ETH"
	}
	scrollVars["time"] = fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))

	total = sneakerTotal(701, 1)
	scrollVars["common_total"] = fmt.Sprintf("%d", total)
	scrollTotal += total
	time.Sleep(time.Second * 5)

	total = sneakerTotal(701, 2)
	scrollVars["uncommon_total"] = fmt.Sprintf("%d", total)
	scrollTotal += total
	time.Sleep(time.Second * 5)

	total = sneakerTotal(701, 3)
	scrollVars["rare_total"] = fmt.Sprintf("%d", total)
	scrollTotal += total
	time.Sleep(time.Second * 5)

	total = sneakerTotal(701, 4)
	scrollVars["epic_total"] = fmt.Sprintf("%d", total)
	scrollTotal += total
	time.Sleep(time.Second * 5)

	total = sneakerTotal(701, 5)
	scrollVars["legendary_total"] = fmt.Sprintf("%d", total)
	scrollTotal += total
	time.Sleep(time.Second * 5)

	rate := CalcRate("scroll-total.txt", fmt.Sprintf("%d", scrollTotal))
	Insert("scroll-total.txt", fmt.Sprintf("%d", scrollTotal))
	scrollVars["total"] = fmt.Sprintf("%d", scrollTotal)
	scrollVars["rate"] = SetFloatColor(rate)

	// 卷轴地板价
	minPrice = 999999999
	price = floorPrice(701, 1, 100)
	scrollVars["common_price"] = fmt.Sprintf("%.2fGMT", price)
	minPrice = comparePrice(minPrice, price)
	scrollMinPrice = minPrice

	price = floorPrice(701, 2, 100)
	scrollVars["uncommon_price"] = fmt.Sprintf("%.2fGMT", price)
	minPrice = comparePrice(minPrice, price)

	price = floorPrice(701, 3, 100)
	scrollVars["rare_price"] = fmt.Sprintf("%.2fGMT", price)
	minPrice = comparePrice(minPrice, price)

	price = floorPrice(701, 4, 100)
	scrollVars["epic_price"] = fmt.Sprintf("%.2fGMT", price)
	minPrice = comparePrice(minPrice, price)

	price = floorPrice(701, 5, 100)
	scrollVars["legendary_price"] = fmt.Sprintf("%.2fGMT", price)
	minPrice = comparePrice(minPrice, price)

	rate = CalcRate("scroll-floor.txt", fmt.Sprintf("%f", minPrice))
	Insert("scroll-floor.txt", fmt.Sprintf("%f", minPrice))
	scrollVars["floor_price"] = fmt.Sprintf("%.2f", minPrice)
	scrollVars["floor_rate"] = SetFloatColor(rate)

	template := "templates/scroll.html"
	newFile := fmt.Sprintf("o-%s-scroll.html", chain)
	newImage := chain + "-scroll.jpg"
	ReplaceVar(template, scrollVars, newFile)
	Html2Image(newFile, newImage)
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
	PushFile(newImage, webhook.String())
}

func HandleSneakerFloor() {

	var minPrice float64 = 999999999
	var price float64 = 0
	var unitName = "BNB"

	if chain == "103" {
		unitName = "SOL"
		sneakerFloorVars["chain_name"] = "SOL"
	} else if chain == "104" {
		unitName = "BNB"
		sneakerFloorVars["chain_name"] = "BSC"
	} else {
		unitName = "ETH"
		sneakerFloorVars["chain_name"] = "ETH"
	}
	sneakerFloorVars["time"] = fmt.Sprintf(`%s`, time.Now().Format("2006-01-02 15:04:05"))

	price = floorPrice(601, 1, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["common_w_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(602, 1, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["common_j_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(603, 1, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["common_r_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(604, 1, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["common_t_floor"] = fmt.Sprintf("%.2f %s", price, unitName)

	price = floorPrice(601, 2, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["uncommon_w_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(602, 2, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["uncommon_j_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(603, 2, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["uncommon_r_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(604, 2, 1000000)
	minPrice = comparePrice(minPrice, price)
	sneakerFloorVars["uncommon_t_floor"] = fmt.Sprintf("%.2f %s", price, unitName)

	sneakerMinPrice = minPrice

	price = floorPrice(601, 3, 1000000)
	sneakerFloorVars["rare_w_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(602, 3, 1000000)
	sneakerFloorVars["rare_j_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(603, 3, 1000000)
	sneakerFloorVars["rare_r_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(604, 3, 1000000)
	sneakerFloorVars["rare_t_floor"] = fmt.Sprintf("%.2f %s", price, unitName)

	price = floorPrice(601, 4, 1000000)
	sneakerFloorVars["epic_w_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(602, 4, 1000000)
	sneakerFloorVars["epic_j_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(603, 4, 1000000)
	sneakerFloorVars["epic_r_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(604, 4, 1000000)
	sneakerFloorVars["epic_t_floor"] = fmt.Sprintf("%.2f %s", price, unitName)

	price = floorPrice(601, 5, 1000000)
	sneakerFloorVars["legendary_w_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(602, 5, 1000000)
	sneakerFloorVars["legendary_j_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(603, 5, 1000000)
	sneakerFloorVars["legendary_r_floor"] = fmt.Sprintf("%.2f %s", price, unitName)
	price = floorPrice(604, 5, 1000000)
	sneakerFloorVars["legendary_t_floor"] = fmt.Sprintf("%.2f %s", price, unitName)

	rate := CalcRate("shoe-floor.txt", fmt.Sprintf("%f", minPrice))
	Insert("shoe-floor.txt", fmt.Sprintf("%f", minPrice))

	sneakerFloorVars["floor"] = fmt.Sprintf("%.2f%s", minPrice, unitName)
	sneakerFloorVars["rate"] = SetFloatColor(rate)

	template := "templates/sneaker-floor.html"
	newFile := fmt.Sprintf("o-%s-sneaker-floor.html", chain)
	newImage := chain + "-sneaker-floor.jpg"
	ReplaceVar(template, sneakerFloorVars, newFile)
	Html2Image(newFile, newImage)

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
	PushFile(newImage, webhook.String())
}

func SetFloatColor(rate string) string {
	var output = ""
	if strings.Contains(rate, "-") {
		output = `<label style="color:red;">` + fmt.Sprintf(`增幅 %s`, rate) + `</label>`
	} else {
		output = `<label style="color:green;">` + fmt.Sprintf(`增幅 %s`, rate) + `</label>`
	}
	return output
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
	sneakerNumVars["rate"] = SetFloatColor(rate)

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
	fmt.Println(webhook)
	PushFile(newImage, webhook.String())
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
	return output
}

func sneakerTotal(types int, quantity int) int {

	var page = 0
	var url = ""

	// https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=1&chain=104&page=0&refresh=true
	// https://apilb.stepn.com/run/orderlist?order=2001&type=501&gType=2&chain=104&page=0&refresh=true

	for {
		if page == 0 {
			if types == 501 {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&gType=%d&chain=%s&page=%d&refresh=true", types, quantity, chain, page)
			} else {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&quality=%d&chain=%s&page=%d&refresh=true", types, quantity, chain, page)
			}
		} else {
			if types == 501 {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&gType=%d&chain=%s&page=%d&refresh=false", types, quantity, chain, page)
			} else {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2001&type=%d&quality=%d&chain=%s&page=%d&refresh=false", types, quantity, chain, page)
			}
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

			if types != 701 && types != 501 {
				newSneakerPrice[data.Otd] = data.SellPrice
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
			if types == 501 {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=%d&gType=%d&chain=%s&page=%d&refresh=true", types, quantity, chain, page)
			} else {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=%d&quality=%d&chain=%s&page=%d&refresh=true", types, quantity, chain, page)
			}
		} else {
			if types == 501 {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=%d&gType=%d&chain=%s&page=%d&refresh=false", types, quantity, chain, page)
			} else {
				url = fmt.Sprintf("https://apilb.stepn.com/run/orderlist?order=2002&type=%d&quality=%d&chain=%s&page=%d&refresh=false", types, quantity, chain, page)
			}
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

			if types != 701 && types != 501 {
				newSneakerPrice[data.Otd] = data.SellPrice
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

func PushFile(filePath string, webhook string) {

	url := webhook
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
