package ChatGPT

import (
	"WeChatAiDrawBot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// const BASEURL = "https://api.openai.com/v1/"
// const BASEURL = "https://apiproxy.chuheng.tech/proxy/https://api.openai.com/v1/chat/completions" //cloudflare反代
// const BASEURL = "http://proxy-api.chuheng.tech/proxy/api.openai.com/v1/" // vercel反代
//const BASEURL = "https://openai.api2d.net/v1/chat/completions" //api2d接口

var BASEURL = config.LoadConfig().BASE_URL
var CHATURL = config.LoadConfig().CHATGPT_URL

type ChoiceItem struct {
}

// 响应体
type gptResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
		PreTokenCount    int `json:"pre_token_count"`
		PreTotal         int `json:"pre_total"`
		AdjustTotal      int `json:"adjust_total"`
		FinalTotal       int `json:"final_total"`
	} `json:"usage"`
}

type Choice struct {
	Index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

// API2D req模型
type ApifoxModel struct {
	MaxTokens        *int64    `json:"max_tokens,omitempty"` // 如果不指定，每次请求会预先冻结 模型支持的最大个数 Token，如果是 3.5 也就是 41P ，在请求完成后按 usage 字段多退少补。如果账户余额少于 41P; 就会报错，此时可以指定 max_tokens，这样就会按 max_tokens 来冻结
	Messages         []Message `json:"messages"`
	Model            Model     `json:"model"`
	Moderation       *bool     `json:"moderation,omitempty"`      // 默认为 false，为 true 时会调用文本安全接口对内容进行判定，并将审核结果添加到返回值中的 moderation; 字段，开发者可以根据值自行判断如何处理。审核输出的详细解释：https://cloud.tencent.com/document/product/1124/51860 开启后每; 9000 字符会增加 10P 的消耗
	ModerationStop   *bool     `json:"moderation_stop,omitempty"` // 默认为 false，在 moderation 为 true 且自身也为 true 时，如果审核结果不是 Pass，将自动进行内容拦截，对流也生效
	SafeMode         *bool     `json:"safe_mode,omitempty"`       // 默认为 false，为 true 时会尝试让 GPT 自己审查内容，不输出违规结果。由于 GPT; 的调性，效果好坏比较随机。总的来说对暴力、色情内容效果较好，政治类效果一般。开启后会每次访问会增加约 1P 的消耗
	Stream           *bool     `json:"stream,omitempty"`          // 流方式返回，兼容官方参数，但因为有计费和审核逻辑，比官方流慢。
	Temperature      float32   `json:"temperature"`               // 哪个采样温度，在 0和2之间。较高的值，如0.8会使输出更随机，而较低的值，如0.2会使其更加集中和确定性。
	TopP             int       `json:"top_p"`                     // 一种替代温度采样的方法叫做核心采样，模型会考虑到具有 top_p 概率质量的标记结果。因此，0.1 表示只有占前 10% 概率质量的标记被考虑
	FrequencyPenalty int       `json:"frequency_penalty"`         // 介于-2.0和2.0之间的数字。正值会根据文本中新令牌的现有频率对其进行惩罚，从而降低模型重复相同行的可能性。
	PresencePenalty  int       `json:"presence_penalty"`          // 介于 -2.0 和 2.0 之间的数字。正值会根据它们是否出现在文本中迄今为止来惩罚新令牌，从而增加模型谈论新主题的可能性。
}

// Message模型
//
//	Role可选 ['system', 'assistant', 'user', 'function']
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// 模型类型
type Model string

// 模型列表
const (
	GPT35Turbo        Model = "gpt-3.5-turbo"
	GPT35Turbo0301    Model = "gpt-3.5-turbo-0301"
	GPT35Turbo0613    Model = "gpt-3.5-turbo-0613"
	GPT35Turbo16K     Model = "gpt-3.5-turbo-16k"
	GPT35Turbo16K0613 Model = "gpt-3.5-turbo-16k-0613"
	GPT4              Model = "gpt-4"
	GPT40613          Model = "gpt-4-0613"
)

// Completions gtp文本模型回复
// curl https://api.openai.com/v1/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer your chatGPT key"
// -d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
	mes := []Message{
		{msg, "user"},
	}

	/*
		requestBody := ApifoxModel{
			Model:          GPT35Turbo,
			Messages:       mes,
			MaxTokens:      2048,
			Moderation:     false,
			ModerationStop: false,
			SafeMode:       false,
			Stream:         true,
		}*/

	requestBody := ApifoxModel{
		Model:       GPT35Turbo,
		Messages:    mes,
		Temperature: 0.7,
	}

	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request ChatGPT json string : %v", string(requestData))
	req, err := http.NewRequest("POST", CHATURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().OPENAI_API_KEY
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

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

	//gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))

	var res gptResponse
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", err
	}

	// 访问解析后的数据
	fmt.Println("ID:", res.ID)
	fmt.Println("Object:", res.Object)
	fmt.Println("Created:", res.Created)
	fmt.Println("Model:", res.Model)

	// 提取choices中的内容
	if len(res.Choices) > 0 {
		fmt.Println("Role:", res.Choices[0].Message.Role)
		fmt.Println("Content:", res.Choices[0].Message.Content)
		fmt.Println("Finish Reason:", res.Choices[0].FinishReason)
	}

	// 提取usage中的内容
	fmt.Println("Prompt Tokens:", res.Usage.PromptTokens)
	fmt.Println("Completion Tokens:", res.Usage.CompletionTokens)
	fmt.Println("Total Tokens:", res.Usage.TotalTokens)
	fmt.Println("Pre Token Count:", res.Usage.PreTokenCount)
	fmt.Println("Pre Total:", res.Usage.PreTotal)
	fmt.Println("Adjust Total:", res.Usage.AdjustTotal)
	fmt.Println("Final Total:", res.Usage.FinalTotal)

	reply := res.Choices[0].Message.Content

	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
