package models

import (
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	cron "github.com/robfig/cron/v3"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func init() {
	c = cron.New()
	c.Start()
}

type Function struct {
	Rules   []string
	FindAll bool
	Admin   bool
	Handle  func(sender *Sender) interface{}
}

var functions = []Function{}

func initFunction() {
	files, err := ioutil.ReadDir(ExecPath + "/replies")
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		jr := string(ExecPath + "/replies/" + v.Name())
		data := ""
		if strings.Contains(jr, "http") {
			data, err = httplib.Get(jr).String()
			if err != nil {
				logs.Warn("回复：%s获取失败%v", jr, err)
				continue
			}
		} else {
			f, err := os.Open(jr)
			if err != nil {
				logs.Warn("回复：%s打开失败%v", jr, err)
				continue
			}
			v, _ := ioutil.ReadAll(f)
			data = string(v)
		}
		rules := []string{}
		for _, res := range regexp.MustCompile(`\[rule:(.+)\]`).FindAllStringSubmatch(data, -1) {
			rules = append(rules, strings.Trim(res[1], " "))
		}
		admin := false
		if res := regexp.MustCompile(`\[admin:([^\[\]]+)\]`).FindStringSubmatch(data); len(res) != 0 {
			admin = strings.Trim(res[1], " ") == "true"
		}
		if len(rules) == 0 {
			logs.Warn("回复：%s找不到规则", jr, err)
			continue
		}
		var handler = func(s *Sender) interface{} {

			return "11111111"
		}
		logs.Warn("回复：%s添加成功", jr)
		AddCommand("", []Function{
			{
				Handle: handler,
				Rules:  rules,
				Admin:  admin,
			},
		})
	}
}

func AddCommand(prefix string, cmds []Function) {
	for j := range cmds {
		for i := range cmds[j].Rules {
			if strings.Contains(cmds[j].Rules[i], "raw ") {
				cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "raw ", "", -1)
				continue
			}
			if strings.Contains(cmds[j].Rules[i], "$") {
				continue
			}
			if prefix != "" {
				cmds[j].Rules[i] = prefix + `\s+` + cmds[j].Rules[i]
			}
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "(", `[(]`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], ")", `[)]`, -1)
			cmds[j].Rules[i] = regexp.MustCompile(`\?$`).ReplaceAllString(cmds[j].Rules[i], `(.+)`)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], " ", `\s+`, -1)
			cmds[j].Rules[i] = strings.Replace(cmds[j].Rules[i], "?", `(\S+)`, -1)
			cmds[j].Rules[i] = "^" + cmds[j].Rules[i] + "$"
		}
		functions = append(functions, cmds[j])
	}
}

func FetchCookieValue(ps ...string) string {
	var key, cookies string
	if len(ps) == 2 {
		if len(ps[0]) > len(ps[1]) {
			key, cookies = ps[1], ps[0]
		} else {
			key, cookies = ps[0], ps[1]
		}
	}
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}
