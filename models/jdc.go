package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
	"os"
	"time"
)

var SMSCodes map[string]chan string

type jdcLogin struct {
	Ckcount int    `json:"ckcount"`
	Qlkey   string `json:"qlkey"`
	Phone   string `json:"Phone"`
}
type SendSMSResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Status   int `json:"status"`
		Captcha  int `json:"captcha"`
		Ckcount  int `json:"ckcount"`
		Tabcount int `json:"tabcount"`
	} `json:"data"`
}
type JdcResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Status int `json:"status"`
	} `json:"data"`
}
type VerifyCodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Qlid      string `json:"qlid"`
		Nickname  string `json:"nickname"`
		Timestamp string `json:"timestamp"`
		Remarks   string `json:"remarks"`
		Qlkey     int    `json:"qlkey"`
		Ckcount   int    `json:"ckcount"`
		Tabcount  int    `json:"tabcount"`
	} `json:"data"`
}

func JdcConfig() {
	req := httplib.Get(fmt.Sprintf("%s/api/Config", Config.JDCAddress))
	data, _ := req.String()
	logs.Info(data)
	//sender.Reply("获取容器NvJDC配置成功")
}

func JdcSendSMS(sender *Sender, phone string) error {
	req := httplib.Post(fmt.Sprintf("%s/api/SendSMS", Config.JDCAddress))
	req.Header("Content-Type", "application/json")
	body, _ := json.Marshal(struct {
		Phone string `json:"Phone"`
		Qlkey int    `json:"qlkey"`
	}{
		Phone: phone,
		Qlkey: 1,
	})
	req.Body(body)
	rsp, err := req.Response()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(rsp.Body)
	obj := SendSMSResponse{}
	err = json.Unmarshal(data, &obj)
	if err == nil {
		if obj.Success {
			sender.Reply("请输入验证码")
		} else {
			if obj.Message != "" {
				sender.Reply(obj.Message)
				logs.Info(obj.Message)
			}
			if obj.Data.Status == 666 {
				if obj.Data.Captcha == 2 {
					sender.Reply("出现文字验证码请前往手动验证：" + Config.JDCAddress + "/Captcha/" + phone)
				} else {
					JdcAutoCaptcha(sender, phone, 1)
				}
			}
		}
	} else {
		sender.Reply("登录出现错误")
	}
	return nil
}

func JdcAutoCaptcha(sender *Sender, phone string, number int) {
	sender.Reply(fmt.Sprintf("第%d次,尝试中******", number))
	req := httplib.Post(fmt.Sprintf("%s/api/AutoCaptcha", Config.JDCAddress))
	req.Header("Content-Type", "application/json")
	body, _ := json.Marshal(struct {
		Phone string `json:"Phone"`
	}{
		Phone: phone,
	})
	req.Body(body)
	rsp, err := req.Response()
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(rsp.Body)
	obj := JdcResponse{}
	err = json.Unmarshal(data, &obj)
	if err == nil {
		if obj.Success || obj.Data.Status != 666 {
			sender.Reply("安全认证破解成功 请输入验证码______")
		} else if !obj.Success && number <= 5 {
			logs.Info("验证失败")
			time.Sleep(time.Millisecond * 1000)
			JdcAutoCaptcha(sender, phone, number+1)
		} else {
			sender.Reply("安全认证破解失败，请联系管理员或者前往网页手动滑块")
		}
	}
}

func JdcVerifyCode(phone string, code string, qq string) bool {
	req := httplib.Post(fmt.Sprintf("%s/api/VerifyCode", Config.JDCAddress))
	req.Header("Content-Type", "application/json")
	body, _ := json.Marshal(struct {
		Phone string `json:"Phone"`
		Qlkey int    `json:"qlkey"`
		Code  string `json:"Code"`
		QQ    string `json:"QQ"`
	}{
		Phone: phone,
		Qlkey: 0,
		Code:  code,
		QQ:    qq,
	})
	req.Body(body)
	rsp, err := req.Response()
	if err != nil {
		return false
	}
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return false
	}
	obj := VerifyCodeResponse{}
	err = json.Unmarshal(data, &obj)
	if err == nil {
		if obj.Success {
			logs.Info(obj.Message)
		}
	}
	return obj.Success
}
func GetInput() string {
	//使用os.Stdin开启输入流
	//函数原型 func NewReader(rd io.Reader) *Reader
	//NewReader创建一个具有默认大小缓冲、从r读取的*Reader 结构见官方文档
	in := bufio.NewReader(os.Stdin)
	//in.ReadLine函数具有三个返回值 []byte bool error
	//分别为读取到的信息 是否数据太长导致缓冲区溢出 是否读取失败
	str, _, err := in.ReadLine()
	if err != nil {
		return err.Error()
	}
	return string(str)
}
