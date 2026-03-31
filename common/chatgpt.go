package common

import (
	"context"

	"github.com/lexkong/log"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

// 发送ChatGPT
func ChatGPT(keyword string) (bool, string) {
	api := viper.GetString("chat_gpt_api")
	client := openai.NewClient(api)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: keyword,
				},
			},
		},
	)
	if err != nil {
		log.Error("发送ChatGPT失败", err)
		return false, ""
	}
	result := resp.Choices[0].Message.Content
	return true, result
}