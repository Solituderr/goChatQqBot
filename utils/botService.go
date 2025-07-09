package utils

import (
	"fmt"
	"go-svc-tpl/model"
	"go-svc-tpl/service"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LagrangeDev/LagrangeGo/client"
)

// 控制两个天天宝
func AddEnableBot1(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	Enable1[cm.GroupId] = true
	cm.Message = "天天宝1号已启用"
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}
func AddEnableBot2(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	Enable2[cm.GroupId] = true
	cm.Message = "天天宝2号已启用"
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// bot开机关机
func AddStartBot(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	IsOn[cm.GroupId] = true
	cm.Message = "天天宝已开机"
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

func AddStopBot(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	IsOn[cm.GroupId] = false
	cm.Message = "天天宝已关机"
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 测试bot是否运行
func AddTestContext(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId == "742277855" {
		cm.Message = "测试成功！"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	} else {
		return nil
	}
}

// 菜单
func AddMenu(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	m := service.GetJsonStr("menu")
	if strings.Contains(cm.Message, "菜单") {
		rsp := m["功能菜单"]
		cm.Message = rsp.(string)
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	} else {
		rsp := m[cm.Message]
		cm.Message = rsp.(string)
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	}
}

// 签到
func AddSignIn(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	var f = false
	m := service.GetJsonStr("qiandao")
	//m1 := m["签到"].([]string)
	mSign := service.RandMessage(m["签到"])
	mNoSign := strings.Replace(mSign, "{nick}", nickname, 1)
	mAfter := m["签到后"].(string)
	if strings.Contains(mNoSign, "来得早不如来得巧") {
		mAfter = "*天天宝将外卖一口吃完*\n" + mAfter
		f = true
	} else {
		mAfter = mNoSign + "\n" + mAfter
	}
	m2 := m["已签到"].(string)
	mHasSign := strings.Replace(m2, "{nick}", nickname, 1)
	if err := model.MinusSign(cm.UserId); err != nil {
		if strings.Contains(err.Error(), "签到次数") {
			cm.Message = mHasSign
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		} else {
			cm.Message = err.Error()
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		}
	} else {
		num := service.RandNum(50, 100)
		err = model.UpdateLoveRate(cm.UserId, num)
		if f {
			cm.Message = mNoSign
			qqServe.SendMsg(cm, flag, qqclient)
			cm.Message = mAfter
			qqServe.SendMsg(cm, flag, qqclient)
		} else {
			cm.Message = mAfter
			qqServe.SendMsg(cm, flag, qqclient)
		}
		return nil
	}

}

// 好感
func AddSelectLove(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	if favor, err := model.QueryLoveRate(cm.UserId); err != nil {
		cm.Message = err.Error()
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	} else {
		if favor < 500 {
			cm.Message = fmt.Sprintf("对%s的好感度只有%d，要多多来和本剑圣聊天！", nickname, favor)
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		} else if favor < 1000 {
			cm.Message = fmt.Sprintf("对%s的好感度有%d，谢谢你一直支持本剑圣", nickname, favor)
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		} else if favor < 5000 {
			cm.Message = fmt.Sprintf("好感度到%d了，不愧是%s！我的忠实粉丝！", favor, nickname)
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		} else {
			cm.Message = fmt.Sprintf("对%s的好感度已经有%v了，宣布你就是我的粉丝中第一名", nickname, favor)
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		}
	}
}

// 交互
func AddAction(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	m := service.GetJsonStr("jiaohu")
	reply := service.RandMessage(m[cm.Message])
	reply = strings.Replace(reply, "{nick}", nickname, 1)
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 一会吃啥
func AddEatWhat(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	m := service.GetJsonStr("chisha")
	reply := m["一会吃啥"].(string)
	eat := service.RandMessage(m["吃啥"])
	reply = strings.Replace(reply, "{chisha}", eat, 1)
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 随机骰点
func AddDiceRand(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	m := cm.Message
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	lis := strings.Split(m, " ")
	if len(lis) == 1 {
		return nil
	}
	name := lis[1]

	re := regexp.MustCompile("\\.r(\\d*)d(\\d*)")
	match := re.FindStringSubmatch(m)
	var diceNum int
	var score int
	if match[1] == "" {
		diceNum = 1
		score, _ = strconv.Atoi(match[2])
	} else if match[2] == "" {
		diceNum, _ = strconv.Atoi(match[1])
		score = 100
	} else {
		diceNum, _ = strconv.Atoi(match[1])
		score, _ = strconv.Atoi(match[2])
	}
	var touzi = fmt.Sprintf("R%dD%d=", diceNum, score)
	rand.Seed(time.Now().UnixNano())
	tmp := rand.Perm(score)
	for i := 0; i < diceNum; i++ {
		num := tmp[i]
		touzi = touzi + strconv.Itoa(num)
		if i != diceNum-1 {
			touzi += "，"
		}
	}
	reply := fmt.Sprintf("%s掷骰 %s: %s", nickname, name, touzi)
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 石头剪刀布
func AddStartRock(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	s := "来石头剪刀布吧~\n发送 我出石头/剪刀/布 来和天天宝比划比划,三局两胜也是可以的哟！"
	cm.Message = s
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

func AddRockGame(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	s := strings.Replace(cm.Message, "我出", "", 1)
	var s1 = []string{"石头", "剪刀", "布"}
	num := service.RandNum(0, 2)
	me := s1[num]
	tiantian := "天天宝出" + me + "!\n"
	iswin := JudgeWin(me, s)
	m := service.GetJsonStr("shitou")
	reply := m[iswin].(string)
	reply = strings.Replace(reply, "{nick}", nickname, 1)
	cm.Message = tiantian + reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 周易算卦
func AddTellerZhou(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	m := service.GetJsonStr("suangua")
	reply := m["周易算卦"].(string)
	msg := service.RandMessage(m["六十四卦"])
	lis := strings.Split(msg, "抽得卦爻：{")
	s1 := lis[0]
	s2 := lis[1]
	s2 = strings.Replace(s2, "}", "", 1)
	msg1 := service.RandMessage(m[s2])
	cm.Message = reply + s1 + msg1
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// abo测试
func AddABOTest(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	m := service.GetJsonStr("abo")
	reply := m["abo"].(string)
	gender := service.RandMessage(m["abo性别"])
	smell := service.RandMessage(m["abo味道"])
	reply = strings.Replace(reply, "{nick}", nickname, 1)
	reply = strings.Replace(reply, "{abo性别}", gender, 1)
	reply = strings.Replace(reply, "{abo味道}", smell, 1)
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// cp关键词
func AddCPKeyWord(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	m := service.GetJsonStr("cpword")
	reply := m["cp梗"].(string)
	s := m["cp"]
	word1 := service.RandMessage(s)
	word2 := service.RandMessage(s)
	word3 := service.RandMessage(s)
	replacements := map[string]string{
		"{cp1}": word1,
		"{cp2}": word2,
		"{cp3}": word3,
	}

	for old, n := range replacements {
		reply = strings.ReplaceAll(reply, old, n)
	}
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// cp名
func AddCPName(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	msg := strings.Replace(cm.Message, "天天宝今日", "", 1)
	m := service.GetJsonStr("cpname")
	what := service.RandMessage(m["what"])
	how := service.RandMessage(m["how"])
	where := service.RandMessage(m["where"])
	reply := fmt.Sprintf("今天为%s抽取的%s关键词如下：\n%s\n%s\n%s", nickname, msg, what, how, where)
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 天天宝原著
func AddOriginalWork(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	m := service.GetJsonStr("yuanzhu")
	reply := service.RandMessage(m["原著"])
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 生日
func AddBirthday(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	m := service.GetJsonStr("shengri")
	reply := service.RandMessage(m["生日"])
	reply = strings.Replace(reply, "{nick}", nickname, 1)
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 发表情 TODO
func AddSendEmoji(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	picpath := GetRandPic()
	if picpath == "g" {
		cm.Message = "出错了"
		qqServe.SendMsg(cm, 2, qqclient)
		return nil
	} else {
		cqCode := fmt.Sprintf("[CQ:image,file=file:///%v]", picpath)
		cm.Message = cqCode
		qqServe.SendMsg(cm, 2, qqclient)
		return nil
	}
}

// 礼物
func AddRcvGift(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	gift := strings.Replace(cm.Message, "送天天宝", "", 1)
	notRcv := []string{"肖战", "秋葵", "刘皓", "黄文"}
	var f = false
	var addfavor = 5
	for _, s := range notRcv {
		if gift == s {
			f = true
			addfavor = -10
			break
		}
	}

	if f {
		m := service.GetJsonStr("liwu")
		reply := m[gift].(string)
		reply = strings.Replace(reply, "{nick}", nickname, -1)
		if err := model.UpdateLoveRate(cm.UserId, addfavor); err != nil {
			cm.Message = err.Error()
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		} else {
			cm.Message = reply
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		}
	} else {
		if gift == "" || gift == "礼物" {
			cm.Message = fmt.Sprintf("嗯嗯嗯？%s要送天天宝什么？", nickname)
			qqServe.SendMsg(cm, flag, qqclient)
			return nil
		} else {
			if err := model.MinusGift(cm.UserId); err != nil {
				if strings.Contains(err.Error(), "送礼次数用完") {
					cm.Message = fmt.Sprintf("可以了可以了，%s今天已经送得太多了，已经收到你的心意了！你站在此不要走动，本剑圣去给你买杯奶茶！", nickname)
					qqServe.SendMsg(cm, flag, qqclient)
					return nil
				} else {
					cm.Message = err.Error()
					qqServe.SendMsg(cm, flag, qqclient)
					return nil
				}
			} else {
				RcvGift += 1
				reply := fmt.Sprintf("哇塞！居然是我最喜欢的%s!!%s真是宇宙无敌帅气大好人，虽然只比我差了一点点哈哈哈哈！~\n今日收礼:%d件", gift, nickname, RcvGift)
				rate, _ := model.QueryLoveRate(cm.UserId)
				if rate+addfavor > 100 && rate <= 100 {
					addfavor = 100 - rate
				}
				if rate > 100 {
					addfavor = 0
				}
				err := model.UpdateLoveRate(cm.UserId, addfavor)
				fmt.Println(err)
				cm.Message = reply
				qqServe.SendMsg(cm, flag, qqclient)
				return nil
			}
		}
	}
}

// 添加推送
func AddPush(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	t := cm.Message
	split := strings.SplitN(t, " ", 3)

	if len(split) != 3 {
		cm.Message = "格式不正确"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	}
	gid := split[0]
	send := split[2]
	lis := strings.Split(split[1], ":")
	h, _ := strconv.Atoi(lis[0])
	m, _ := strconv.Atoi(lis[1])
	addTime := time.Date(2023, 1, 1, h, m, 0, 0, time.Local)
	pm := model.PushMsg{
		MsgTime: addTime,
		Message: send,
		GroupId: gid,
	}
	if err := model.AddPush(pm); err != nil {
		cm.Message = "添加推送失败"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	} else {
		cm.Message = "添加推送成功"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	}
}

func AddDelPush(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	t := cm.Message
	split := strings.SplitN(t, " ", 2)
	if len(split) != 2 {
		cm.Message = "格式不正确"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	}
	gid := split[0]
	lis := strings.Split(split[1], ":")
	h, _ := strconv.Atoi(lis[0])
	m, _ := strconv.Atoi(lis[1])
	addTime := time.Date(2023, 1, 1, h, m, 0, 0, time.Local)
	pm := model.PushMsg{
		MsgTime: addTime,
		GroupId: gid,
	}
	if err := model.DelPush(pm); err != nil {
		cm.Message = "删除推送失败"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	} else {
		cm.Message = "删除推送成功"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	}

}

func AddListPush(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if push, err := model.ListPush(); err != nil {
		cm.Message = "查询失败"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	} else {
		for _, v := range push {
			t := fmt.Sprintf("时间：%d:%d\n", v.MsgTime.Hour(), v.MsgTime.Minute())
			g := fmt.Sprintf("群号：%s\n", v.GroupId)
			m := fmt.Sprintf("消息：%s", v.Message)
			cm.Message = t + g + m
			qqServe.SendMsg(cm, flag, qqclient)
		}
		return nil
	}
}

func AddListPushTime(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	queryTime, _ := model.QueryTime()
	var reply = "所有推送时间：\n"
	for _, t := range queryTime {
		hour := t.Hour()
		minute := t.Minute()
		reply += fmt.Sprintf("%d时%d分\n", hour, minute)
	}
	cm.Message = reply[:len(reply)-1]
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

func AddStartPush(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	PushOn[PushGid] = true
	cm.Message = "推送功能已开启"
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

func AddStopPush(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	if cm.UserId != authQq {
		return nil
	}
	PushOn[PushGid] = false
	cm.Message = "推送功能已关闭"
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

// 改称呼
func AddChangeName(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	name := cm.Message
	if err := model.ChangeName(cm.UserId, name); err != nil {
		cm.Message = "修改称呼失败"
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	} else {
		cm.Message = fmt.Sprintf("已将%s的名称更改为%s", nickname, name)
		qqServe.SendMsg(cm, flag, qqclient)
		return nil
	}
}

// 幸运值
func AddLuckyNum(cm model.CommonMsg, flag int, qqclient *client.QQClient) error {
	nickname := cm.Sender.NickName
	nickname = CheckIfChangeName(cm.UserId, nickname)
	num := service.RandNum(0, 100)
	reply := fmt.Sprintf("看剑看剑看剑！本剑圣宣布%s今天的幸运值是%d", nickname, num)
	cm.Message = reply
	qqServe.SendMsg(cm, flag, qqclient)
	return nil
}

//所有图片爬取

//公用记录，群聊游戏

//自动删除clear聊天记录
