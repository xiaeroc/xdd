package main

import (
	"github.com/Mrs4s/go-cqhttp/cmd/gocq"
	"github.com/Mrs4s/go-cqhttp/coolq"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/beego/beego/v2/server/web"
	"github.com/xiaeroc/xdd/controllers"
	"github.com/xiaeroc/xdd/models"
)

var theme = ""

func main() {
	go func() {
		models.Save <- &models.JdCookie{}
	}()
	web.Get("/count", func(ctx *context.Context) {
		ctx.WriteString(models.Count())
	})
	web.Get("/", func(ctx *context.Context) {
		if models.Config.Theme == "" {
			models.Config.Theme = models.GhProxy + "https://raw.githubusercontent.com/xiaeroc/xdd/master/myTheme/kuduan.html"
		}
		if theme != "" {
			ctx.WriteString(theme)
			return
		}
		if strings.Contains(models.Config.Theme, "http") {
			logs.Info("下载最新主题")
			s, _ := httplib.Get(models.Config.Theme).String()
			if s != "" {
				theme = s
				ctx.WriteString(s)
				return
			}
			logs.Warn("主题下载失败，使用默认主题")
		}
		f, err := os.Open(models.Config.Theme)
		if err == nil {
			d, _ := ioutil.ReadAll(f)
			theme = string(d)
			ctx.WriteString(string(d))
			return
		}
	})
	web.Router("/api/login/qrcode", &controllers.LoginController{}, "get:GetQrcode")
	web.Router("/api/login/qrcode.png", &controllers.LoginController{}, "get:GetQrcode")
	web.Router("/api/login/query", &controllers.LoginController{}, "get:Query")
	web.Router("/api/login/cookie", &controllers.LoginController{}, "get:Cookie")
	web.Router("/api/account", &controllers.AccountController{}, "get:List")
	web.Router("/api/account", &controllers.AccountController{}, "post:CreateOrUpdate")
	web.Router("/admin", &controllers.AccountController{}, "get:Admin")
	web.Router("/admin", &controllers.AccountController{}, "get:Admin")
	web.Router("/api/login/ck", &controllers.LoginController{}, "POST:CkLogin")
	web.Router("/getSign", &controllers.LoginController{}, "get:GetSign")
	web.Router("/api/login/smslogin", &controllers.LoginController{}, "post:SMSLogin")
	web.Router("/api/appCkLogin", &controllers.LoginController{}, "put:AppCkLogin")
	web.Router("/api/appUpdate", &controllers.LoginController{}, "get:AppUpdate")
	web.Router("/open/auth/token", &controllers.OtherController{}, "*:AuthToken")
	web.Router("/open/envs", &controllers.OtherController{}, "*:Envs")
	if models.Config.Static == "" {
		models.Config.Static = "./static"
	}
	web.BConfig.WebConfig.StaticDir["/static"] = models.Config.Static
	web.BConfig.AppName = models.AppName
	web.BConfig.WebConfig.AutoRender = false
	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	web.BConfig.WebConfig.Session.SessionName = models.AppName
	go func() {
		time.Sleep(time.Second * 4)
		(&models.JdCookie{}).Push("小滴滴已启动")
	}()
	go gocq.Main()
	coolq.PrivateMessageEventCallback = models.ListenQQPrivateMessage
	coolq.GroupMessageEventCallback = models.ListenQQGroupMessage
	web.Run()
}
