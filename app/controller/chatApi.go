package controller

import (
	"go-svc-tpl/logs"
	"go-svc-tpl/model"
	"go-svc-tpl/service"
	"go-svc-tpl/utils"
	"strconv"
	"strings"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
)

func ClassifyReq(qqclient *client.QQClient) error {

	// TODO(Echo) 订阅群消息事件
	qqclient.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		if event.ToString() == "114514" {
			img, _ := message.NewFileImage("testgroup.png")
			_, err := client.SendGroupMessage(event.GroupUin, []message.IMessageElement{img})
			if err != nil {
				return
			}
		}
	})

	qqclient.GroupNotifyEvent.Subscribe(func(client *client.QQClient, event event.INotifyEvent) {

	})

	qqclient.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		// 先打印出消息的结构信息
		logs.Info("msg: %s", event.ToString())
		var flag int
		var cMsg model.CommonMsg
		var newMsg string
		//转化为common message
		gMsg := &model.GroupMsg{
			Message: event.ToString(),
			Sender: struct {
				UserId   int64  "json:\"user_id\""
				NickName string "json:\"nickname\""
				Sex      string "json:\"sex\""
				Age      int32  "json:\"age\""
			}{
				UserId:   int64(event.Sender.Uin),
				NickName: event.Sender.Nickname,
			},
		}
		//判断是否为at
		if strings.Contains(gMsg.Message, utils.AtQqUid) {
			//为at
			newMsg = gMsg.Message[utils.LenAtQqUid:]
			flag = 1
		} else {
			//非at
			newMsg = gMsg.Message
			flag = 2
		}

		cMsg = model.CommonMsg{
			UserId:  strconv.Itoa(int(event.Sender.Uin)),
			GroupId: strconv.Itoa(int(event.GroupUin)),
			Message: newMsg,
			Sender:  gMsg.Sender,
		}
		//根据 cMsg 和 flag 写功能类 要根据有无groupId选择不同的flag
		//无指定内容
		var err error
		switch newMsg {
		case ".enable1":
			err = utils.AddEnableBot1(cMsg, flag, qqclient)
			return
			//case ".enable2":
			//	err := utils.AddEnableBot2(cMsg, flag, qqclient)
			//	return
		}

		if !utils.Enable1[cMsg.GroupId] {
			logs.Info("not enable bot")
			return
		}
		//if !utils.Enable2[cMsg.GroupId] {
		//	return c.String(http.StatusOK, "okk")
		//}

		switch newMsg {
		case ".bot on":
			utils.AddStartBot(cMsg, flag, qqclient)
			return
		case ".bot off":
			err = utils.AddStopBot(cMsg, flag, qqclient)
			return
		}
		if !utils.IsOn[cMsg.GroupId] {
			logs.Info("not on bot")
			return
		}
		switch newMsg {
		case "/测":
			err = utils.AddTestContext(cMsg, flag, qqclient)
			return
		case ".draw一会吃啥":
			err = utils.AddEatWhat(cMsg, flag, qqclient)
			return
		case ".draw 一会吃啥":
			err = utils.AddEatWhat(cMsg, flag, qqclient)
			return
		case ".draw 周易算卦":
			err = utils.AddTellerZhou(cMsg, flag, qqclient)
			return
		case ".draw abo":
			err = utils.AddABOTest(cMsg, flag, qqclient)
			return
		case ".draw cp关键词":
			err = utils.AddCPKeyWord(cMsg, flag, qqclient)
			return
		case "/listTime":
			err = utils.AddListPushTime(cMsg, flag, qqclient)
			return
		case ".jrrp":
			err = utils.AddLuckyNum(cMsg, flag, qqclient)
			return
		case "/listPush":
			err = utils.AddListPush(cMsg, flag, qqclient)
			return
		case "/startPush":
			err = utils.AddStartPush(cMsg, flag, qqclient)
			return
		case "/stopPush":
			err = utils.AddStopPush(cMsg, flag, qqclient)
			return
		}
		//有指定内容
		switch {
		case service.IsMenu(cMsg.Message):
			err = utils.AddMenu(cMsg, flag, qqclient)
			return
		case service.IsQianDao(cMsg.Message):
			err = utils.AddSignIn(cMsg, flag, qqclient)
			return
		case service.IsHaoGan(cMsg.Message):
			err = utils.AddSelectLove(cMsg, flag, qqclient)
			return
		case service.IsHuDong(cMsg.Message):
			err = utils.AddAction(cMsg, flag, qqclient)
			return
		case service.IsToudian(cMsg.Message):
			err = utils.AddDiceRand(cMsg, flag, qqclient)
			return
		case service.IsStartShitou(cMsg.Message):
			err = utils.AddStartRock(cMsg, flag, qqclient)
			return
		case service.IsShitou(cMsg.Message):
			err = utils.AddRockGame(cMsg, flag, qqclient)
			return
		case service.IsYuanzuo(cMsg.Message):
			err = utils.AddOriginalWork(cMsg, flag, qqclient)
			return
		case service.IsShengRi(cMsg.Message):
			err = utils.AddBirthday(cMsg, flag, qqclient)
			return
		case service.IsBiaoQing(cMsg.Message):
			err = utils.AddSendEmoji(cMsg, flag, qqclient)
			return
		case strings.Contains(newMsg, "天天宝今日"):
			err = utils.AddCPName(cMsg, flag, qqclient)
			return
		case strings.Contains(newMsg, "送天天宝"):
			err = utils.AddRcvGift(cMsg, flag, qqclient)
			return
		case strings.Contains(newMsg, "/addPush"):
			cMsg.Message = newMsg[9:]
			err = utils.AddPush(cMsg, flag, qqclient)
			return
		case strings.Contains(newMsg, "/delPush"):
			cMsg.Message = newMsg[9:]
			err = utils.AddDelPush(cMsg, flag, qqclient)
			return
		case newMsg[:3] == ".nn":
			cMsg.Message = newMsg[4:]
			err = utils.AddChangeName(cMsg, flag, qqclient)
			return
		default:
		}
		if err != nil {
			logs.Error("[ClassifyReq] %v", err)
		}
		return
	})

	qqclient.GroupJoinEvent.Subscribe(func(client *client.QQClient, event *event.GroupMemberIncrease) {
		gid := strconv.Itoa(int(event.GroupUin))

		if !utils.Enable1[gid] {
			return
		}
		//if !utils.Enable2[gid] {
		//	return c.String(http.StatusOK, "okk")
		//}
		nMsg := model.EnterGroup{
			GroupId: int64(event.GroupUin),
			UserId:  int64(event.UserUin),
		}
		err := utils.AddWelcomePerson(nMsg, 1, qqclient)
		if err != nil {
			logs.Error("[GroupJoinEvent] %v", err)
		}
		return
	})

	qqclient.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		logs.Info("连接已断开：%v", event.Message)
	})

	return nil
}
