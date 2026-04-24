// 公共方法
package common

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

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

func ByteEncoder(s interface{}) []byte {
	var enc_result bytes.Buffer
	enc := gob.NewEncoder(&enc_result)
	if err := enc.Encode(s); err != nil {
		Log.Fatal("encode error:", err)
	}

	return enc_result.Bytes()
}

//获取用户id
func GetUserID(c *gin.Context) string { 
	return c.GetString("userID") 
} 

//获取用户角色
func GetRole(c *gin.Context) string { 
	return c.GetString("role") 
}