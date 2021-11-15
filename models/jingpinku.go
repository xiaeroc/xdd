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
type GoodsLink struct {
	ShortUrl string `json:"short_url"`
}
type HistoryPrice struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Price float64 `json:"price"`
		Date  string  `json:"date"`
	} `json:"data"`
}

func GetRebateLink(str string) RebateLink {
	req := httplib.Get(fmt.Sprintf("https://api.jingpinku.com/get_powerful_coup_link/api?appid=%s&appkey=%s&union_id=%s&content=%s", Config.AppId, Config.AppKey, Config.UnionId, str))
	rsp, _ := req.Response()
	data, _ := ioutil.ReadAll(rsp.Body)
	so := RebateLink{}
	_ = json.Unmarshal(data, &so)
	logs.Info(req.String())
	//if so.Content!="" {
	//	GetHistoryPrice(str)
	//}
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
func getGoodsLink(str string) string {
	req := httplib.Get(fmt.Sprintf("https://api.jingpinku.com/get_powerful_coup_link/api?appid=%s&appkey=%s&union_id=%s&material_url=%s", Config.AppId, Config.AppKey, Config.UnionId, str))
	rsp, _ := req.Response()
	data, _ := ioutil.ReadAll(rsp.Body)
	so := GoodsLink{}
	_ = json.Unmarshal(data, &so)
	logs.Info(req.String())
	return so.ShortUrl
}

func GetHistoryPrice(str string) string {
	req := httplib.Get(fmt.Sprintf("https://api.jingpinku.com/get_history_price/api?appid=%s&appkey=%s&sku_url=%s", Config.AppId, Config.AppKey, str))
	rsp, _ := req.Response()
	data, _ := ioutil.ReadAll(rsp.Body)
	so := HistoryPrice{}
	_ = json.Unmarshal(data, &so)
	logs.Info(req.String())
	return so.Code
}
