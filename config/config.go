package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Configuration 项目配置
type Configuration struct {
	// openai api key (required)
	OPENAI_API_KEY string `json:"OPENAI_API_KEY"`

	// You can start service behind a proxy
	PROXY_URL string `json:"PROXY_URL"`

	// Override openai api request base url. (optional)
	// Default: https://api.openai.com
	// Examples: http://your-openai-proxy.com
	BASE_URL string `json:"BASE_URL"`

	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}

		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("OPENAI_API_KEY")
		AutoPass := os.Getenv("AutoPass")
		if ApiKey != "" {
			config.OPENAI_API_KEY = ApiKey
		}
		if AutoPass == "true" {
			config.AutoPass = true
		}
	})

	//fmt.Print(config.OPENAI_API_KEY)
	//fmt.Print(config.BASE_URL)
	return config
}
