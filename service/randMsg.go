package service

import (
	"fmt"
	"math/rand"
	"time"
)

func RandNum(a int, b int) int {
	// 使用当前时间作为随机数种子，确保每次运行程序都生成不同的随机数
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(b-a) + a
	time.Sleep(1 * time.Microsecond)
	return randNum
}

func RandMessage(msg interface{}) string {
	// 使用当前时间作为随机数种子，确保每次运行程序都生成不同的随机数
	slice := []string{}
	for _, item := range msg.([]interface{}) {
		slice = append(slice, fmt.Sprintf("%v", item))
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(slice))
	fmt.Println(randNum)
	time.Sleep(1 * time.Microsecond)
	return slice[randNum]
}
