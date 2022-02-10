package models

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"strings"

	"time"
)

type Dyjtx struct {
	linkId        string
	redEnvelopeId string
	inviter       string
}

func Dyj_tx(tx Dyjtx, sender *Sender) {
	cks := GetJdCookies()
	body := "%7B%22linkId%22%3A%22" + tx.linkId + "%22%2C%22redEnvelopeId%22%3A%22" + tx.redEnvelopeId + "%22%2C%22inviter%22%3A%22" + tx.inviter + "%22%2C%22helpType%22%3A2%7D"
	str := ""
	for i := range cks {
		if cks[i].PtKey != "" {
			req := httplib.Get(fmt.Sprintf(`https://api.m.jd.com/?functionId=openRedEnvelopeInteract&body=%s&t=%d&appid=activities_platform&clientVersion=3.6.0`, body, time.Now().Unix()))
			req.Header("Host", "618redpacket.jd.com;")
			req.Header("Accept-Language", "zh-cn")
			req.Header("Accept-Encoding", "gzip, deflate, br")
			req.Header("Referer", "https://618redpacket.jd.com/")
			req.Header("User-Agent", fmt.Sprintf("jdpingou;iPhone;4.13.0;14.4.2;%d-%d;network/wifi;model/iPhone10,2;appBuild/100609;ADID/00000000-0000-0000-0000-000000000000;supportApplePay/1;hasUPPay/0;pushNoticeIsOpen/1;hasOCPay/0;supportBestPay/0;session/c%dk;pap/JA2019_3111789;brand/apple;supportJDSHWK/1;Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X)AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148", time.Now().Unix(), time.Now().Unix(), cks[i].ID))
			req.Header("Connection", "keep-alive;")
			req.Header("cookie", fmt.Sprintf("pt_key=%s;pt_pin=%s;", cks[i].PtKey, cks[i].PtPin))

			data, _ := req.String()

			if strings.Contains(data, "助力成功") {
				str = str + "助力成功\n"
			} else if strings.Contains(data, "活动太火爆了") {
				str = str + "活动太火爆了，请稍后重试\n"
			} else if strings.Contains(data, "这次运气不好") {
				str = str + "这次运气不好，试试拆自己的红包吧\n"
			} else if strings.Contains(data, "您今日已帮该好友提现资格助力过") {
				str = str + "您今日已帮该好友提现资格助力过\n"
			} else if strings.Contains(data, "你的好友今日已成功提现") {
				str = str + "你的好友今日已成功提现\n"
				sender.Reply(str)
				return
			}

		}
	}
	sender.Reply(str)
}

type JCommandDate struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Img      string `json:"img"`
		HeadImg  string `json:"headImg"`
		Title    string `json:"title"`
		UserName string `json:"userName"`
		JumpUrl  string `json:"jumpUrl"`
	} `json:"data"`
}

func JCommand(code string) string {
	req := httplib.Post(fmt.Sprintf(`https://api.jds.codes/jd/jcommand`))
	req.Header("content-type", "application/json")
	req.Body(fmt.Sprintf(`{"code":"%s"}`, code))
	data, err := req.Bytes()
	jCommandDate := JCommandDate{}
	err = json.Unmarshal(data, &jCommandDate)
	if err == nil && jCommandDate.Code == 200 {
		return fmt.Sprintf("活动: %s \n 用户: %s \n 地址: %s ", jCommandDate.Data.Title, jCommandDate.Data.UserName, jCommandDate.Data.JumpUrl)
	}
	return ""
}
