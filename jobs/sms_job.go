// 队列任务
package jobs

import (
	"bytes"
	"encoding/gob"
	"gin-synolux/common"
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
	ok := common.SendSMS(scb.To, scb.Body)
	if !ok {
		common.Log.Error("发送短信"+scb.To+"失败", nil)
	}
	return nil
}
