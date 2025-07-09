package cronjob

import (
	"go-svc-tpl/app"
	"go-svc-tpl/logs"
	"go-svc-tpl/model"
	"go-svc-tpl/utils"
	"time"
)

func PushPeripheral() {
	qqclient := app.GetQQClient()
	qqServe := utils.GetQQServer()
	for {
		if !utils.PushOn[utils.PushGid] {
			time.Sleep(10 * time.Second)
			continue
		}
		localTime := time.Now().Local()
		hour, min, _ := localTime.Clock()

		var f = false
		var nowT time.Time
		if allt, err := model.QueryTime(); err != nil {
			//fmt.Println(err)
			logs.Error(err.Error())
		} else {
			for _, t := range allt {
				if t.Hour() == hour && t.Minute() == min {
					nowT = t
					f = true
				}
			}
		}
		if f {
			// 如果时间到，则找到所有满足该时刻的struct，发送
			msg, _ := model.QueryMsg(nowT)
			for _, v := range msg {
				send := v.Message
				gid := v.GroupId
				cm := model.CommonMsg{
					GroupId: gid,
					Message: send,
				}
				qqServe.SendMsg(cm, 2, qqclient)
			}
			time.Sleep(1 * time.Minute)
		}
		time.Sleep(3 * time.Second)
	}
}
