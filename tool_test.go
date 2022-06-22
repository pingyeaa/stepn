package main

import (
	"fmt"
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
