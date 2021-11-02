package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

type OtherController struct {
	BaseController
}

func (c *OtherController) AuthToken() {
	logs.Info(c.Ctx.Request.GetBody)
	logs.Info(c.Ctx.Request.URL)
	logs.Info(c.Ctx.Request.PostForm)
	c.Ctx.WriteString(fmt.Sprintf("{\"code\":200,\"data\":{\"token\":\"3e60261b-dece-4a1e-942c-690ba7205f76\",\"token_type\":\"Bearer\",\"expiration\":%d}}", time.Now().Unix()+2626560))
}

func (c *OtherController) Envs() {
	logs.Info(c.Ctx.Request.URL)
	logs.Info(c.Ctx.Request.PostForm)
	logs.Info(c.GetString("value"))
	logs.Info(c.GetString("name"))
	if c.Ctx.Request.Method == "POST" {
		body := c.Ctx.Request.PostForm
		value := body.Get("value")
		name := body.Get("name")
		c.Ctx.WriteString(fmt.Sprintf("{\"code\":200,\"data\":[{\"value\":\"%s\",\"_id\":\"0\",\"created\":0,\"status\":0,\"timestamp\":\"0\",\"position\":0,\"name\":\"%s\"}]}", value, name))
	} else if c.Ctx.Request.Method == "PUT" {
		c.Ctx.WriteString(fmt.Sprintf("{\"code\":200,\"data\":[{\"value\":\"%s\",\"_id\":\"0\",\"created\":0,\"status\":0,\"timestamp\":\"0\",\"position\":0,\"name\":\"%s\"}]}"))
	} else {
		c.Ctx.WriteString("{\"code\":200,\"data\":[]}")
	}
}
