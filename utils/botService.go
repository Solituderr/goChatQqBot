package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go-svc-tpl/model"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// 添加好友，需首先添加权限
func AddServeQqFri(addfri *model.AddFriMsg, c echo.Context) error {
	uid := strconv.FormatInt(addfri.UserId, 10)
	fmt.Println(addfri)
	if addfri.RequestType == "friend" && model.CheckAccess(uid) == nil {
		t := qqServe.AddQqFri(uid, addfri.Flag)
		return c.String(http.StatusOK, "okk"+t)
	} else {
		logrus.Error("not friend req or no access")
		return c.String(http.StatusOK, "not friend request")
	}
}

// 添加权限 功能
func AddDealAccess(cm model.CommonMsg, flag int, c echo.Context) error {
	message := cm.Message
	uid := cm.UserId
	if uid == "742277855" {
		dealUid := message
		if err := model.AddAccess(dealUid); err != nil {
			logrus.Error(err.Error())
			return c.String(http.StatusOK, "addAcc gg")
		} else {
			cm.Message = "添加权限成功"
			qqServe.SendMsg(cm, flag)
			return c.String(http.StatusOK, "ok")
		}
	} else {
		cm.Message = "添加权限失败，你可能不是管理员"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "no access or not addAcc req")
	}
}

func AddDeleteAccess(cm model.CommonMsg, flag int, c echo.Context) error {
	message := cm.Message
	uid := cm.UserId
	if uid == "742277855" {
		dealUid := message
		if err := model.DeleteAccess(dealUid); err != nil {

			logrus.Error(err.Error())
			return c.String(http.StatusOK, "delAcc gg")
		} else {
			cm.Message = "删除权限成功"
			qqServe.SendMsg(cm, flag)
			return c.String(http.StatusOK, "ok")
		}
	} else {
		cm.Message = "删除权限失败，你可能不是管理员"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "no access or not addAcc req")
	}
}

// 群聊bot 功能

// 菜单
func AddMenu(cm model.CommonMsg, flag int, c echo.Context) error {
	//要填的菜单内容
	cm.Message = "本CBOTv3.0是由菜鸡开发\n目前功能：\n1.记忆对话\n2.更改设定\n3.群聊模式\n4.私聊模式(需要权限)\n5.增加设定(需要权限)\n6.权限设置(管理员)\n\nat机器人然后问问题(不能复制at，私聊模式不用)\n部分参数(可以也可以不at机器人后发送)：\n1./system [你需要的设定] （注意system后的空格）\n2./clear  清除会话，保留设定\n3./clearSys  清除设定及会话，系统会分配默认设定\n4./listSet  显示默认设定 \n5./set [序号]  选择默认设定（注意空格）\n6./addSet [概要] [内容]  添加设定（注意空格）\n7./updates     查看更新内容\n\n注：使用/system 和 /set 后请等待2条消息\n------请大家合法使用------"
	qqServe.SendMsg(cm, flag)
	return c.String(http.StatusOK, "okkk")
}

// 本次更新内容
func AddUpdates(cm model.CommonMsg, flag int, c echo.Context) error {
	cm.Message = "本次更新概要：\n1.优化了系统架构，增强代码可读性和可改性，但可能降低了响应速度。\n2.添加了私聊，默认设定，权限设置等功能，请注意空格。\n3.解决部分prompt重复回答问题，改进了返回error提示功能。\n4./set 和 /system 使用后会返回两个内容，一是清空确认，二是回复你的内容，请耐心等待。\n5.使用腾讯云函数反代chat，可能节约了成本。"
	qqServe.SendMsg(cm, flag)
	return c.String(http.StatusOK, "okkk")
}

// 测试bot是否运行
func AddTestContext(cm model.CommonMsg, flag int, c echo.Context) error {
	if cm.UserId == "742277855" {
		cm.Message = "测试成功！"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "ok")
	} else {
		return c.String(http.StatusOK, "not ok")
	}
}

// clear功能
func AddDealClear(cm model.CommonMsg, flag int, c echo.Context) error {
	err := model.DeleteMsg(cm.UserId, false)
	if err != nil {
		logrus.Error(err.Error())
		cm.Message = "当前对话清除失败"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "not ok")
	} else {
		cm.Message = "当前对话清除完成"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "ok")
	}
}

// clearSys功能
func AddDealClearSys(cm model.CommonMsg, flag int, c echo.Context) error {
	err := model.DeleteMsg(cm.UserId, true)
	if err != nil {
		cm.Message = "设定及对话清除失败"
		logrus.Error(err.Error())
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "not ok")
	} else {
		cm.Message = "设定及对话清除完成"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "ok")
	}
}

func AddChangeSys(cm model.CommonMsg, flag int, c echo.Context) error {
	uid := cm.UserId
	//先清除数据
	if err := model.DeleteMsg(uid, true); err != nil {
		logrus.Error(err.Error())
		cm1 := model.CommonMsg{
			UserId:  uid,
			GroupId: cm.GroupId,
			Message: "数据清除失败",
		}
		qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, "not ok")
	}
	//再设定
	var sysMsg = model.UserInfo{UserId: uid, Role: "system", Content: cm.Message}
	if err := model.SaveMsg(sysMsg); err != nil {
		logrus.Error(err.Error())
		cm.Message = "修改设定失败"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "not ok")
	} else {
		cm.Message = "修改设定完成，对话数据已清空"
		qqServe.SendMsg(cm, flag)
		return c.String(http.StatusOK, "ok")
	}
}

// chatgpt功能
func AddServeText(cm model.CommonMsg, flag int, c echo.Context) error {
	//一些需要用的数据
	uid := cm.UserId
	gid := cm.GroupId
	prompt := cm.Message
	//查看用户是否已经沟通过，没沟通过则初始化设定
	if err := model.CheckUser(uid); err != nil {
		err1 := model.SaveMsg(model.UserInfo{
			UserId:  uid,
			Role:    "system",
			Content: "你是一个万能助手",
		})
		if err1 != nil {
			return c.String(http.StatusOK, err1.Error())
		}
	}

	//问题存入database
	saveData := model.UserInfo{UserId: uid, Role: "user", Content: prompt}
	err := model.SaveMsg(saveData)
	if err != nil {
		var cm1 = model.CommonMsg{UserId: uid, GroupId: gid, Message: "数据库储存问题失败"}
		resp := qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, resp)
	}
	//找到全部prompt
	allMsg, err := model.FindMsg(uid)
	if err != nil {
		var cm1 = model.CommonMsg{UserId: uid, GroupId: gid, Message: "数据库寻找聊天记录失败"}
		resp := qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, resp)
	}
	//发送到python服务端
	gptReply := TestApi(cm, flag, allMsg)
	fmt.Println(gptReply)
	saveData1 := model.UserInfo{UserId: uid, Role: "assistant", Content: gptReply}
	err = model.SaveMsg(saveData1)
	if err != nil {
		var cm1 = model.CommonMsg{UserId: uid, GroupId: gid, Message: "回复保存到数据库失败"}
		resp := qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, resp)
	}
	cm.Message = gptReply
	resp := qqServe.SendMsg(cm, flag)
	return c.String(http.StatusOK, resp)
}

// 发送给gpt服务器
func TestApi(cm model.CommonMsg, flag int, users []model.UserInfo) string {
	//var sendPrompt []model.UserInfo
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		var cm1 = model.CommonMsg{UserId: cm.UserId, GroupId: cm.GroupId, Message: "json数据流化失败"}
		res := qqServe.SendMsg(cm1, flag)
		return res
	}
	url := fmt.Sprintf("http://%v:1928/chat", ip)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))

	if err != nil {
		var cm1 = model.CommonMsg{UserId: cm.UserId, GroupId: cm.GroupId, Message: "向python server传输数据失败"}
		res := qqServe.SendMsg(cm1, flag)
		return res
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	return string(bodyBytes)
}

//增加设定功能

//列举设定

func ListSetting(cm model.CommonMsg, flag int, c echo.Context) error {
	tmp := ""
	for i, data := range setting {
		tmp += fmt.Sprintf("%v : %v\n", i, data.Summary)
	}
	tmp += "发送/set 序号  选择设定"
	cm.Message = tmp
	qqServe.SendMsg(cm, flag)
	return c.String(http.StatusOK, "okkkk")
}

// 选择的设定  /set 序号

func GetSetting(cm model.CommonMsg, flag int, c echo.Context) error {
	index, _ := strconv.Atoi(cm.Message)
	if index >= len(setting) {
		var cm1 = model.CommonMsg{UserId: cm.UserId, GroupId: cm.GroupId, Message: "序号不对，别瞎搞！"}
		qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, "index error")
	}
	define := setting[index].Content
	cm.Message = define
	fmt.Println(cm)
	err := AddChangeSys(cm, flag, c)
	return err
}

// 添加设定 要权限  格式 /addSet [summary] [setting]

func AddSetting(cm model.CommonMsg, flag int, c echo.Context) error {
	split := strings.Split(cm.Message, " ")
	if len(split) > 2 {
		var cm1 = model.CommonMsg{UserId: cm.UserId, GroupId: cm.GroupId, Message: "格式错误，请重试"}
		qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, "geshi error")
	}
	if err := model.CheckAccess(cm.UserId); err != nil {
		var cm1 = model.CommonMsg{UserId: cm.UserId, GroupId: cm.GroupId, Message: "兄弟，你没权限，请向管理员申请"}
		qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, "add failed")
	} else {
		setting = append(setting, model.Setting{Summary: split[0], Content: split[1]})
		var cm1 = model.CommonMsg{UserId: cm.UserId, GroupId: cm.GroupId, Message: "添加已完成"}
		qqServe.SendMsg(cm1, flag)
		return c.String(http.StatusOK, "add success")
	}
}

//公用记录，群聊游戏

//自动删除clear聊天记录
