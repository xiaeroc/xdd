package models

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
	"time"
)

type MYAPP struct {
	Appid       int
	QVersion    string
	CountryCode int
}

func getMyApp() MYAPP {
	return MYAPP{
		Appid:       959,
		QVersion:    "1.0.0",
		CountryCode: 86,
	}
}

type SmsQuick struct {
	Data struct {
		Autologin  int    `json:"autologin"`
		Gsalt      string `json:"gsalt"`
		Guid       string `json:"guid"`
		Lsid       string `json:"lsid"`
		NeedAuth   int    `json:"need_auth"`
		ReturnPage string `json:"return_page"`
		RsaModulus string `json:"rsa_modulus"`
	} `json:"data"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}
type SmSCodeQuick struct {
	Data struct {
		ExpireTime int    `json:"expire_time"`
		Guid       string `json:"guid"`
		Lsid       string `json:"lsid"`
	} `json:"data"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}
type SmSCK struct {
	Data struct {
		Autologin   int    `json:"autologin"`
		ExpireTime  int    `json:"expire_time"`
		Guid        string `json:"guid"`
		Lsid        string `json:"lsid"`
		PtKey       string `json:"pt_key"`
		PtPin       string `json:"pt_pin"`
		RefreshTime int    `json:"refresh_time"`
		ReturnPage  string `json:"return_page"`
	} `json:"data"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

// SendSMS 发送短信接口
func SendSMS(sender *Sender, phone string) SmsQuick {
	quick, err := getQuick()
	ck := fmt.Sprintf("guid=%s;  lsid=%s;  gsalt=%s;  rsa_modulus=%s;", quick.Data.Guid, quick.Data.Lsid, quick.Data.Gsalt, quick.Data.RsaModulus)
	if err != nil {
		sender.Reply("发送验证码失败! 请联系管理员")
		return quick
	}
	app := getMyApp()
	ts := time.Now().Unix()
	subCmd := 2
	gsign := md5V(fmt.Sprintf(`%d%s%d36%d%s`, app.Appid, app.QVersion, ts, subCmd, quick.Data.Gsalt))
	sign := md5V(fmt.Sprintf("%d%s%d%s4dtyyzKF3w6o54fJZnmeW3bVHl0$PbXj", app.Appid, app.QVersion, app.CountryCode, phone))
	body := fmt.Sprintf("country_code=%d&client_ver=1.0.0&gsign=%s&appid=%d&mobile=%s&sign=%s&cmd=36&sub_cmd=%d&qversion=%s&ts=%d", app.CountryCode, gsign, app.Appid, phone, sign, subCmd, app.QVersion, ts)
	req := httplib.Post(fmt.Sprintf("https://qapplogin.m.jd.com/cgi-bin/qapp/quick"))
	req.Header("Host", "qapplogin.m.jd.com")
	req.Header("cookie", ck)
	req.Header("user-agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({\"packageName\":\"com.vivo.hybrid\",\"type\":\"deeplink\",\"extra\":{}})")
	req.Header("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header("content-type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header("content-length", string(len(body)))
	req.Header("accept-encoding", "")
	req.Body(body)
	rsp, err := req.Response()
	logs.Info(rsp.Body)
	smsCodeQuick := SmSCodeQuick{}
	if err != nil {
		sender.Reply("发送验证码失败! 请联系管理员")
		return quick
	}
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		sender.Reply("发送验证码失败! 请联系管理员")
		return quick
	}
	err = json.Unmarshal(data, &smsCodeQuick)
	if smsCodeQuick.ErrCode == 0 {
		sender.Reply("验证码发送成功! ")
		return quick
	} else {
		sender.Reply("发送验证码失败! 请联系管理员")
		return quick
	}
	return quick
}
func getQuick() (SmsQuick, error) {
	app := getMyApp()
	ts := time.Now().Unix()
	subCmd := 1
	sprintf := fmt.Sprintf(`%d%s%d36%dsb2cwlYyaCSN1KUv5RHG3tmqxfEb8NKN`, app.Appid, app.QVersion, ts, subCmd)
	gsign := md5V(sprintf)
	returnPage := "https%3A%2F%2Fcrpl.jd.com%2Fn%2Fmine%3FpartnerId%3DWBTF0KYY%26ADTAG%3Dkyy_mrqd%26token%3D"
	body := fmt.Sprintf("client_ver=1.0.0&gsign=%s&appid=%d&return_page=%s&cmd=36&sdk_ver=1.0.0&sub_cmd=%d&qversion=%s&ts=%d", gsign, app.Appid, returnPage, subCmd, app.QVersion, ts)
	req := httplib.Post(fmt.Sprintf("https://qapplogin.m.jd.com/cgi-bin/qapp/quick"))
	req.Header("Host", "qapplogin.m.jd.com")
	req.Header("cookie", "")
	req.Header("user-agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({\"packageName\":\"com.vivo.hybrid\",\"type\":\"deeplink\",\"extra\":{}})")
	req.Header("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header("content-type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header("content-length", string(len(body)))
	req.Header("accept-encoding", "")
	req.Body(body)
	rsp, err := req.Response()
	smsQuick := SmsQuick{}
	if err != nil {
		return smsQuick, err
	}
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return smsQuick, err
	}
	err = json.Unmarshal(data, &smsQuick)
	if err != nil || smsQuick.ErrCode != 0 {
		return smsQuick, err
	}
	return smsQuick, nil
}

func VerifyCode(phone string, code string, quick SmsQuick) SmSCK {
	ck := fmt.Sprintf("guid=%s;  lsid=%s;  gsalt=%s;  rsa_modulus=%s;", quick.Data.Guid, quick.Data.Lsid, quick.Data.Gsalt, quick.Data.RsaModulus)
	app := getMyApp()
	ts := time.Now().Unix()
	subCmd := 3
	gsign := md5V(fmt.Sprintf(`%d%s%d36%d%s`, app.Appid, app.QVersion, ts, subCmd, quick.Data.Gsalt))
	body := fmt.Sprintf("country_code=%d&client_ver=1.0.0&gsign=%s&smscode=%s&appid=%d&mobile=%s&cmd=36&sub_cmd=%d&qversion=%s&ts=%d", app.CountryCode, gsign, code, app.Appid, phone, subCmd, app.QVersion, ts)
	req := httplib.Post(fmt.Sprintf("https://qapplogin.m.jd.com/cgi-bin/qapp/quick"))
	req.Header("Host", "qapplogin.m.jd.com")
	req.Header("cookie", ck)
	req.Header("user-agent", "Mozilla/5.0 (Linux; Android 10; V1838T Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.87 mobile Safari/537.36 hap/1.9/vivo com.vivo.hybrid/1.9.6.302 com.jd.crplandroidhap/1.0.3 ({\"packageName\":\"com.vivo.hybrid\",\"type\":\"deeplink\",\"extra\":{}})")
	req.Header("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header("content-type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header("content-length", string(len(body)))
	req.Header("accept-encoding", "")
	req.Body(body)
	rsp, _ := req.Response()
	logs.Info(rsp.Body)
	data, _ := ioutil.ReadAll(rsp.Body)
	smsCk := SmSCK{}
	_ = json.Unmarshal(data, &smsCk)
	return smsCk
}
