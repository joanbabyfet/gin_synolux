package common

import "github.com/mojocn/base64Captcha"

// 验证码存放
var Store = base64Captcha.DefaultMemStore

// 获取图片验证码
func GetCaptcha() (string, string, string, error) {
	driver := base64Captcha.DefaultDriverDigit
	captcha := base64Captcha.NewCaptcha(driver, Store)
	id, b64s, answer, err := captcha.Generate()
	return id, b64s, answer, err
}