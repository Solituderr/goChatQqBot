package controller

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go-svc-tpl/model"
	"go-svc-tpl/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func ClassifyReq(c echo.Context) error {
	ch := new(model.CommonReq)
	//必须转化为二进制流，否则在读取一次后，json流关闭
	data, _ := io.ReadAll(c.Request().Body)

	if err := json.Unmarshal(data, ch); err != nil {
		logrus.Error("初始化解析失败：" + err.Error())
		return c.String(http.StatusOK, "jiexi gg")
	}
	//测试ch
	fmt.Println(ch)

	if ch.PostType == "message" {
		tMsg := new(model.TypeMsg)
		if err := json.Unmarshal(data, tMsg); err != nil {
			logrus.Error("初始化解析失败：" + err.Error())
			return c.String(http.StatusOK, "jiexi gg")
		}
		var flag int
		var cMsg model.CommonMsg
		var newMsg string
		if tMsg.MessageType == "group" {
			//转化为common message
			gMsg := new(model.GroupMsg)
			if err := json.Unmarshal(data, gMsg); err != nil {
				logrus.Error(err.Error())
				return c.String(http.StatusOK, "not ok")
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
				UserId:  strconv.FormatInt(gMsg.UserId, 10),
				GroupId: strconv.FormatInt(gMsg.GroupId, 10),
				Message: newMsg,
			}
		} else {
			flag = 0
			//转化为common message
			pMsg := new(model.PersonMsg)
			if err := json.Unmarshal(data, pMsg); err != nil {
				logrus.Error(err.Error())
				return c.String(http.StatusOK, "not ok")
			}
			newMsg = pMsg.Message
			cMsg = model.CommonMsg{
				UserId:  strconv.FormatInt(pMsg.UserId, 10),
				GroupId: "0",
				Message: newMsg,
			}
		}
		//根据 cMsg 和 flag 写功能类 要根据有无groupId选择不同的flag
		//无指定内容
		switch newMsg {
		case "/菜单":
			err := utils.AddMenu(cMsg, flag, c)
			return err
		case "/测":
			err := utils.AddTestContext(cMsg, flag, c)
			return err
		case "/clear":
			err := utils.AddDealClear(cMsg, flag, c)
			return err
		case "/clearSys":
			err := utils.AddDealClearSys(cMsg, flag, c)
			return err
		case "/listSet":
			err := utils.ListSetting(cMsg, flag, c)
			return err
		case "/updates":
			err := utils.AddUpdates(cMsg, flag, c)
			return err
		default:
			break
		}
		//有指定内容
		switch {
		case strings.Contains(newMsg, "/addAcc"):
			msg := newMsg[8:]
			cMsg.Message = msg
			err := utils.AddDealAccess(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "/delAcc"):
			uid := newMsg[8:]
			cMsg.Message = uid
			err := utils.AddDeleteAccess(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "/system"):
			define := newMsg[8:]
			cMsg.Message = define
			err := utils.AddChangeSys(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "/set"):
			index := newMsg[5:]
			cMsg.Message = index
			err := utils.GetSetting(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "/addSet"):
			m := newMsg[8:]
			cMsg.Message = m
			err := utils.AddSetting(cMsg, flag, c)
			return err
		default:
			//群组里需要at才行，不然不回复 或者 私聊直接发送
			if flag == 1 || flag == 0 {
				err := utils.AddServeText(cMsg, flag, c)
				return err
			} else {
				return c.String(http.StatusOK, "not at me")
			}
		}

	} else if ch.PostType == "request" {
		//这种其实是地址格式，需要用*转化为值
		var addFri = new(model.AddFriMsg)
		if err := json.Unmarshal(data, addFri); err != nil {
			logrus.Error("json 解析 gg" + err.Error())
			return c.String(http.StatusOK, "jiexi gg")
		}
		err := utils.AddServeQqFri(addFri, c)
		return err
	} else {
		return c.String(http.StatusOK, "no useful sending")
	}

}
