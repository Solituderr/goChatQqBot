package cronjob

import (
	"bufio"
	"fmt"
	"go-svc-tpl/app"
	"go-svc-tpl/model"
	"go-svc-tpl/utils"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var WeiBoId = make(map[string]string)

func GetUserName(data string) string {
	res := gjson.Get(data, "data.userInfo.screen_name")
	return res.String()
}

func GetContainerId(data string) string {
	var container string
	res := gjson.Get(data, "data.tabsInfo.tabs")
	for _, m := range res.Array() {
		ma := m.Map()
		if ma["tab_type"].Str == "weibo" {
			container = ma["containerid"].Str
		}
	}
	return container

}

func GetWei(data string, uid string) string {
	t := gjson.Get(data, "data.cards")
	for _, v := range t.Array() {
		if v.Get("card_type").Int() == 9 && !strings.Contains(v.Get("profile_type_id").Str, "top") {
			id := v.Get("mblog").Get("id").Str
			fmt.Printf("当前uid：%s，当前消息id：%s\n", uid, id)
			// 更新了
			if id != WeiBoId[uid] {
				//createTime := v.Get("mblog").Get("created_at").Time()
				text := v.Get("mblog").Get("text").Str
				fmt.Printf("当前uid：%s 检测到新微博。\n", uid)
				WeiBoId[uid] = id
				return text
			} else {
				return ""
			}
		}
	}
	return ""
}

func GetWeibo(uid string) {
	qqclient := app.GetQQClient()
	qqServe := utils.GetQQServer()
	url := "https://m.weibo.cn/api/container/getIndex"
	client := req.C()
	resp, err := client.R().SetQueryParams(map[string]string{
		"type":  "uid",
		"value": uid,
	}).Get(url)
	//username := GetUserName(resp.String())
	container := GetContainerId(resp.String())
	//fmt.Println(container)
	//fmt.Printf("用户名为：%s", username)
	wbData, err := client.R().SetQueryParams(map[string]string{
		"type":        "uid",
		"value":       uid,
		"containerid": container,
	}).Get(url)
	//fmt.Println(wbData.String())

	if err != nil {
		logrus.Error(err)
	}
	s := GetWei(wbData.String(), uid)
	s = regexp.MustCompile("<a.*?a>").ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "<br />", "\n")
	msg := bluemonday.StrictPolicy().Sanitize(s)
	if msg != "" {
		send := fmt.Sprintf("%s", msg)
		for _, gid := range needSendGroup {
			cm := model.CommonMsg{
				GroupId: gid,
				Message: send,
			}
			qqServe.SendMsg(cm, 2, qqclient)
			time.Sleep(200 * time.Millisecond)
		}
	}
	//fmt.Println(resp.String())
}

func StartWeibo() {
	f, _ := os.Open("./uid.txt")
	reader := bufio.NewReader(f)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			logrus.Error(err)
			break
		}
		uid := string(line)
		url := "https://m.weibo.cn/api/container/getIndex"
		client := req.C()
		resp, err := client.R().SetQueryParams(map[string]string{
			"type":  "uid",
			"value": uid,
		}).Get(url)
		//username := GetUserName(resp.String())
		container := GetContainerId(resp.String())
		//fmt.Println(container)
		//fmt.Printf("用户名为：%s", username)
		wbData, err := client.R().SetQueryParams(map[string]string{
			"type":        "uid",
			"value":       uid,
			"containerid": container,
		}).Get(url)
		t := gjson.Get(wbData.String(), "data.cards")
		for _, v := range t.Array() {
			if v.Get("card_type").Int() == 9 && !strings.Contains(v.Get("profile_type_id").Str, "top") {
				id := v.Get("mblog").Get("id").Str
				WeiBoId[uid] = id
				break
			}
		}
	}
	fmt.Println(666)
	fmt.Println(WeiBoId)

	f2, _ := os.Open("./uid.txt")
	reader2 := bufio.NewReader(f2)
	for {
		line, _, err := reader2.ReadLine()
		if err != nil {
			logrus.Error(err)
			break
		}
		go func(uid string) {
			for {
				GetWeibo(uid)
				time.Sleep(10 * time.Second)
			}
		}(string(line))
	}
}
