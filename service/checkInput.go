package service

import "strings"

func IsMenu(s string) bool {
	m := GetJsonStr("menu")
	for k := range m {
		if s == k {
			return true
		}
	}
	return false
}

func IsQianDao(s string) bool {
	if s == "签到" || s == "天天宝签到" {
		return true
	}
	return false
}

func IsHaoGan(s string) bool {
	if s == "天天好感" || s == "天天宝好感" {
		return true
	}
	return false
}

func IsHuDong(s string) bool {
	if s == "摸天天" || s == "抱天天" || s == "拍天天" || s == "亲天天" {
		return true
	}
	return false
}

func IsToudian(s string) bool {
	if len(s) <= 2 {
		return false
	}
	if s[:2] == ".r" && strings.Contains(s, "d") {
		return true
	}
	return false
}

func IsStartShitou(s string) bool {
	if s == "和天天玩石头剪刀布" || s == "和天天宝玩石头剪刀布" || s == "和天天宝玩剪刀石头布" || s == "和天天玩剪刀石头布" {
		return true
	}
	return false
}

func IsShitou(s string) bool {
	if s == "我出石头" || s == "我出剪刀" || s == "我出布" {
		return true
	}
	return false
}

func IsYuanzuo(s string) bool {
	if s == "原著天天" || s == "原作天天" || s == "天天宝原作" || s == "天天宝原著" || s == "天天宝原文" || s == "天天原著" || s == "天天原作" {
		return true
	}
	return false
}

func IsShengRi(s string) bool {
	if s == "天天宝生日快乐" || s == "天天生日快乐" {
		return true
	}
	return false
}

func IsBiaoQing(s string) bool {
	if s == "天天宝表情" || s == "天天表情" || s == "天天宝表情包" || s == "天天表情包" {
		return true
	}
	return false
}
