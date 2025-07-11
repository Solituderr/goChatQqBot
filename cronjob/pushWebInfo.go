package cronjob

import (
	"go-svc-tpl/app"
	"go-svc-tpl/client/http"
	"go-svc-tpl/logs"
	"go-svc-tpl/model"
	"go-svc-tpl/utils"
	"regexp"
	"strings"
	"time"

	"github.com/asmcos/requests"
)

var latestUrls []string

var needSendGroup = []string{"724938478"}

func PushWebInfo() {
	qqclient := app.GetQQClient()
	qqServe := utils.GetQQServer()
	latestUrls = extractUrl()
	for {
		time.Sleep(60 * time.Second)
		nowUrls := extractUrl()
		// 如果当前不一样，说明有更新
		nPtr := 0
		for nPtr < len(nowUrls) {
			if len(latestUrls) == 0 || nowUrls[nPtr] != latestUrls[0] {
				nPtr++
			}
		}
		for _, url := range nowUrls[:nPtr] {
			for _, gid := range needSendGroup {
				cm := model.CommonMsg{
					GroupId: gid,
					Message: "天天宝推送：all黄tag有新内容!\n" + url,
				}
				qqServe.SendMsg(cm, 2, qqclient)
				time.Sleep(200 * time.Millisecond)
			}
		}
		latestUrls = nowUrls
	}
}

func extractUrl() []string {
	headers := requests.Header{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":           "zh-CN,zh;q=0.9,en;q=0.8",
		"cache-control":             "max-age=0",
		"priority":                  "u=0, i",
		"referer":                   "https://xiangcaoweilanmeiqishui.lofter.com/post/1fc57442_2bed87bcb?sharefrom=LOFTER-Android 8.2.18&shareto=qq",
		"sec-ch-ua":                 `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "Windows",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		"Cookie":                    `usertrack=c+53SWhg4n2yk3x2CoGfAg==; JSESSIONID-WLF-XXD=0a45958493767851152e9fe317e687319eb89008ecae8665c84d85f484493e9f3fda97bb245b57f6df5bdc1d156742cb0f9cb3a3f2cd711b98cd524c5177710e2b2d789349a07b743d64b085b7a807b97154ea78d311ab5badec5613a45c0edade82a5f70c27632dfb0397af9ab34d74362151fea287f45df38740e7cf3c74e7ff677365; Hm_lvt_c5c55f9c94fbca8efd7d891afb3210e8=1751179949; LOFTER_SESS=BgxElu_2XGebDicZeA1bomKq9Ac8a-cyVJClCVDQmRM-OcPzAjhmoFnYMq1RdDovwGqkRlz5-UPZXnOQ10HnMvjHD-abxi1_VrGX6hpWEgIRywGI77RPfkxE_FWfDnHL; __LOFTER_TRACE_UID=AE096E1145FE4EE5A132CEBC5C224A45#2335448209#3; firstentry=%2Fpost.do%3FloftBlogName%3Dxiangcaoweilanmeiqishui%26loftPostUrl%3D1fc57442_2bed87bcb%26sharefrom%253DLOFTER-Android%25208.2.18%2526shareto%253Dqq|; regtoken=2000; hb_MA-BFD7-963BF6846668_source=xiangcaoweilanmeiqishui.lofter.com; reglogin_isLoginFlag=1; reglogin_isLoginFlag=1; NTESwebSI=07A54C5019D0154549BAA9BBC3E606FB.lofter-webapp-web-old-docker-lftpro-3-3nhsm-7sly4-5f89ffcfkrl49-8080`,
	}
	html := http.HttpService.HttpGetRequest("https://xiangcaoweilanmeiqishui.lofter.com/", headers)
	re := regexp.MustCompile(`<a href="(https?://[^\s"']+)">查看全文<\/a>`)

	// 查找所有匹配项
	matches := re.FindAllString(html, -1)
	// 取最新的数据
	if len(matches) > 0 {
		urls := []string{}
		for _, match := range matches {
			url := strings.ReplaceAll(match, `<a href="`, "")
			url = strings.ReplaceAll(url, `">查看全文</a>`, "")
			urls = append(urls, url)
		}
		logs.Info("[PushWebInfo] latestUrl: %v", urls)
		return urls
	}
	return nil
}
