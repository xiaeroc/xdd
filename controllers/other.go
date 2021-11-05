package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiaeroc/xdd/models"
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
	if c.Ctx.Request.Method == "POST" || c.Ctx.Request.Method == "PUT" {
		requestBody := c.Ctx.Input.RequestBody
		var dataArr = []interface{}{}
		json.Unmarshal(requestBody, &dataArr)
		if len(dataArr) <= 0 {
			c.Ctx.WriteString("{\"code\":200,\"data\":[]}")
			return
		}
		name, nameErr := dataArr[0].(map[string]interface{})["name"].(string)
		value, valueErr := dataArr[0].(map[string]interface{})["value"].(string)
		remarks, _ := dataArr[0].(map[string]interface{})["remarks"].(string)
		if !nameErr || !valueErr {
			c.Ctx.WriteString("{\"code\":200,\"data\":[]}")
			return
		}
		if models.GetEnv("webSend") == models.True {
			go models.SendQQ(models.Config.QQID, fmt.Sprintf("JDC：添加，%s", value))
		}
		if value != "" {
			ptKey := FetchJdCookieValue("pt_key", value)
			ptPin := FetchJdCookieValue("pt_pin", value)
			if ptKey != "" && ptPin != "" {
				ck := models.JdCookie{
					PtKey: ptKey,
					PtPin: ptPin,
					QQ:    0,
					Hack:  models.False,
					Note:  remarks,
				}
				if models.CookieOK(&ck) {
					if !models.HasKey(ck.PtKey) {
						if nck, err := models.GetJdCookie(ck.PtPin); err == nil {
							nck.InPool(ck.PtKey)
							msg := fmt.Sprintf("更新账号，%s", ck.PtPin)
							logs.Info(msg)
							if models.GetEnv("webSend") == models.True {
								go models.SendQQ(models.Config.QQID, fmt.Sprintf("APP：更新账号，%s", ck.PtPin))
							}
						} else {
							models.NewJdCookie2(&ck)
							msg := fmt.Sprintf("添加账号，%s", ck.PtPin)
							logs.Info(msg)
							if models.GetEnv("webSend") == models.True {
								go models.SendQQ(models.Config.QQID, fmt.Sprintf("APP：添加账号，%s", ck.PtPin))
							}
						}
					}
				}
			}
		}
		c.Ctx.WriteString(fmt.Sprintf("{\"code\":200,\"data\":[{\"value\":\"%s\",\"_id\":\"0\",\"created\":0,\"status\":0,\"timestamp\":\"0\",\"position\":0,\"name\":\"%s\"}]}", value, name))
	} else {
		c.Ctx.WriteString("{\"code\":200,\"data\":[]}")
	}
}
