package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
)

//func (sender *Sender)Get_rebatel_link(str string) {
//
//	req := httplib.Get(fmt.Sprintf(`https://api.m.jd.com/?functionId=openRedEnvelopeInteract&body=%s&t=%d&appid=activities_platform&clientVersion=3.6.0`, body, time.Now().Unix()))
//
//	data, _ := req.String()
//}
var appid = "2110052102069332"
var unionId = "2011685381"
var appKey = "YLiFUk7uUqLElJ3N7MLypQj1d5PUK4Of"

type RebateLink struct {
}

func (sender *Sender) GetRebateLink(str string) {
	req := httplib.Get(fmt.Sprintf("https://api.jingpinku.com/get_rebate_link/api?appid=%s&appkey=%s&union_id=%s&content=%s", appid, appKey, unionId, str))
	data, _ := req.String()
	logs.Info(data)
}

// [rule: raw https://item\.m\.jd\.com/product/(\d+).html]
// [rule: raw https://item\.m\.jd\.com/(\d+).html]
// [rule: raw https://m\.jingxi\.com/item/jxview\?sku=(\d+)]
// [rule: raw https://kpl\.m\.jd\.com/product\?wareId=(\d+)]
// [rule: raw https://wq\.jd\.com/item/view\?sku=(\d+)]
//
//var id = param(1)
//
//var content = ""
//
//if (!isNaN(id) && (parseInt(id).toString().length === id.length)) {
//content = "https://item.jd.com/" + id + ".html"
//} else {
//content = "https://u.jd.com/" + id
//}
//
//var data = request({
//"url": "https://api.jingpinku.com/get_rebate_link/api?" +
//"appid=" + get("jingpinku_appid") +
//"&appkey=" + get("jingpinku_appkey") +
//"&union_id=" + get("jd_union_id") +
//"&content=" + content,
//"dataType": "json"
//})
//if (data && data.code == 0) {
//if (data.official) {
//if (data.images.length > 0) {
//sendImage(data.images[0])
//}
//var finals = [];
//var lines = data.official.split("\n");
//for (var i = 0; i < lines.length; i++) {
//if (lines[i].indexOf("佣金") == -1) {
//finals.push(lines[i])
//}
//}
//sendText(finals.join("\n"))
//} else {
//sendText("暂无商品信息。")
//}
//} else {
//sendText("异常。")
//}
