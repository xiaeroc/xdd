package models

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
)

type RebateLink struct {
	Code     int      `json:"code"`
	Content  string   `json:"content"`
	Images   []string `json:"images"`
	Official string   `json:"official"`
}

func GetRebateLink(str string) RebateLink {
	req := httplib.Get(fmt.Sprintf("https://api.jingpinku.com/get_rebate_link/api?appid=%s&appkey=%s&union_id=%s&content=%s", Config.AppId, Config.AppKey, Config.UnionId, str))
	rsp, _ := req.Response()
	data, _ := ioutil.ReadAll(rsp.Body)
	so := RebateLink{}
	_ = json.Unmarshal(data, &so)
	logs.Info(req.String())
	return so
}
func GetRebateLinkStr(str string) string {
	link := GetRebateLink(str)
	msg := "暂无商品信息"
	if link.Content != "" {
		msg = link.Content
	}
	if link.Official != "" {
		msg = link.Official
	}
	return msg
}
