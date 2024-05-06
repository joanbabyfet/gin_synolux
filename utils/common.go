// 公共方法
package utils

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/lexkong/log"
	"github.com/mojocn/base64Captcha"
	"github.com/sashabaranov/go-openai"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// 验证码存放
var Store = base64Captcha.DefaultMemStore

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	//占位待%x为整型以十六进制方式显示
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// 密码加密
func PasswordHash(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// 密码验证
func PasswordVerify(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

// 生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

// 获取时间戳
func Timestamp() int {
	t := time.Now().Unix()
	return int(t)
}

// 获取当前日期时间
func DateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取当前日期
func Date() string {
	return time.Now().Format("2006-01-02")
}

// 时间戳转日期
func UnixToDateTime(timestramp int) string {
	t := time.Unix(int64(timestramp), 0)
	return t.Format("2006-01-02 15:04:05") //通用时间模板定义
}

// 时间戳转日期
func UnixToDate(timestramp int) string {
	t := time.Unix(int64(timestramp), 0)
	return t.Format("2006-01-02") //通用时间模板定义
}

// 日期转时间戳
func DateToUnix(str string) int {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		return 0
	}
	return int(t.Unix())
}

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
		log.Error("发送短信失败", err)
		return false
	}
	response, _ := json.Marshal(*resp) //将数据编码为json字符串, 返回[]uint8
	//返回json字符串写入日志
	log.Info(string(response))

	return true
}

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

// 获取图片完整地址
func DisplayImg(filename string) string {
	if filename == "" {
		return ""
	}
	return viper.GetString("file_url") + "/image/" + filename
}

// 获取视频完整地址
func DisplayVideo(filename string) string {
	if filename == "" {
		return ""
	}
	return viper.GetString("file_url") + "/video/" + filename
}

// 获取图片验证码
func GetCaptcha() (string, string, string, error) {
	driver := base64Captcha.DefaultDriverDigit
	captcha := base64Captcha.NewCaptcha(driver, Store)
	id, b64s, answer, err := captcha.Generate()
	return id, b64s, answer, err
}

// 加密
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 对明文进行 padding
	plaintext = pkcs7Padding(plaintext, block.BlockSize())

	// CBC 模式下的加密对象
	iv := []byte("1234567890123456")
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

// 解密
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// CBC 模式下的解密对象
	iv := []byte("1234567890123456")
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除 padding
	plaintext = pkcs7UnPadding(plaintext)
	return plaintext, nil
}

// 对明文进行 padding
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去除 padding
func pkcs7UnPadding(plaintext []byte) []byte {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	return plaintext[:(length - unpadding)]
}

// 获取cpu使用率
func GetCpuPercent() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	return percent[0], nil
}

// 获取内存使用率
func GetRamPercent() (float64, error) {
	ram_info, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return ram_info.UsedPercent, nil
}

// 获取cpu温度
func GetCpuTemp() (int, error) {
	cmd := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	tempStr := strings.Replace(out.String(), "\n", "", -1)
	temp, err := strconv.Atoi(tempStr)
	if err != nil {
		return 0, err
	}
	temp = temp / 1000
	return temp, nil
}
