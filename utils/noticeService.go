package utils

import (
	"go-svc-tpl/model"
	"go-svc-tpl/service"
	"strconv"

	"github.com/LagrangeDev/LagrangeGo/client"
)

func AddWelcomePerson(cm model.EnterGroup, flag int, qqclient *client.QQClient) error {
	uid := strconv.Itoa(int(cm.UserId))
	gid := strconv.Itoa(int(cm.GroupId))
	m := service.GetJsonStr("ruqun")
	var reply string

	if gid == "773444291" {
		reply = m["入群欢迎2"].(string)
	} else if gid == "658227422" || gid == "429093558" {
		reply = m["入群欢迎1"].(string)
	} else {
		return nil
	}

	cmg := model.CommonMsg{
		UserId:  uid,
		GroupId: gid,
		Message: reply,
	}
	qqServe.SendMsg(cmg, flag, qqclient)
	return nil
}
