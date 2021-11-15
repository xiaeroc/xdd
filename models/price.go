package models

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type JdPrice struct {
	PriceStatus int `json:"price_status"`
	Store       []struct {
		FirstPrice       float64   `json:"first_price"`
		SecondPrice      float64   `json:"second_price"`
		Lowest           float64   `json:"lowest"`
		Highest          float64   `json:"highest"`
		LastPrice        float64   `json:"last_price"`
		MinStamp         string    `json:"min_stamp"`
		MaxStamp         int64     `json:"max_stamp"`
		AllLine          []float64 `json:"all_line"`
		AllLineBeginTime int64     `json:"all_line_begin_time"`
		HalfYearLineTime int64     `json:"half_year_line_time"`
		HalfYearLine     []float64 `json:"half_year_line"`
		MonthLineTime    int64     `json:"month_line_time"`
		MonthLine        []float64 `json:"month_line"`
		LowestDate       int       `json:"lowest_date"`
		Name             string    `json:"name"`
		PriceRange       string    `json:"price_range"`
		CurrentPrice     string    `json:"current_price"`
		PriceStatus      int       `json:"price_status"`
	} `json:"store"`
	Analysis struct {
		Tip       string `json:"tip"`
		PromoDays []struct {
			Show  string  `json:"show"`
			Price float64 `json:"price"`
			Date  string  `json:"date"`
		} `json:"promo_days"`
	} `json:"analysis"`
	Promo []struct {
		Price    int `json:"price"`
		OriPrice int `json:"ori_price"`
		Msg      struct {
			Promotion string `json:"promotion,omitempty"`
			Coupon    string `json:"coupon,omitempty"`
		} `json:"msg"`
		Time int `json:"time"`
		S    int `json:"_s,omitempty"`
	} `json:"promo"`
	NopuzzlePromo []struct {
		Price    int `json:"price"`
		OriPrice int `json:"ori_price"`
		Msg      struct {
			Promotion string `json:"promotion,omitempty"`
			Coupon    string `json:"coupon,omitempty"`
		} `json:"msg"`
		Time int `json:"time"`
		S    int `json:"_s,omitempty"`
	} `json:"nopuzzle_promo"`
	NowDay  int64 `json:"now_day"`
	NowHour int64 `json:"now_hour"`
	ItemQr  struct {
		Type string `json:"type"`
		Src  string `json:"src"`
	} `json:"itemQr"`
}

func JdPriceFunc(str string) string {
	materialUrl := ""
	if IsNum(str) {
		materialUrl = fmt.Sprintf("https://item.jd.com/%s.html", str)
	} else {
		materialUrl = fmt.Sprintf("https://u.jd.com/%s", str)
	}
	rebateLink := GetRebateLink(materialUrl)
	req := httplib.Get(fmt.Sprintf("https://browser.bijiago.com/extension/price_towards?dp_ids=undefined&dp_id=%s-3&ver=1&format=jsonp&union=union_bijiago&version=1594190525099&from_device=bijiago&from_type=bjg_ser&crc64=1&_=%d", str, time.Now().Unix()))
	req.Header("Connection", "keep-alive")
	req.Header("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header("sec-ch-ua-mobile", "?0")
	req.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header("Accept", "*/*")
	req.Header("Sec-Fetch-Site", "cross-site")
	req.Header("Sec-Fetch-Mode", "no-cors")
	req.Header("Sec-Fetch-Dest", "script")
	req.Header("Referer", "https://item.jd.com/")
	req.Header("Accept-Language", "zh-CN,zh;q=0.9")
	req.SetTimeout(time.Second*3, time.Second*3)
	rsp, _ := req.Response()
	data, _ := ioutil.ReadAll(rsp.Body)
	price := JdPrice{}
	_ = json.Unmarshal(data, &price)
	logs.Info(req.String())
	fmt.Println()
	req.Response()
	var build strings.Builder
	if strs := strings.Split(rebateLink.Official, "\n"); len(strs) > 0 {
		for i, str := range strs {
			if i == 1 && price.Store[0].CurrentPrice != "" {
				s := price.Store[0]
				if len(price.Store) >= 2 {
					s = price.Store[1]
				}
				_startDate, _ := strconv.ParseInt(s.MinStamp, 10, 64)
				var ninStamp = time.Unix(_startDate, 0).Format("2006-01-02 15:04:05")
				tip := ""
				if price.Analysis.Tip != "" {
					tip = price.Analysis.Tip
				}
				build.WriteString(fmt.Sprintf("\n%s \n最高价：%.0f  \n最低价：%.0f  %s\n%s：%.0f %s\n%s：%.0f %s\n比价结果仅供参考\n",
					tip, s.Highest, s.Lowest, ninStamp[0:10],
					price.Analysis.PromoDays[0].Show, price.Analysis.PromoDays[0].Price, price.Analysis.PromoDays[0].Date,
					price.Analysis.PromoDays[1].Show, price.Analysis.PromoDays[1].Price, price.Analysis.PromoDays[1].Date))
			}
			if !strings.Contains(str, "佣金") && str != "" {
				build.WriteString("\n")
				build.WriteString(str)
			}
		}
	}
	return build.String()
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
