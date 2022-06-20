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
