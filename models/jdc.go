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
		Ckcount  int `json:"ckcount"`
		Tabcount int `json:"tabcount"`
	} `json:"data"`
}
type JdcResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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

var JdcUrl = Config.

func JdcConfig() {
	req := httplib.Get(fmt.Sprintf("%s/api/Config", JdcUrl))
	data, _ := req.String()
	logs.Info(data)
	//sender.Reply("获取容器NvJDC配置成功")
}

func JdcSendSMS(sender *Sender, phone string) error {
	req := httplib.Post(fmt.Sprintf("%s/api/SendSMS", JdcUrl))
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
				logs.Info("安全验证")
				JdcAutoCaptcha(sender, phone, 1)
			}
		}

	} else {
		sender.Reply("登录出现错误")
	}
	return nil
}

func JdcAutoCaptcha(sender *Sender, phone string, number int) {
	req := httplib.Post(fmt.Sprintf("%s/api/AutoCaptcha", JdcUrl))
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
		if obj.Success {
			sender.Reply("安全认证破解成功 请输入短信验证码")
		} else if !obj.Success && number <= 5 {
			logs.Info("验证失败")
			time.Sleep(time.Millisecond * 1000)
			JdcAutoCaptcha(sender, phone, number+1)
		} else {
			sender.Reply("安全认证破解失败，请联系管理员")
		}
	}
}

func JdcVerifyCode(phone string, code string) string {
	req := httplib.Post(fmt.Sprintf("%s/api/VerifyCode", JdcUrl))
	req.Header("Content-Type", "application/json")
	body, _ := json.Marshal(struct {
		Phone string `json:"Phone"`
		Qlkey int    `json:"qlkey"`
		Code  string `json:"Code"`
	}{
		Phone: phone,
		Qlkey: 1,
		Code:  code,
	})
	req.Body(body)
	rsp, err := req.Response()
	data, err := ioutil.ReadAll(rsp.Body)
	obj := VerifyCodeResponse{}
	err = json.Unmarshal(data, &obj)
	if err == nil {
		logs.Info(obj.Data)
		logs.Info(obj.Message)
		if obj.Success {
			logs.Info("登录成功")
			return obj.Data.Nickname
		}
	}
	return ""
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
