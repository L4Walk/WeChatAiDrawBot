/*
package main

import (

	"WeChatAiDrawBot/bootstrap"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"

)

	func main() {
		bootstrap.Run()
	}
*/
package main

import (
	"WeChatAiDrawBot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type openAiMessage struct {
	role    string `json:"role"`
	content string `json:"content"`
}

type openAiBody struct {
	model       string        `json:"model"`
	message     openAiMessage `json:"message"`
	temperature float32       `json:"temperature"`
}

type Payload struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int64   `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int64   `json:"top_p"`
	FrequencyPenalty int64   `json:"frequency_penalty"`
	PresencePenalty  int64   `json:"presence_penalty"`
	Model            string  `json:"model"`
}

func p2() {
	data := Payload{
		Prompt:           "腾讯是一家怎样的公司",
		MaxTokens:        2048,
		Temperature:      0.5,
		TopP:             0,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Model:            "gpt-3.5-turbo",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	//req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", body)
	req, err := http.NewRequest("POST", config.LoadConfig().BASE_URL+"chat/completions", body)
	//req, err := http.NewRequest("POST", "https://openai.api2d.net/v1/chat/completions", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", config.LoadConfig().OPENAI_API_KEY))
	//req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", "fk200248-aQOaMhfKtRxkAEYIc9VNzbkMN3HOtPWV"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(all))
}

func post() {
	url := config.LoadConfig().BASE_URL + "chat/completions"
	//body := strings.NewReader("Hello Wolrd")
	msg := openAiMessage{role: "user", content: "Say this is a test!"}
	bdy := openAiBody{model: "gpt-3.5-turbo", message: msg, temperature: 0.7}
	jsonValue, err := json.Marshal(bdy)
	//jsonStr := []bytes(`{"model":"gpt-3.5-turbo","message"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.LoadConfig().OPENAI_API_KEY)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	bodyString := string(bodyBytes)

	fmt.Print(bodyString)
}

func main() {
	p2()
}

/*
func main() {
	client := openai.NewClient("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")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
*/
