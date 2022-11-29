package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type CodeSignal struct {
	Command []string
	Admin   bool
	Handle  func(sender *Sender) interface{}
}

type Sender struct {
	UserID            int
	ChatID            int
	Type              string
	Contents          []string
	MessageID         int
	Username          string
	IsAdmin           bool
	ReplySenderUserID int
}

type QQuery struct {
	Code int `json:"code"`
	Data struct {
		LSid          string `json:"lSid"`
		QqLoginQrcode struct {
			Bytes string `json:"bytes"`
			Sig   string `json:"sig"`
		} `json:"qqLoginQrcode"`
		RedirectURL string `json:"redirectUrl"`
		State       string `json:"state"`
		TempCookie  string `json:"tempCookie"`
	} `json:"data"`
	Message string `json:"message"`
}

func (sender *Sender) Reply(msg string) {
	switch sender.Type {
	case "tg":
		SendTgMsg(sender.UserID, msg)
	case "tgg":
		SendTggMsg(sender.ChatID, sender.UserID, msg, sender.MessageID, sender.Username)
	case "qq":
		SendQQ(int64(sender.UserID), msg)
	case "qqg":
		SendQQGroup(int64(sender.ChatID), int64(sender.UserID), msg)
	}
}

func (sender *Sender) JoinContens() string {
	return strings.Join(sender.Contents, " ")
}

func (sender *Sender) IsQQ() bool {
	return strings.Contains(sender.Type, "qq")
}

func (sender *Sender) IsTG() bool {
	return strings.Contains(sender.Type, "tg")
}

func (sender *Sender) handleJdCookies(handle func(ck *JdCookie)) error {
	cks := GetJdCookies()
	a := sender.JoinContens()
	ok := false
	if !sender.IsAdmin || a == "" {
		for i := range cks {
			if strings.Contains(sender.Type, "qq") {
				if cks[i].QQ == sender.UserID {
					if !ok {
						ok = true
					}
					handle(&cks[i])
				}
			} else if strings.Contains(sender.Type, "tg") {
				if cks[i].Telegram == sender.UserID {
					if !ok {
						ok = true
					}
					handle(&cks[i])
				}
			}
		}
		if !ok {
			sender.Reply("你尚未绑定🐶东账号，请抓取CK(不会抓的私聊群主，wsKey上车请私聊群主)私聊发机器人后即可查询账户资产信息。 请前往 " + Config.JDCAddress + " 进行登录")
			return errors.New("你尚未绑定🐶东账号，请抓取CK(不会抓的私聊群主，wsKey上车请私聊群主)私聊发机器人后即可查询账户资产信息。 请前往 " + Config.JDCAddress + " 进行登录")
		}
	} else {
		cks = LimitJdCookie(cks, a)
		if len(cks) == 0 {
			sender.Reply("没有匹配的账号")
			return errors.New("没有匹配的账号")
		} else {
			for i := range cks {
				handle(&cks[i])
			}
		}
	}
	return nil
}

func (sender *Sender) handleTenRead(handle func(ck *TenRead)) error {
	if strings.Contains(sender.Type, "qq") {
		ck, _ := GetTenRead(sender.UserID)
		handle(ck)
	}
	return nil
}

var codeSignals = []CodeSignal{
	{
		Command: []string{"登录", "登陆", "短信登录", "账号登录"},
		Handle: func(s *Sender) interface{} {
			s.Reply(fmt.Sprintf("请输入手机号___________ 或者前往 %s 进行登录\n请私聊机器人ck进行登录,不会请找管理", Config.JDCAddress))
			return nil
		},
	},

	{
		Command: []string{"status", "状态"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			return Count()
		},
	},
	{
		Command: []string{"QQ扫码", "qq扫码"},
		Handle: func(sender *Sender) interface{} {
			rsp, err := httplib.Post("https://api.kukuqaq.com/jd/qrcode").Response()
			if err != nil {
				return nil
			}
			body, err1 := ioutil.ReadAll(rsp.Body)
			if err1 == nil {
				fmt.Println(string(body))
			}
			s := &QQuery{}
			if len(body) > 0 {
				json.Unmarshal(body, &s)
			}
			logs.Info(s.Data.QqLoginQrcode.Bytes)
			ddd, _ := base64.StdEncoding.DecodeString(s.Data.QqLoginQrcode.Bytes) //成图片文件并把文件写入到buffer
			err2 := ioutil.WriteFile("./output.jpg", ddd, 0666)                   //buffer输出到jpg文件中（不做处理，直接写到文件）
			if err2 != nil {
				logs.Error(err2)
			}
			//ddd, _ := base64.StdEncoding.DecodeString("data:image/png;base64,"+s.Data.QqLoginQrcode.Bytes)
			return "data:image/png;base64," + s.Data.QqLoginQrcode.Bytes
		},
	},
	{
		Command: []string{"qrcode", "扫码", "二维码", "scan"},
		Handle: func(sender *Sender) interface{} {
			url := fmt.Sprintf("http://127.0.0.1:%d/api/login/qrcode.png?tp=%s&uid=%d&gid=%d", web.BConfig.Listen.HTTPPort, sender.Type, sender.UserID, sender.ChatID)
			if sender.Type == "tgg" {
				url += fmt.Sprintf("&mid=%v&unm=%v", sender.MessageID, sender.Username)
			}
			rsp, err := httplib.Get(url).Response()
			if err != nil {
				return nil
			}
			return rsp
		},
	},
	{
		Command: []string{"升级", "更新", "update", "upgrade"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			if err := Update(sender); err != nil {
				return err.Error()
			}
			sender.Reply("小滴滴重启程序")
			Daemon()
			return nil
		},
	},
	{
		Command: []string{"重启", "reload", "restart", "reboot"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.Reply("小滴滴重启程序")
			Daemon()
			return nil
		},
	},
	{
		Command: []string{"get-ua", "ua"},
		Handle: func(sender *Sender) interface{} {
			if !sender.IsAdmin {
				coin := GetCoin(sender.UserID)
				if coin < 0 {
					return "许愿币不足以查看UserAgent。"
				}
				sender.Reply("查看一次扣1个许愿币。")
				RemCoin(sender.UserID, 1)
			}
			return ua
		},
	},
	{
		Command: []string{"set-ua"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			ctt := sender.JoinContens()
			db.Create(&UserAgent{Content: ctt})
			ua = ctt
			return "已更新User-Agent。"
		},
	},
	{
		Command: []string{"任务列表"},
		Admin:   true,
		Handle: func(_ *Sender) interface{} {
			rt := ""
			for i := range Config.Repos {
				for j := range Config.Repos[i].Task {
					rt += fmt.Sprintf("%s\t%s\n", Config.Repos[i].Task[j].Title, Config.Repos[i].Task[j].Cron)
				}
			}
			return rt
		},
	},
	{
		Command: []string{"查询", "query"},
		Handle: func(sender *Sender) interface{} {
			if !sender.IsAdmin && GetEnv("query") == False {
				str := GetEnv("queryMsg")
				sender.Reply(str)
			} else {
				if sender.IsAdmin || getLimit(sender.UserID, 1) {
					sender.handleJdCookies(func(ck *JdCookie) {
						query := ck.Query()
						if sender.IsAdmin {
							query = query + fmt.Sprintf("\n优先级：%v", ck.Priority)
							query = query + fmt.Sprintf("\n绑定QQ：%v", ck.QQ)
						}
						sender.Reply(query)
					})
				} else {
					sender.Reply(fmt.Sprintf("鉴于东哥对接口限流，为了不影响大家的任务正常运行，即日起每日限流%d次，已超过今日限制", Config.Lim))
				}

			}
			return nil
		},
	},
	{
		Command: []string{"详细查询", "query"},
		Handle: func(sender *Sender) interface{} {
			if !sender.IsAdmin && GetEnv("query") == False {
				str := GetEnv("queryMsg")
				sender.Reply(str)
			} else {
				if sender.IsAdmin || getLimit(sender.UserID, 1) {
					sender.handleJdCookies(func(ck *JdCookie) {
						sender.Reply(ck.Query1())
					})
				} else {
					sender.Reply(fmt.Sprintf("鉴于东哥对接口限流，为了不影响大家的任务正常运行，即日起每日限流%d次，已超过今日限制", Config.Lim))
				}
			}
			return nil
		},
	},
	{
		Command: []string{"编译", "build"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.Reply("小滴滴正在编译程序")
			_, err := exec.Command("sh", "-c", "cd "+ExecPath+" && go build -o "+pname).Output()
			if err != nil {
				return errors.New("小滴滴编译失败：" + err.Error())
			} else {
				sender.Reply("小滴滴编译成功")
			}
			return nil
		},
	},
	{
		Command: []string{"备注", "bz"},
		Handle: func(sender *Sender) interface{} {
			if len(sender.Contents) > 1 {
				note := sender.Contents[0]
				sender.Contents = sender.Contents[1:]
				str := sender.Contents[0]
				number, err := strconv.Atoi(str)
				count := 0
				sender.handleJdCookies(func(ck *JdCookie) {
					count++
					if (err == nil && number == count) || ck.PtPin == str || sender.IsAdmin {
						ck.Update("Note", note)
						sender.Reply(fmt.Sprintf("已设置账号%s(%s)的备注为%s。", ck.PtPin, ck.Nickname, note))
					}
				})
			}
			return nil
		},
	},
	{
		Command: []string{"发送", "通知", "notify", "send"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			if len(sender.Contents) < 2 {
				sender.Reply("发送指令格式错误")
			} else {
				rt := strings.Join(sender.Contents[1:], " ")
				sender.Contents = sender.Contents[0:1]
				if sender.handleJdCookies(func(ck *JdCookie) {
					ck.Push(rt)
				}) == nil {
					return "操作成功"
				}
			}
			return nil
		},
	},
	{
		Command: []string{"run", "执行", "运行"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			name := sender.Contents[0]
			pins := ""
			if len(sender.Contents) > 1 {
				sender.Contents = sender.Contents[1:]
				err := sender.handleJdCookies(func(ck *JdCookie) {
					pins += "&" + ck.PtPin
				})
				if err != nil {
					return nil
				}
			}
			envs := []Env{}
			if pins != "" {
				envs = append(envs, Env{
					Name:  "pins",
					Value: pins,
				})
			}
			runTask(&Task{Path: name, Envs: envs}, sender)
			return nil
		},
	},
	{
		Command: []string{"upck", "刷新ck"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			updateCookie()
			return nil
		},
	},
	{
		Command: []string{"dlWskey", "dlwskey", "删除wskey"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			a := sender.JoinContens()
			if a == "" {
				sender.Reply(fmt.Sprintf("请指定要删除的账号"))
				return nil
			}
			sender.handleJdCookies(func(ck *JdCookie) {
				ck.Update("Wskey", "")
				sender.Reply(fmt.Sprintf("已清理WsKey%s ", ck.Nickname))
			})
			return nil
		},
	},
	{
		Command: []string{"qq", "QQ", "绑定qq"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			qqNum := Int(sender.Contents[0])
			if len(sender.Contents) > 1 {
				sender.Contents = sender.Contents[1:]
				sender.handleJdCookies(func(ck *JdCookie) {
					ck.Update(QQ, qqNum)
				})
			}
			if qqNum > 0 {
				return "绑定成功"
			} else {
				return "解绑成功"
			}
		},
	},
	{
		Command: []string{"cq", "CQ"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			str := ""
			sender.Contents = sender.Contents[0:]
			sender.handleJdCookies(func(ck *JdCookie) {
				str = str + fmt.Sprintf("账号：%s (%s) QQ：%d \n", ck.Nickname, ck.PtPin, ck.QQ)
			})
			return str
		},
	},
	{
		Command: []string{"优先级", "priority"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			priority := Int(sender.Contents[0])
			if len(sender.Contents) > 1 {
				sender.Contents = sender.Contents[1:]
				sender.handleJdCookies(func(ck *JdCookie) {
					ck.Update(Priority, priority)
					sender.Reply(fmt.Sprintf("已设置账号%s(%s)的优先级为%d。", ck.PtPin, ck.Nickname, priority))
				})
			}
			return nil
		},
	},
	{
		Command: []string{"cmd", "command", "命令"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			ct := sender.JoinContens()
			if regexp.MustCompile(`rm\s+-rf`).FindString(ct) != "" {
				return "over"
			}
			cmd(ct, sender)
			return nil
		},
	},
	{
		Command: []string{"环境变量", "environments", "envs"},
		Admin:   true,
		Handle: func(_ *Sender) interface{} {
			rt := []string{}
			envs := GetEnvs()
			if len(envs) == 0 {
				return "未设置任何环境变量"
			}
			for _, env := range envs {
				rt = append(rt, fmt.Sprintf(`%s="%s"`, env.Name, env.Value))
			}
			return strings.Join(rt, "\n")
		},
	},
	{
		Command: []string{"get-env", "env", "e"},
		Handle: func(sender *Sender) interface{} {
			ct := sender.JoinContens()
			if ct == "" {
				return "未指定变量名"
			}
			value := GetEnv(ct)
			if value == "" {
				return "未设置环境变量"
			}
			return fmt.Sprintf("环境变量的值为：" + value)
		},
	},
	{
		Command: []string{"set-env", "se", "export"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			env := &Env{}
			if len(sender.Contents) >= 2 {
				env.Name = sender.Contents[0]
				env.Value = strings.Join(sender.Contents[1:], " ")
			} else if len(sender.Contents) == 1 {
				ss := regexp.MustCompile(`^([^'"=]+)=['"]?([^=]+?)['"]?$`).FindStringSubmatch(sender.Contents[0])
				if len(ss) != 3 {
					return "无法解析"
				}
				env.Name = ss[1]
				env.Value = ss[2]
			} else {
				return "???"
			}
			ExportEnv(env)
			return "操作成功"
		},
	},
	{
		Command: []string{"unset-env", "ue", "unexport", "de"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			UnExportEnv(&Env{
				Name: sender.JoinContens(),
			})
			return "操作成功"
		},
	},
	{
		Command: []string{"reply", "回复"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			if len(sender.Contents) >= 2 {
				replies[sender.Contents[0]] = strings.Join(sender.Contents[1:], " ")
			} else {
				return "操作失败"
			}
			return "操作成功"
		},
	},
	{
		Command: []string{"help", "助力"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.handleJdCookies(func(ck *JdCookie) {
				ck.Update(Help, True)
				sender.Reply(fmt.Sprintf("已设置助力账号%s(%s)", ck.PtPin, ck.Nickname))
			})
			return nil
		},
	},
	{
		Command: []string{"tool", "工具人", "unhelp", "取消助力"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.handleJdCookies(func(ck *JdCookie) {
				ck.Update(Help, False)
				sender.Reply(fmt.Sprintf("已设置取消助力账号%s(%s)", ck.PtPin, ck.Nickname))
			})
			return nil
		},
	},
	{
		Command: []string{"屏蔽", "hack"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.handleJdCookies(func(ck *JdCookie) {
				ck.Update(Hack, True)
				sender.Reply(fmt.Sprintf("已设置屏蔽助力账号%s(%s)", ck.PtPin, ck.Nickname))
			})
			return nil
		},
	},
	{
		Command: []string{"取消屏蔽", "unhack"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.handleJdCookies(func(ck *JdCookie) {
				ck.Update(Hack, False)
				sender.Reply(fmt.Sprintf("已设置取消屏蔽助力账号%s(%s)", ck.PtPin, ck.Nickname))
			})
			return nil
		},
	},
	{
		Command: []string{"献祭", "导出"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.handleJdCookies(func(ck *JdCookie) {
				sender.Reply(fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin))
			})
			return nil
		},
	},
	{
		Command: []string{"删除账号", "删除", "清理过期"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			sender.Reply(fmt.Sprintf("PtKey为空并且是false的账号"))
			msg := "已清理账号"
			for _, ck := range GetJdCookies() {
				if ck.PtKey == "" && ck.Wskey == "" {
					ck.Removes(ck)
					if ck.Nickname == "" {
						msg += "\n" + ck.PtPin
					} else {
						msg += "\n" + ck.Nickname
					}
				}
			}
			return msg
		},
	},
	{
		Command: []string{"delete", "dl"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			a := sender.JoinContens()
			if a == "" {
				sender.Reply(fmt.Sprintf("请指定要删除的账号"))
				return nil
			}
			sender.handleJdCookies(func(ck *JdCookie) {
				ck.Removes(ck)
				sender.Reply(fmt.Sprintf("已清理账号%s", ck.Nickname))
			})
			return nil
		},
	},
	{
		Command: []string{"口令", "kl"},
		Handle: func(sender *Sender) interface{} {
			code := sender.JoinContens()
			command := JCommand(code)
			if command != "" {
				return command
			}
			return nil
		},
	},
	{
		Command: []string{"设置管理员"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			ctt := sender.JoinContens()
			db.Create(&UserAdmin{Content: ctt})
			return "已设置管理员"
		},
	},
	{
		Command: []string{"取消管理员"},
		Admin:   true,
		Handle: func(sender *Sender) interface{} {
			ctt := sender.JoinContens()
			RemoveUserAdmin(ctt)
			return "已取消管理员"
		},
	},
}

var mx = map[int]bool{}

func LimitJdCookie(cks []JdCookie, a string) []JdCookie {
	ncks := []JdCookie{}
	if s := strings.Split(a, "-"); len(s) == 2 {
		for i := range cks {
			if i+1 >= Int(s[0]) && i+1 <= Int(s[1]) {
				ncks = append(ncks, cks[i])
			}
		}
	} else if x := regexp.MustCompile(`^[\s\d,]+$`).FindString(a); x != "" {
		xx := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(a, -1)
		for i := range cks {
			for _, x := range xx {
				if fmt.Sprint(i+1) == x[1] {
					ncks = append(ncks, cks[i])
				} else if strconv.Itoa(cks[i].QQ) == x[1] {
					ncks = append(ncks, cks[i])
				}
			}

		}
	} else if a != "" {
		a = strings.Replace(a, " ", "", -1)
		for i := range cks {
			if strings.Contains(cks[i].Note, a) || strings.Contains(cks[i].Nickname, a) || strings.Contains(cks[i].PtPin, a) {
				ncks = append(ncks, cks[i])
			}
		}
	}
	return ncks
}

func ReturnCoin(sender *Sender) {
	tx := db.Begin()
	ws := []Wish{}
	if err := tx.Where("status = 0 and user_number = ?", sender.UserID).Find(&ws).Error; err != nil {
		tx.Rollback()
		sender.Reply(err.Error())
	}
	for _, w := range ws {
		if tx.Model(User{}).Where("number = ? ", sender.UserID).Update(
			"coin", gorm.Expr(fmt.Sprintf("coin + %d", w.Coin)),
		).RowsAffected == 0 {
			tx.Rollback()
			sender.Reply("愿望未达成退还许愿币失败。")
			return
		}
		sender.Reply(fmt.Sprintf("愿望未达成退还%d枚许愿币。", w.Coin))
		if tx.Model(&w).Update(
			"status", 1,
		).RowsAffected == 0 {
			tx.Rollback()
			sender.Reply("愿望未达成退还许愿币失败。")
			return
		}
	}
	tx.Commit()
}
