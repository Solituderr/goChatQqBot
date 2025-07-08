package controller

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go-svc-tpl/model"
	"go-svc-tpl/service"
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
				Sender:  gMsg.Sender,
			}
		} else {
			return c.String(http.StatusOK, "not group")
		}
		//根据 cMsg 和 flag 写功能类 要根据有无groupId选择不同的flag
		//无指定内容
		switch newMsg {
		case ".enable1":
			err := utils.AddEnableBot1(cMsg, flag, c)
			return err
			//case ".enable2":
			//	err := utils.AddEnableBot2(cMsg, flag, c)
			//	return err
		}

		if !utils.Enable1[cMsg.GroupId] {
			return c.String(http.StatusOK, "okk")
		}
		//if !utils.Enable2[cMsg.GroupId] {
		//	return c.String(http.StatusOK, "okk")
		//}

		switch newMsg {
		case ".bot on":
			err := utils.AddStartBot(cMsg, flag, c)
			return err
		case ".bot off":
			err := utils.AddStopBot(cMsg, flag, c)
			return err
		}
		if !utils.IsOn[cMsg.GroupId] {
			return c.String(http.StatusOK, "okk")
		}
		switch newMsg {
		case "/测":
			err := utils.AddTestContext(cMsg, flag, c)
			return err
		case ".draw一会吃啥":
			err := utils.AddEatWhat(cMsg, flag, c)
			return err
		case ".draw 一会吃啥":
			err := utils.AddEatWhat(cMsg, flag, c)
			return err
		case ".draw 周易算卦":
			err := utils.AddTellerZhou(cMsg, flag, c)
			return err
		case ".draw abo":
			err := utils.AddABOTest(cMsg, flag, c)
			return err
		case ".draw cp关键词":
			err := utils.AddCPKeyWord(cMsg, flag, c)
			return err
		case "/listTime":
			err := utils.AddListPushTime(cMsg, flag, c)
			return err
		case ".jrrp":
			err := utils.AddLuckyNum(cMsg, flag, c)
			return err
		case "/listPush":
			err := utils.AddListPush(cMsg, flag, c)
			return err
		case "/startPush":
			err := utils.AddStartPush(cMsg, flag, c)
			return err
		case "/stopPush":
			err := utils.AddStopPush(cMsg, flag, c)
			return err
		}
		//有指定内容
		switch {
		case service.IsMenu(cMsg.Message):
			err := utils.AddMenu(cMsg, flag, c)
			return err
		case service.IsQianDao(cMsg.Message):
			err := utils.AddSignIn(cMsg, flag, c)
			return err
		case service.IsHaoGan(cMsg.Message):
			err := utils.AddSelectLove(cMsg, flag, c)
			return err
		case service.IsHuDong(cMsg.Message):
			err := utils.AddAction(cMsg, flag, c)
			return err
		case service.IsToudian(cMsg.Message):
			err := utils.AddDiceRand(cMsg, flag, c)
			return err
		case service.IsStartShitou(cMsg.Message):
			err := utils.AddStartRock(cMsg, flag, c)
			return err
		case service.IsShitou(cMsg.Message):
			err := utils.AddRockGame(cMsg, flag, c)
			return err
		case service.IsYuanzuo(cMsg.Message):
			err := utils.AddOriginalWork(cMsg, flag, c)
			return err
		case service.IsShengRi(cMsg.Message):
			err := utils.AddBirthday(cMsg, flag, c)
			return err
		case service.IsBiaoQing(cMsg.Message):
			err := utils.AddSendEmoji(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "天天宝今日"):
			err := utils.AddCPName(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "送天天宝"):
			err := utils.AddRcvGift(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "/addPush"):
			cMsg.Message = newMsg[9:]
			err := utils.AddPush(cMsg, flag, c)
			return err
		case strings.Contains(newMsg, "/delPush"):
			cMsg.Message = newMsg[9:]
			err := utils.AddDelPush(cMsg, flag, c)
			return err
		case newMsg[:3] == ".nn":
			cMsg.Message = newMsg[4:]
			err := utils.AddChangeName(cMsg, flag, c)
			return err
		default:
			return c.String(http.StatusOK, "ok")
		}
	} else if ch.PostType == "notice" {
		nMsg := new(model.EnterGroup)

		if err := json.Unmarshal(data, nMsg); err != nil {
			logrus.Error(err.Error())
			return c.String(http.StatusOK, "not ok")
		}

		gid := strconv.Itoa(int(nMsg.GroupId))

		if !utils.Enable1[gid] {
			return c.String(http.StatusOK, "okk")
		}
		//if !utils.Enable2[gid] {
		//	return c.String(http.StatusOK, "okk")
		//}
		if nMsg.NoticeType == "group_increase" {
			err := utils.AddWelcomePerson(*nMsg, 1, c)
			return err
		}
		return c.String(http.StatusOK, "okk")
	} else {
		return c.String(http.StatusOK, "no useful sending")
	}

}
