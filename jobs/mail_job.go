// 队列任务
package jobs

import (
	"bytes"
	"encoding/gob"
	"gin-synolux/utils"

	"github.com/lexkong/log"
)

type SubscribeMail struct {
	To      string
	Subject string
	Body    string
}

// 发送邮件
func (s *SubscribeMail) ActionMail(args interface{}) error {
	var scb SubscribeMail
	b := args.([]byte)
	decoder := gob.NewDecoder(bytes.NewReader(b))
	decoder.Decode(&scb)

	//发送
	ok := utils.SendMail(scb.To, scb.Subject, scb.Body)
	if !ok {
		log.Error("发送邮件"+scb.To+"失败", nil)
	}
	return nil
}
