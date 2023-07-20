package main

import (
	"WeChatAiDrawBot/bootstrap"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	bootstrap.Run()
}

func p3() {

	url := "https://oa.api2d.net/v1/chat/completions"
	method := "POST"

	payload := strings.NewReader(`{
    "model": "gpt-3.5-turbo",
    "messages": [
        {
            "role": "user",
            "content": "讲个笑话"
        }
    ],
    "safe_mode": false
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer fk200248-AiV5NErCx6Y9kYWLRn8ZMrBqHXrH35iM|ck242-aa7b04a")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func p4() {

	url := "https://oa.api2d.net/dashboard/billing/credit_grants"
	method := "GET"

	payload := strings.NewReader(`{
    "model": "text-davinci-edit-001",
    "instruction": "请修改文本中的拼写错误",
    "input": "What tim is it"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer fk200248-AiV5NErCx6Y9kYWLRn8ZMrBqHXrH35iM|ck242-aa7b04a")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
