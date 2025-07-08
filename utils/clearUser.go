package utils

import (
	"fmt"
	"go-svc-tpl/model"
	"time"
)

// 零点的时候清空用户
func ClearUserInfo() {
	for {
		localTime := time.Now().Local()
		hour, min, _ := localTime.Clock()
		if hour == 0 && min == 0 {
			RcvGift = 0
			if err := model.ResetUser(); err != nil {
				fmt.Println(err)
			} else {
				time.Sleep(1 * time.Minute)
			}
		}
	}
}
