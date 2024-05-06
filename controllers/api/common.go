package controllers

import (
	"encoding/json"
	"gin-synolux/utils"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonController struct {
	BaseController
}

// 获取列表
func (c *CommonController) ChatGPT(ctx *gin.Context) {
	keyword := ctx.Query("keyword")

	stat, content := utils.ChatGPT(keyword)
	if !stat {
		c.ErrorJson(ctx, -1, "发送错误", nil)
		return
	}
	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["content"] = content

	c.SuccessJson(ctx, "success", resp)
}

// 返回客户端ip
func (c *CommonController) Ip(ctx *gin.Context) {
	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["ip"] = ctx.ClientIP()
	c.SuccessJson(ctx, "success", resp)
}

// 检测用,可查看是否返回信息及时间戳
func (c *CommonController) Ping(ctx *gin.Context) {
	c.SuccessJson(ctx, "pong", nil)
}

// 获取图形验证码
func (c *CommonController) Captcha(ctx *gin.Context) {
	id, b64s, _, err := utils.GetCaptcha()
	if err != nil {
		c.ErrorJson(ctx, -1, "生成验证码错误", nil)
		return
	}

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["key"] = id
	resp["img"] = b64s
	c.SuccessJson(ctx, "success", resp)
}

// 获取彩云天气数据
func (c *CommonController) Weather(ctx *gin.Context) {
	// dailysteps返回多少天数据
	//url := "https://api.caiyunapp.com/v2.6/3mvCFBpeZ4qWAtim/121.4159,31.0281/weather?alert=true&dailysteps=1&hourlysteps=24"
	url := "https://api.caiyunapp.com/v2.6/3mvCFBpeZ4qWAtim/121.4159,31.0281/weather?alert=true"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		c.ErrorJson(ctx, -1, "请求错误", nil)
		return
	}
	// body := res.Body
	// contentLength := res.ContentLength
	// contentType := res.Header.Get("Content-Type")
	// //数据写入响应体
	// ctx.DataFromReader(http.StatusOK, contentLength, contentType, body, nil)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var info interface{}
	json.Unmarshal(body, &info)
	c.SuccessJson(ctx, "success", &info)
}

// 获取系统信息
func (c *CommonController) Hardware(ctx *gin.Context) {
	//cpu使用率
	cpu_usage, _ := utils.GetCpuPercent()
	//cpu温度
	cpu_temp, _ := utils.GetCpuTemp()
	//内存使用率
	ram_usage, _ := utils.GetRamPercent()

	//获取wifi信息
	wifi_status := make(map[string]interface{})
	wifi_status["name"] = "乙太网路"
	wifi_status["value"] = 100

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["cpu_usage"] = cpu_usage
	resp["ram_usage"] = ram_usage
	resp["cpu_temp"] = cpu_temp
	resp["wifi_status"] = wifi_status
	c.SuccessJson(ctx, "success", resp)
}

// 获取Gist数据
func (c *CommonController) Gist(ctx *gin.Context) {
	url := "https://api.github.com/gists/82d4bb9fce4cd96752d66d7b3c832415"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		c.ErrorJson(ctx, -1, "请求错误", nil)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var info interface{}
	json.Unmarshal(body, &info)
	c.SuccessJson(ctx, "success", &info)
}
