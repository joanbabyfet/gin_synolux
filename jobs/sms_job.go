// 队列任务
package jobs

import (
	"bytes"
	"encoding/gob"
	"gin-synolux/utils"

	"github.com/lexkong/log"
)

type SubscribeSMS struct {
	To   string
	Body string
}

// 发送邮件
func (s *SubscribeSMS) ActionSMS(args interface{}) error {
	var scb SubscribeSMS
	b := args.([]byte)
	decoder := gob.NewDecoder(bytes.NewReader(b))
	decoder.Decode(&scb)

	//发送
	ok := utils.SendSMS(scb.To, scb.Body)
	if !ok {
		log.Error("发送短信"+scb.To+"失败", nil)
	}
	return nil
}
