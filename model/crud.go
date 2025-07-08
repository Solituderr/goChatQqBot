package model

import (
	"errors"
	"fmt"
	"time"
)

// 用户处理

func QueryLoveRate(uid string) (int, error) {
	var user UserInfo
	user.UserId = uid
	err := DB.Model(&UserInfo{}).Where("user_id = ?", uid).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("查询失败，请联系管理员。")
	}
	return user.LoveRate, nil
}

func UpdateLoveRate(uid string, increase int) error {
	var user UserInfo
	user.UserId = uid
	err := DB.Model(&UserInfo{}).Where("user_id = ?", uid).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("查询失败，请联系管理员。")
	}
	lr := user.LoveRate
	lr += increase
	if lr < 0 {
		lr = 0
	}
	err = DB.Model(&UserInfo{}).Where("user_id = ?", uid).Update("love_rate", lr).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("更新好感度失败，请联系管理员。")
	}
	return nil
}

func MinusSign(uid string) error {
	var user UserInfo
	user.UserId = uid
	err := DB.Model(&UserInfo{}).Where("user_id = ?", uid).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("查询用户失败，请联系管理员。")
	}
	st := user.SignTime
	if st == 0 {
		t := errors.New("签到次数用完")
		return t
	}
	st = st - 1
	err = DB.Model(&UserInfo{}).Where("user_id = ?", uid).Update("sign_time", st).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("更新礼物次数失败，请联系管理员。")
	}
	return nil
}

func MinusGift(uid string) error {
	var user UserInfo
	user.UserId = uid
	err := DB.Model(&UserInfo{}).Where("user_id = ?", uid).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("查询用户失败，请联系管理员。")
	}
	gt := user.GiftTime
	if gt == 0 {
		t := errors.New("送礼次数用完")
		return t
	}
	gt = gt - 1
	err = DB.Model(&UserInfo{}).Where("user_id = ?", uid).Update("gift_time", gt).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("更新礼物次数失败，请联系管理员。")
	}
	return nil
}

// 周边推送
func AddPush(p PushMsg) error {
	err := DB.Model(&PushMsg{}).Create(&p).Error
	return err
}

// 删除推送
func DelPush(p PushMsg) error {
	err := DB.Model(&PushMsg{}).Where("group_id = ? and msg_time = ?", p.GroupId, p.MsgTime).Delete(&p).Error
	return err
}

// 查询推送
func ListPush() ([]PushMsg, error) {
	var p []PushMsg
	err := DB.Model(&PushMsg{}).Find(&p).Error
	return p, err
}

func QueryTime() ([]time.Time, error) {
	var allTimes []time.Time
	err := DB.Model(&PushMsg{}).Pluck("msg_time", &allTimes).Error
	return allTimes, err
}

func QueryMsg(t time.Time) ([]PushMsg, error) {
	var p []PushMsg
	err := DB.Model(&PushMsg{}).Where("msg_time = ?", t).Find(&p).Error
	return p, err
}

// 改名字
func ChangeName(uid string, name string) error {
	var user UserInfo
	user.UserId = uid
	err := DB.Model(&UserInfo{}).Where("user_id = ?", uid).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Println(err)
		return errors.New("查询失败，请联系管理员。")
	}
	err = DB.Model(&UserInfo{}).Where("user_id = ?", uid).Update("name", name).Error
	return err
}

func CheckName(uid string) (string, error) {
	var user UserInfo
	user.UserId = uid
	err := DB.Model(&UserInfo{}).Where("user_id = ?", uid).FirstOrCreate(&user).Error
	if err != nil {
		fmt.Println(err)
		return "", errors.New("查询失败，请联系管理员。")
	}
	return user.Name, nil
}

// 重置数据
func ResetUser() error {
	var s = map[string]interface{}{
		"sign_time": 1,
		"gift_time": 3,
	}
	err := DB.Model(&UserInfo{}).Where("user_id > ?", 0).Updates(s).Error
	return err
}
