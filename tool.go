package main

import (
	"bufio"
	"fmt"
	"os"
)

func StoreJson(filePath string, jsonStr string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(jsonStr)
	write.Flush()
}

func LoadJson(filePath string) (byte, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	return reader.ReadByte()
}
