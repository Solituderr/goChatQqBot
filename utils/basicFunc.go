package utils

import (
	"fmt"
	"go-svc-tpl/model"
)

func JudgeWin(player1 string, player2 string) string {
	if player1 == player2 {
		return "平"
	}

	switch player1 {
	case "石头":
		if player2 == "剪刀" {
			return "赢"
		} else {
			return "输"
		}
	case "剪刀":
		if player2 == "布" {
			return "赢"
		} else {
			return "输"
		}
	case "布":
		if player2 == "石头" {
			return "赢"
		} else {
			return "输"
		}
	default:
		fmt.Println("无效的出拳")
		return "g"
	}
}

func CheckIfChangeName(uid string, nickname string) string {
	name, _ := model.CheckName(uid)
	if name == "name" {
		return nickname
	} else {
		return name
	}
}
