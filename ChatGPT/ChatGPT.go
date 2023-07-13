package ChatGPT

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// const BASEURL = "https://api.openai.com/v1/"
// const BASEURL = "https://apiproxy.chuheng.tech/proxy/https://api.openai.com/v1/chat/completions" //cloudflare反代
// const BASEURL = "http://proxy-api.chuheng.tech/proxy/api.openai.com/v1/" // vercel反代
const BASEURL = "https://openai.api2d.net/v1/chat/completions" //api2d接口

//var BASEURL = config.LoadConfig().BASE_URL

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ChoiceItem struct {
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

// Completions gtp文本模型回复
// curl https://api.openai.com/v1/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer your chatGPT key"
// -d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
	requestBody := ChatGPTRequestBody{
		Model:            "gpt-3.5-turbo",
		Prompt:           msg,
		MaxTokens:        2048,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request ChatGPT json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	//apiKey := config.LoadConfig().OPENAI_API_KEY
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+"fk200248-aQOaMhfKtRxkAEYIc9VNzbkMN3HOtPWV")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v["text"].(string)
			break
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
