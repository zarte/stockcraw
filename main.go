package main

import (
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unsafe"
)

//定义结构体
type StackInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Val string `json:"val"`
	Per string `json:"per"`
	Ptype string `json:"ptype"`
	Daystr string `json:"daystr"`
}

func main()  {
	weburl :="https://www.hkexnews.hk/sdw/search/mutualmarket_c.aspx"
	data := make(map[string]string)
	header := make(map[string]string)
	res,err := curl(weburl,data,header,"GET")

	if err!=nil{
		fmt.Println(err)
	}else{
		//无异常
		//解析页面数据
		if res == ""{
			fmt.Println("页面内容为空")
			return
		}
		dealHtmlContent(res)
		
	}
}

func dealHtmlContent(content string)  {
	doc,_ := htmlquery.Parse(strings.NewReader(content))
	trlist := htmlquery.Find(doc,`//*[@id="mutualmarket-result"]/tbody/tr`)
	for _,tr  := range trlist {
		var stuckinfo StackInfo
		content := htmlquery.FindOne(tr,".//td[1]/div[2]")
		//股票编号
		stuckinfo.Code = htmlquery.InnerText(content)
		if stuckinfo.Code==""{
			fmt.Println("craw detail fail")
			continue
		}
		//名称
		content = htmlquery.FindOne(tr,".//td[2]/div[2]")
		stuckinfo.Name = htmlquery.InnerText(content)
		//fmt.Println(stuckinfo.Name)
		//stuckinfo.Name = strings.Replace(stuckinfo.Name, " ", "", -1)
		//fmt.Println(stuckinfo.Name)
		//持有量
		content = htmlquery.FindOne(tr,".//td[3]/div[2]")
		stuckinfo.Val = htmlquery.InnerText(content)
		stuckinfo.Val = strings.Replace(stuckinfo.Val, ",", "", -1)
		//百分比
		content = htmlquery.FindOne(tr,".//td[4]/div[2]")
		stuckinfo.Per = strings.Trim(htmlquery.InnerText(content),"%")
		stuckinfo.Ptype = "1"
		//Todo存入数据库
		fmt.Println(stuckinfo)
	}
}

/**
请求函数
curltype: post get json
*/
func curl(weburl string,data map[string]string,header map[string]string,curltype string) (string,error) {
	var ctype string
	reader := strings.NewReader("")
	switch strings.ToUpper(curltype) {
	case "JSON":
	case "POST":
		ctype = "POST"
		str, err := json.Marshal(data)
		if err != nil {
			fmt.Println("json.Marshal failed:", err)
			return "",err
		}
		reader = strings.NewReader(string(str))
		break
	case "GET":
		ctype = "GET"
		params := url.Values{}
		parseURL, err := url.Parse(weburl)
		if err != nil {
			log.Println("err")
		}
		for key,val := range data {
			params.Set(key, val)
		}
		//对中文进行url编码
		parseURL.RawQuery = params.Encode()
		weburl = parseURL.String()
		break
	default:
		return "",errors.New("未定义请求类型")
	}
	request, err := http.NewRequest(ctype, weburl, reader)
	if err != nil {
		return "",err
	}
	//添加头部参数
	for key,item := range header{
		request.Header.Set(key,item)
	}
	//正式发起请求
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "",err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	res := (*string)(unsafe.Pointer(&respBytes))
	return *res,nil
}