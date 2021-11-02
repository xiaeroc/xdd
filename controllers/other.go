package controllers

import "github.com/beego/beego/v2/core/logs"

type OtherController struct {
	BaseController
}

func (c *OtherController) AuthToken() {
	logs.Info(c.Ctx.Request.GetBody)
	logs.Info(c.Ctx.Request.URL)
	logs.Info(c.Ctx.Request.PostForm)
}

func (c *OtherController) Envs() {
	logs.Info(c.Ctx.Request.URL)
	logs.Info(c.Ctx.Request.GetBody)
	logs.Info(c.Ctx.Request.PostForm)
}
