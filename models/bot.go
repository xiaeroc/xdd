package models

import (
	"fmt"
	"github.com/Mrs4s/go-cqhttp/coolq"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
)

var SendQQ = func(a int64, b interface{}) {
	coolq.SendQQ(a, b, Config.QQGroupID)
}
var SendQQGroup = func(a int64, b int64, c interface{}) {
	coolq.SendQQGroup(a, b, c)
}
var ListenQQPrivateMessage = func(uid int64, msg string) {
	SendQQ(uid, handleMessage(msg, "qq", int(uid)))
}

var ListenQQGroupMessage = func(gid int64, uid int64, msg string) {
	logs.Info(Config.QQGroupIDS)
	logs.Info(gid)
	logs.Info(strings.Contains(Config.QQGroupIDS, fmt.Sprintf("%d", gid)))
	if gid == Config.QQGroupID || strings.Contains(Config.QQGroupIDS, fmt.Sprintf("%d", gid)) {
		if Config.QbotPublicMode {
			SendQQGroup(gid, uid, handleMessage(msg, "qqg", int(uid), int(gid)))
		} else {
			SendQQ(uid, handleMessage(msg, "qq", int(uid)))
		}
	}

}

var replies = map[string]string{}

func InitReplies() {
	f, err := os.Open(ExecPath + "/conf/reply.php")
	if err == nil {
		defer f.Close()
		data, _ := ioutil.ReadAll(f)
		ss := regexp.MustCompile("`([^`]+)`\\s*=>\\s*`([^`]+)`").FindAllStringSubmatch(string(data), -1)
		for _, s := range ss {
			replies[s[1]] = s[2]
		}
	}
	if _, ok := replies["壁纸"]; !ok {
		replies["壁纸"] = "https://acg.toubiec.cn/random.php"
	}
}

var handleMessage = func(msgs ...interface{}) interface{} {
	msg := msgs[0].(string)
	args := strings.Split(msg, " ")
	head := args[0]
	contents := args[1:]
	sender := &Sender{
		UserID:   msgs[2].(int),
		Type:     msgs[1].(string),
		Contents: contents,
	}
	if len(msgs) >= 4 {
		sender.ChatID = msgs[3].(int)
	}
	if sender.Type == "tgg" {
		sender.MessageID = msgs[4].(int)
		sender.Username = msgs[5].(string)
		sender.ReplySenderUserID = msgs[6].(int)
	}
	if sender.UserID == Config.TelegramUserID || sender.UserID == int(Config.QQID) {
		sender.IsAdmin = true
	}
	if sender.IsAdmin == false {
		if IsUserAdmin(strconv.Itoa(sender.UserID)) {
			sender.IsAdmin = true
		}
	}
	for i := range codeSignals {
		for j := range codeSignals[i].Command {
			if codeSignals[i].Command[j] == head {
				return func() interface{} {
					if codeSignals[i].Admin && !sender.IsAdmin {
						return "你没有权限操作"
					}
					return codeSignals[i].Handle(sender)
				}()
			}
		}
	}
	switch msg {
	default:
		{
			ss := regexp.MustCompile(`^(\d{11})$`).FindStringSubmatch(msg)
			if len(ss) > 0 {
				if num := 5; len(codes) >= num {
					return fmt.Sprintf("%v坑位全部在使用中，请排队。", num)
				}
				id := "qq" + strconv.Itoa(sender.UserID)
				if _, ok := codes[id]; ok {
					return "你已在登录中。"
				}
				go func() {
					c := make(chan string, 1)
					codes = make(map[string]chan string)
					codes[id] = c
					defer delete(codes, id)
					phone := ss[0]
					logs.Info(phone)
					sender.Reply("请稍后，正在模拟环境...")
					quick := SendSMS(sender, phone)
					sms_code := ""
					select {
					case sms_code = <-c:
						sender.Reply("正在提交验证码...")
						smsCk := VerifyCode(phone, sms_code, quick)
						if smsCk.ErrCode == 0 {
							ck := JdCookie{
								PtKey: smsCk.Data.PtKey,
								PtPin: smsCk.Data.PtPin,
								Hack:  False,
								QQ:    sender.UserID,
							}
							if CookieOK(&ck) {
								if sender.IsQQ() {
									ck.QQ = sender.UserID
								} else if sender.IsTG() {
									ck.Telegram = sender.UserID
								}
								if nck, err := GetJdCookie(ck.PtPin); err == nil {
									if nck.QQ == 0 {
										nck.InPoolQQ(ck.PtKey, sender.UserID)
									} else {
										nck.InPool(ck.PtKey)
									}
									sender.Reply(fmt.Sprintf("更新账号，%s", ck.PtPin))
								} else {
									if Cdle {
										ck.Hack = True
									}
									NewJdCookie(&ck)
									sender.Reply(fmt.Sprintf("添加账号，%s", ck.PtPin))
								}
								for i := range Config.Containers {
									(&Config.Containers[i]).Write([]JdCookie{ck})
								}
							}
							sender.Reply("登录成功...")
						} else {
							sender.Reply("登录失败...")
						}
					case <-time.After(60 * time.Second):
						sender.Reply("验证码超时。")
						return

					}
					time.Sleep(time.Second)
				}()
			}

		}
		{
			ss := regexp.MustCompile(`^(\d{6})$`).FindStringSubmatch(msg)
			if len(ss) > 0 {
				if Config.JDCAddress == "" {
				}
				if code, ok := codes["qq"+fmt.Sprint(sender.UserID)]; ok {
					code <- ss[0]
					logs.Info(code)
				} else {
					sender.Reply("验证码不存在或过期了，请重新登录。")
				}
			}
		}
		{ //tyt https://pushgold.jd.com/#/helper?packetId=6d899aec8f52481197bc544be93ef69b-amRfNzIyY2M1MDJlOTM4Nw!!&currentActId=d5a8c7198ee54de093d2adb04089d3ec&utm_user=plusmember&ad_od=share&utm_source=androidapp&utm_medium=appshare&utm_campaign=t_335139774&utm_term=QQfriends
			ss := regexp.MustCompile(`activityId=(\S+)(&|&amp;)packetId=(\S+)(&|&amp;)currentActId`).FindStringSubmatch(msg)
			if len(ss) > 0 {
				if !sender.IsAdmin {
					coin := GetCoin(sender.UserID)
					if coin < 88 {
						return "推一推需要88个许愿币。"
					}
					RemCoin(sender.UserID, 88)
					sender.Reply("推一推即将开始，已扣除88个许愿币。")
				}
				runTask(&Task{Path: "jd_tyt.js", Envs: []Env{
					{Name: "actId", Value: ss[1]},
					{Name: "tytpacketId", Value: ss[3]},
				}}, sender)
				return nil
			}
		}
		{ //挖宝
			ss := regexp.MustCompile(`activityId=(\S+)(&|&amp;)inviterId=(\S+)(&|&amp;)inviterCode=(\w+)(&|&amp;)`).FindStringSubmatch(msg)
			if len(ss) > 6 {
				if !sender.IsAdmin {
					coin := GetCoin(sender.UserID)
					if coin < 188 {
						return "发财挖宝需要88个许愿币。"
					}
					RemCoin(sender.UserID, 188)
					sender.Reply("发财挖宝即将开始，已扣除88个许愿币。")
				}
				wbhelp := GetEnv("wbhelp")
				if wbhelp == "" {
					wbhelp = "0"
				}
				wbhelpMax := GetEnv("wbhelpMax")
				if wbhelpMax == "" {
					wbhelpMax = "30"
				}
				runTask(&Task{Path: "jd_fcwb.js", Envs: []Env{
					{Name: "activityId", Value: ss[1]}, {Name: "inviterId", Value: ss[3]}, {Name: "inviteCode", Value: ss[5]}, {Name: "wbhelpMax", Value: wbhelpMax}, {Name: "wbhelp", Value: wbhelp},
				}}, sender)
				return nil
			}
		}
		{
			ss := regexp.MustCompile(`activityId=(\S+)(&|&amp;)redEnvelopeId=(\S+)(&|&amp;)inviterId=(\S+)(&|&amp;)helpType=2`).FindStringSubmatch(msg)
			if len(ss) >= 7 {
				if sender.IsAdmin {
					sender.Reply("极速版大赢家提现即将开始。")
					dyjtx := Dyjtx{linkId: ss[1], redEnvelopeId: ss[3], inviter: ss[5]}
					Dyj_tx(dyjtx, sender)
				}
				return nil
			}
		}
		{ //
			ss := regexp.MustCompile(`pt_key=([^;=\s]+);[ ]*pt_pin=([^;=\s]+)`).FindAllStringSubmatch(msg, -1)

			if len(ss) > 0 {

				xyb := 0
				for _, s := range ss {
					ck := JdCookie{
						PtKey:    s[1],
						PtPin:    s[2],
						Priority: 2,
					}
					if CookieOK(&ck) {
						xyb++
						if sender.IsQQ() {
							ck.QQ = sender.UserID
						} else if sender.IsTG() {
							ck.Telegram = sender.UserID
						}
						if HasKey(ck.PtKey) {
							sender.Reply(fmt.Sprintf("重复提交"))
						} else {
							if nck, err := GetJdCookie(ck.PtPin); err == nil {
								if nck.QQ == 0 {
									nck.InPoolQQ(ck.PtKey, sender.UserID)
									SendQQ(Config.QQID, fmt.Sprintf("更新账号，%s，%d", ck.PtPin, sender.UserID))
								} else {
									nck.InPool(ck.PtKey)
									SendQQ(Config.QQID, fmt.Sprintf("更新账号，%s", ck.PtPin))
								}
								sender.Reply(fmt.Sprintf("更新账号，%s", ck.PtPin))
							} else {
								if Cdle {
									ck.Hack = True
								}
								NewJdCookie(&ck)
								sender.Reply(fmt.Sprintf("添加账号，%s", ck.PtPin))
							}
							for i := range Config.Containers {
								(&Config.Containers[i]).Write([]JdCookie{ck})
							}
						}
					} else {
						sender.Reply(fmt.Sprintf("无效账号，%s", ck.PtPin))
					}
				}
				go func() {
					Save <- &JdCookie{}
				}()
				return nil
			} else {
				ptPin := FetchJdCookieValue("pt_pin", msg)
				ptKey := FetchJdCookieValue("pt_key", msg)
				if ptPin != "" && ptKey != "" {
					ck := JdCookie{
						PtKey:    ptKey,
						PtPin:    ptPin,
						Priority: 2,
					}
					if CookieOK(&ck) {
						if sender.IsQQ() {
							ck.QQ = sender.UserID
						} else if sender.IsTG() {
							ck.Telegram = sender.UserID
						}
						if HasKey(ck.PtKey) {
							sender.Reply(fmt.Sprintf("重复提交"))
						} else {
							if nck, err := GetJdCookie(ck.PtPin); err == nil {
								if nck.QQ == 0 {
									nck.InPoolQQ(ck.PtKey, sender.UserID)
									SendQQ(Config.QQID, fmt.Sprintf("更新账号，%s，%d", ck.PtPin, sender.UserID))
								} else {
									nck.InPool(ck.PtKey)
									SendQQ(Config.QQID, fmt.Sprintf("更新账号，%s，%d", ck.PtPin))
								}
								msg := fmt.Sprintf("更新账号，%s", ck.PtPin)
								sender.Reply(fmt.Sprintf("更新账号，%s", ck.PtPin))
								if !sender.IsAdmin {
									SendQQ(Config.QQID, fmt.Sprintf("更新账号，%s", ck.PtPin))
								}
								(&JdCookie{}).Push(msg)
								logs.Info(msg)
							} else {
								if Cdle {
									ck.Hack = True
								}
								NewJdCookie(&ck)
								msg := fmt.Sprintf("添加账号，%s", ck.PtPin)
								sender.Reply(fmt.Sprintf("添加账号，%s", ck.PtPin))
								(&JdCookie{}).Push(msg)
								logs.Info(msg)
							}
							for i := range Config.Containers {
								(&Config.Containers[i]).Write([]JdCookie{ck})
							}
						}
					} else {
						sender.Reply(fmt.Sprintf("无效账号，%s", ck.PtPin))
					}
				}
			}
		}
		{
			o := findShareCode(msg)
			if o != "" {
				return "导入互助码成功"
			}
		}
		{
			ss := regexp.MustCompile(`pin=([^;=\s]+);[ ]*wskey=([^;=\s]+)`).FindAllStringSubmatch(msg, -1)
			if len(ss) > 0 {
				for _, s := range ss {
					if HasWsKey(msg) {
						sender.Reply(fmt.Sprintf("重复提交"))
						continue
					}
					if fleas, str := WsKeyOK2(&JdCookie{
						Wskey: s[0],
						PtPin: s[1],
					}); fleas {
						ck := JdCookie{
							Wskey: msg,
							PtPin: FetchJdCookieValue("pt_pin", str),
							PtKey: FetchJdCookieValue("pt_key", str),
						}
						if sender.IsQQ() {
							ck.QQ = sender.UserID
						} else if sender.IsTG() {
							ck.Telegram = sender.UserID
						}
						if nck, err := GetJdCookie(ck.PtPin); err == nil {
							nck.InPoolWsKey(ck.PtKey, ck.Wskey)
							msg := fmt.Sprintf("更新账号，%s", ck.PtPin)
							sender.Reply(fmt.Sprintf("更新账号，%s", ck.PtPin))
							logs.Info(msg)
						} else {
							if Cdle {
								ck.Hack = True
							}
							NewJdCookieWsKey(&ck)
							msg := fmt.Sprintf("添加账号，%s", ck.PtPin)
							sender.Reply(fmt.Sprintf("添加账号，%s", ck.PtPin))
							(&JdCookie{}).Push(msg)
							logs.Info(msg)
						}
						for i := range Config.Containers {
							(&Config.Containers[i]).Write([]JdCookie{ck})
						}
					} else {
						sender.Reply(fmt.Sprintf("无效账号，%s", s[1]))
					}
				}
				go func() {
					Save <- &JdCookie{}
				}()
				return nil
			}
		}
		{
			PHPSESSID := FetchJdCookieValue("PHPSESSID", msg)
			udtauth := FetchJdCookieValue("udtauth", msg)
			if PHPSESSID != "" && udtauth != "" {
				tr := TenRead{
					CK:   "PHPSESSID=" + PHPSESSID + "; udtauth=" + udtauth + ";",
					SSID: PHPSESSID,
					QQ:   sender.UserID,
				}
				if nck, err := GetTenRead(tr.QQ); err == nil {
					nck.UpdateTenRead(TenReadCK, nck.CK+"@"+tr.CK)
				} else {
					NewTenRead(&tr)
				}
				sender.Reply(fmt.Sprintf("添加账号10秒阅读账号，%s", &tr.SSID))
				logs.Info(msg)
				return nil
			}

		}
		for k, v := range replies {
			if regexp.MustCompile(k).FindString(msg) != "" {
				if strings.Contains(msg, "妹") && time.Now().Unix()%10 == 0 {
					v = "https://pics4.baidu.com/feed/d833c895d143ad4bfee5f874cfdcbfa9a60f069b.jpeg?token=8a8a0e1e20d4626cd31c0b838d9e4c1a"
				}
				if regexp.MustCompile(`^https{0,1}://[^\x{4e00}-\x{9fa5}\n\r\s]{3,}$`).FindString(v) != "" {
					url := v
					rsp, err := httplib.Get(url).Response()
					if err != nil {
						return nil
					}
					ctp := rsp.Header.Get("content-type")
					if ctp == "" {
						rsp.Header.Get("Content-Type")
					}
					if strings.Contains(ctp, "text") || strings.Contains(ctp, "json") {
						data, _ := ioutil.ReadAll(rsp.Body)
						return string(data)
					}
					return rsp
				}
				return v
			}
		}
	}
	for _, function := range functions {
		for _, rule := range function.Rules {
			var matched bool
			str := ""
			if function.FindAll {
				if res := regexp.MustCompile(rule).FindAllStringSubmatch(msg, -1); len(res) > 0 {
					tmp := [][]string{}
					for i := range res {
						tmp = append(tmp, res[i][1:])
					}
					matched = true
				}
			} else {
				if res := regexp.MustCompile(rule).FindStringSubmatch(msg); len(res) > 0 {
					//sender.SetMatch(res[1:])
					// https://item.jd.com/" + id + ".html"
					str = res[len(res)-1]
					matched = true
				}
			}
			if matched {
				sender.Reply(JdPriceFunc(str))
				return ""
				//rt := function.Handle(sender)
				//if rt != nil {
				//	sender.Reply("")
				//}
				//return ""
			}
		}
	}
	return nil
}

func FetchJdCookieValue(key string, cookies string) string {
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}
