package common

import (
	"crypto/tls"
	"strconv"

	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// 发送邮件
func SendMail(to string, subject string, body string) bool {
	mail_host := viper.GetString("mail_host")
	mail_port, _ := strconv.Atoi(viper.GetString("mail_port"))
	mail_username := viper.GetString("mail_username")
	mail_password := viper.GetString("mail_password")
	mail_from_address := viper.GetString("mail_from_address")
	mail_from_name := viper.GetString("mail_from_name")

	m := gomail.NewMessage()
	m.SetAddressHeader("From", mail_from_address, mail_from_name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(mail_host, mail_port, mail_username, mail_password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Error("发送邮件失败", err)
		//panic(err) //后续代码不会执行
		return false
	}
	return true
}