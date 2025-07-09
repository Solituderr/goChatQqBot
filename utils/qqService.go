package utils

import (
	"go-svc-tpl/logs"
	"go-svc-tpl/model"
	"strconv"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
)

type QqServe interface {
	SendMsg(cm model.CommonMsg, flag int, client *client.QQClient)
	AddQqFri(uid string, flag string) string
	BanPeople(cm model.CommonMsg) string
}

type Deal struct{}

func GetQQServer() QqServe {
	return qqServe
}

// 响应msg id

func (Deal) SendMsg(cm model.CommonMsg, flag int, client *client.QQClient) {
	//1表示群聊at 0表示私聊 2 表示群聊不at
	uid1, err := strconv.Atoi(cm.UserId)
	if err != nil {
		logs.Error("[SendMsg] conv uid i64 to str error")
		return
	}
	uid := uint32(uid1)
	gid1, err := strconv.Atoi(cm.GroupId)
	if err != nil {
		logs.Error("[SendMsg] conv gid i64 to str error")
		return
	}
	gid := uint32(gid1)
	if flag == 1 {
		atMsg := message.NewAt(uid)
		msg := message.NewText(cm.Message)
		_, err := client.SendGroupMessage(gid, []message.IMessageElement{atMsg, msg})
		if err != nil {
			logs.Error("[SendMsg] %s", err.Error())
			return
		}
	} else if flag == 0 {
		msg := message.NewText(cm.Message)
		_, err := client.SendPrivateMessage(uid, []message.IMessageElement{msg})
		if err != nil {
			logs.Error("[SendMsg] %s", err.Error())
			return
		}
	} else {
		msg := message.NewText(cm.Message)
		_, err := client.SendGroupMessage(gid, []message.IMessageElement{msg})
		if err != nil {
			logs.Error("[SendMsg] %s", err.Error())
			return
		}
	}

	return
}

func (Deal) BanPeople(cm model.CommonMsg) string {
	return ""
}

//无返回值

func (Deal) AddQqFri(uid string, flag string) string {
	return ""
}
