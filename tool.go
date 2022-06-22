package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"gopkg.in/fatih/set.v0"
)

func Insert(filePath string, value string) {
	file, err := os.OpenFile(fmt.Sprintf("./db/%s", filePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(value + "\n")
	write.Flush()
}

func FindLatest(filePath string) string {
	file, err := os.OpenFile(fmt.Sprintf("./db/%s", filePath), os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
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
