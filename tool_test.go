package main

import (
	"fmt"
	"log"
	"testing"
)

func TestFindLatest(t *testing.T) {
	res := FindLatest("a.txt")
	fmt.Println(res)
}

func TestInsert(t *testing.T) {
	Insert("a.txt", "20")
}

func TestCalcRate(t *testing.T) {
	res := CalcRate("a.txt", "99")
	fmt.Println(res)
}

func TestCalcDiffNumSneakers(t *testing.T) {
	a := map[int]int{
		5: 2,
		1: 2,
		2: 3,
		3: 4,
	}
	b := map[int]int{
		1: 2,
		//4: 6,
		//8: 3,
		//9: 4,
		//19: 7,
		3: 4,
	}
	old, news, avg, middle := CalcDiffNumSneakers(a, b)
	fmt.Println(old, news, avg, middle)
}

func TestNumBelowTo(t *testing.T) {
	msg, price := NumBelowTo(map[int]int{
		1: 1480000,
		2: 1420000,
		3: 1620000,
		4: 2320000,
	})
	fmt.Println(msg, price)
}

func TestGSTPriceForBSC(t *testing.T) {
	price, pricebnb := GMTPriceForBSC()
	log.Println(price, pricebnb)
}

func TestCalcMintProfitForBSC(t *testing.T) {
	//log.Println(CalcMintProfitForBSC(1.37, 62))
	fmt.Println("sol gst, gmt", GSTPriceForSol(), GMTPriceForSol())
}

func TestGetFileContent(t *testing.T) {
	//var test = map[int]int{
	//	1: 2,
	//	2: 2,
	//	3: 2,
	//}
	//content, _ := json.Marshal(test)
	//Rewrite("abc1.txt", string(content))
	//a := GetFileContent("abc1.txt")
	//_ = json.Unmarshal([]byte(a), &test)
	//fmt.Println(test[1])
	//fmt.Println(BnbPrice(), SolPrice())
	//fmt.Println(IsAwesomeNum(998))

	//resp := RemoveDuplicateElement([]int{1, 3, 4, 5, 3, 4, 1})
	//fmt.Println(resp)

	//vars := map[string]string{
	//	"chain_name":          "BSC",
	//	"time":                "2022-06-28 12:11",
	//	"floor":               "13.45 bnb",
	//	"rate":                "-2%",
	//	"common_w_floor":      "10.00",
	//	"common_j_floor":      "12.00",
	//	"common_r_floor":      "14.23",
	//	"common_t_floor":      "16.11",
	//	"common_avg_floor":    "16.11",
	//	"uncommon_w_floor":    "12.22",
	//	"uncommon_j_floor":    "14.44",
	//	"uncommon_r_floor":    "15.32",
	//	"uncommon_t_floor":    "18.12",
	//	"uncommon_avg_floor":  "18.12",
	//	"rare_w_floor":        "19.53",
	//	"rare_j_floor":        "20.44",
	//	"rare_r_floor":        "21.98",
	//	"rare_t_floor":        "33.00",
	//	"rare_avg_floor":      "33.00",
	//	"epic_w_floor":        "34.13",
	//	"epic_j_floor":        "34.12",
	//	"epic_r_floor":        "34.55",
	//	"epic_t_floor":        "54.33",
	//	"epic_avg_floor":      "54.33",
	//	"legendary_w_floor":   "66.88",
	//	"legendary_j_floor":   "66.99",
	//	"legendary_r_floor":   "66.98",
	//	"legendary_t_floor":   "66.90",
	//	"legendary_avg_floor": "66.90",
	//}
	//template := "templates/sneaker-floor.html"
	//newFile := fmt.Sprintf("o-%s-sneaker-floor.html", chain)
	//newImage := "test1.jpg"
	//ReplaceVar(template, vars, newFile)
	//Html2Image(newFile, newImage)
	//PushFile(newImage)

	//vars := map[string]string{
	//	"chain_name":          "BSC",
	//	"time":                "2022-06-28 12:11",
	//	"total":               "12",
	//	"rate":                "-2%",
	//	"below_price_1":       "1.6b",
	//	"below_price_total_1": "20",
	//	"below_price_2":       "1.6b",
	//	"below_price_total_2": "20",
	//	"new_total":           "20",
	//	"consume_total":       "20",
	//	"new_avg_price":       "20",
	//	"new_middle_price":    "20",
	//	"common_w":            "20",
	//	"common_j":            "20",
	//	"common_r":            "20",
	//	"common_t":            "20",
	//	"uncommon_w":          "20",
	//	"uncommon_j":          "20",
	//	"uncommon_r":          "20",
	//	"uncommon_t":          "20",
	//	"rare_w":              "20",
	//	"rare_j":              "20",
	//	"rare_r":              "20",
	//	"rare_t":              "20",
	//	"epic_w":              "20",
	//	"epic_j":              "20",
	//	"epic_r":              "20",
	//	"epic_t":              "20",
	//	"legendary_w":         "20",
	//	"legendary_j":         "20",
	//	"legendary_r":         "20",
	//	"legendary_t":         "20",
	//	"notice":              "市场通缩中……",
	//}
	//template := "templates/sneaker-num.html"
	//newFile := fmt.Sprintf("o-%s-sneaker-num.html", chain)
	//newImage := "sneaker_num.jpg"
	//ReplaceVar(template, vars, newFile)
	//Html2Image(newFile, newImage)
	//PushFile(newImage)

	//vars := map[string]string{
	//	"chain_name":            "BSC",
	//	"time":                  "2022-06-28 12:11",
	//	"gst_price_u":           "12u",
	//	"gmt_price_u":           "12u",
	//	"mint_price_u":          "12u",
	//	"scroll_price_u":        "12u",
	//	"upgrade_price_u":       "12u",
	//	"sneaker_floor_price_u": "12u",
	//	"formula":               "1.08x220.46x0.94-(0.3257x360+0.8843x40+0.8843x2x43.90)-(20x0.3257+10x0.8843)=-21.82usd",
	//}
	//template := "templates/sneaker-mint.html"
	//newFile := fmt.Sprintf("o-%s-sneaker-mint.html", chain)
	//newImage := "sneaker_mint.jpg"
	//ReplaceVar(template, vars, newFile)
	//Html2Image(newFile, newImage)
	//PushFile(newImage)

	vars := map[string]string{
		"chain_name":      "BSC",
		"time":            "2022-06-28 12:11",
		"total":           "200",
		"rate":            "2%",
		"common_total":    "20",
		"uncommon_total":  "20",
		"rare_total":      "20",
		"epic_total":      "20",
		"legendary_total": "20",
		"common_price":    "20",
		"uncommon_price":  "20",
		"rare_price":      "20",
		"epic_price":      "20",
		"legendary_price": "20",
	}
	template := "templates/scroll.html"
	newFile := fmt.Sprintf("o-%s-scroll.html", chain)
	newImage := "scroll.jpg"
	ReplaceVar(template, vars, newFile)
	Html2Image(newFile, newImage)
	PushFile(newImage)
}
