package utils

import (
	"fmt"
	"github.com/asmcos/requests"
	"go-svc-tpl/model"
)

type QqServe interface {
	SendMsg(cm model.CommonMsg, flag int) string
	AddQqFri(uid string, flag string) string
}

type Deal struct{}

// 响应msg id

func (Deal) SendMsg(cm model.CommonMsg, flag int) string {
	var resp *requests.Response
	//1表示群聊at 0表示私聊 2 表示群聊不at
	if flag == 1 {
		cqmsg := fmt.Sprintf("[CQ:at,qq=%s] %s", cm.UserId, cm.Message)
		url := fmt.Sprintf("http://%v:5700/send_group_msg", ip)
		data := requests.Datas{
			"group_id": cm.GroupId,
			"message":  cqmsg,
		}
		resp, _ = requests.Post(url, data)
	} else if flag == 0 {
		url := fmt.Sprintf("http://%v:5700/send_private_msg", ip)
		data := requests.Datas{
			"user_id": cm.UserId,
			"message": cm.Message,
		}
		resp, _ = requests.Post(url, data)
	} else {
		url := fmt.Sprintf("http://%v:5700/send_group_msg", ip)
		data := requests.Datas{
			"group_id": cm.GroupId,
			"message":  cm.Message,
		}
		resp, _ = requests.Post(url, data)
	}

	return resp.Text()
}

//无返回值

func (Deal) AddQqFri(uid string, flag string) string {
	url := fmt.Sprintf("http://%v:5700/set_friend_add_request", ip)
	data := requests.Datas{
		"flag":    flag,
		"approve": "true",
	}
	resp, _ := requests.Post(url, data)
	return resp.Text()
}
