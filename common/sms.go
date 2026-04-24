package common

import (
	"encoding/json"

	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// 发送短信
func SendSMS(to string, body string) bool {
	twilio_sid := viper.GetString("twilio_sid")
	twilio_token := viper.GetString("twilio_token")
	twilio_from := viper.GetString("twilio_from")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilio_sid,
		Password: twilio_token,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(twilio_from)
	params.SetBody(body)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		Log.Error("发送短信失败", err)
		return false
	}
	response, _ := json.Marshal(*resp) //将数据编码为json字符串, 返回[]uint8
	//返回json字符串写入日志
	Log.Info(string(response))

	return true
}