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
)

var latestUrls []string

var needSendGroup = []string{"658227422", "429093558", "724938478"}

func PushWebInfo() {
	qqclient := app.GetQQClient()
	qqServe := utils.GetQQServer()
	latestUrls = extractUrl()
	for {
		time.Sleep(60 * time.Second)
		nowUrls := extractUrl()
		if len(nowUrls) == 0 {
			continue
		}
		nowUrls = nowUrls[:1]
		// 如果当前不一样，说明有更新
		nPtr := 0
		for nPtr < len(nowUrls) {
			if len(latestUrls) == 0 {
				nPtr = len(nowUrls)
				break
			}
			if nowUrls[nPtr] == latestUrls[0] {
				break
			}
			nPtr++
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
	headers := map[string]string{
		"accept":                    "*/*",
		"accept-language":           "zh-CN,zh;q=0.9,en;q=0.8",
		"cache-control":             "max-age=0",
		"priority":                  "u=0, i",
		"referer":                   "https://www.lofter.com/tag/all%E9%BB%84",
		"sec-ch-ua":                 `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "Windows",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		"Cookie":                    `usertrack=c+53SWhg4n2yk3x2CoGfAg==; Hm_lvt_c5c55f9c94fbca8efd7d891afb3210e8=1751179949; hb_MA-BFD7-963BF6846668_source=xiangcaoweilanmeiqishui.lofter.com; NTESwebSI=92EFEF92C61052E02C2CA9781C269E36.lofter-webapp-web-old-docker-lftpro-3-3nhsm-czkrl-6c4478ffw9kc5-8080; firstentry=%2Fblogindex.do%3FloftBlogName%3Dxiangcaoweilanmeiqishui%26|; JSESSIONID-WLF-XXD=e94fda8f9d16d4a21c1fa6b29a289f4b537c4d725a6b520f543093a55815ebd328c127628ef0254153f03c7e63b635fc337a4c19a044e1e5bd44a687ae9ba96a356916f74d488a46ac4f4738a41e2644eb2387a03d46c62efe36768a559fea15930b450f3d41987464e5598eb8dfdb8ce27edf3cc46ddee0be6d9dc8143fbe5acbda7629; LOFTER-PHONE-LOGINNUM=15161159168; LOFTER-PHONE-LOGIN-FLAG=1; LOFTER-PHONE-LOGIN-AUTH=E6otP4BGrhAXrQrK7CxWhpg6vBrMbOxsCtMDHY1GvuaDXql0z0T9W7Ldtyx3pDoZQBSVIBe4Sf99mSj2nCUTHISp_-dvb9cQ; token=E6otP4BGrhAXrQrK7CxWhpg6vBrMbOxsCtMDHY1GvuaDXql0z0T9W7Ldtyx3pDoZQBSVIBe4Sf99mSj2nCUTHISp_-dvb9cQ; phone=15161159168; deviceid=eeef0f93-c131-47e5-8cb8-6db974e8722e; __LOFTER_TRACE_UID=BD5C4D72D4B24B83B1661B67612AE99C#2351499150#14; reglogin_isLoginFlag=1; reglogin_isLoginFlag=1`,
	}
	body := `callCount=1
scriptSessionId=${scriptSessionId}187
httpSessionId=
c0-scriptName=TagBean
c0-methodName=search
c0-id=0
c0-param0=string:all%E9%BB%84
c0-param1=number:0
c0-param2=string:
c0-param3=string:new
c0-param4=boolean:false
c0-param5=number:0
c0-param6=number:20
c0-param7=number:0
c0-param8=number:0
batchId=951941`
	html, err := http.HttpService.HttpPostRequest("https://www.lofter.com/dwr/call/plaincall/TagBean.search.dwr", headers, body)
	if err != nil {
		logs.Error("[extractUrl] err: %v", err)
		return nil
	}
	// logs.Info("html: %s", html)
	re := regexp.MustCompile(`https?://[a-zA-Z0-9-]+\.lofter\.com/post/[a-zA-Z0-9_]+`)

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
	logs.Error("[PushWebInfo] 无获取到指定信息")
	return nil
}
