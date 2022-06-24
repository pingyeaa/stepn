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

	genesShoes = append(genesShoes, &Shoe{
		ID:        1,
		TypeID:    601,
		Otd:       1234,
		SellPrice: 9999999,
		Level:     5,
		Quantity:  1,
		Mint:      2,
	})

	genesShoes = append(genesShoes, &Shoe{
		ID:        2,
		TypeID:    603,
		Otd:       1235,
		SellPrice: 2999999,
		Level:     5,
		Quantity:  1,
		Mint:      2,
	})

	genesShoes = append(genesShoes, &Shoe{
		ID:        3,
		TypeID:    601,
		Otd:       1034,
		SellPrice: 9999999,
		Level:     3,
		Quantity:  1,
		Mint:      0,
	})
	res := GenesShoes()
	log.Print(res)
}
