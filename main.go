package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	SignInURL          = "https://www.hifini.com/sg_sign.htm"
	CookieEnvVariable  = "COOKIE"
	PayloadEnvVariable = "PAYLOAD"
	DingDingWebhook    = "DINGDING_WEBHOOK"
)

type DingDingMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	success := SignIn(client)
	if success {
		result := "签到成功"
		fmt.Println(result)
		dingding(result)
	} else {
		result := "签到失败"
		fmt.Println(result)
		dingding(result)
		os.Exit(3)
	}
}

// SignIn 签到
func SignIn(client *http.Client) bool {
	cookie := os.Getenv(CookieEnvVariable)
	if cookie == "" {
		log.Println("COOKIE不存在，请检查是否添加")
		return false
	}
	payload_str := os.Getenv(PayloadEnvVariable)
	log.Println(payload_str)
	payload := strings.NewReader(payload_str)

	req, err := http.NewRequest("POST", SignInURL, payload)
	if err != nil {
		log.Println("创建请求失败:", err)
		return false
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("发送请求失败:", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取响应失败:", err)
		return false
	}

	log.Println(string(body))

	return strings.Contains(string(body), "成功")
}

func dingding(result string) {
	message := DingDingMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: "HiFiNi" + result,
		},
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Println("转换消息为JSON失败:", err)
		return
	}

	webhook := os.Getenv(DingDingWebhook)
	if webhook == "" {
		log.Println("DINGDING_WEBHOOK不存在，请检查是否添加")
		return
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		log.Println("发送消息失败:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("消息发送成功")
}
