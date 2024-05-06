package admin

import (
	"gin-synolux/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type CommonController struct {
	AdminBaseController
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

// 发送消息
func (c *CommonController) SendMsg(ctx *gin.Context) {
	message := ctx.PostForm("message")

	url := viper.GetString("rabbitmq_host")
	rabbitmq, err := utils.NewRabbitMQ("queue1", "", "", url)
	defer rabbitmq.Destroy()
	if err != nil {
		c.ErrorJson(ctx, -1, err.Error(), nil)
		return
	}
	//发送消息
	err = rabbitmq.Publish(message)
	if err != nil {
		c.ErrorJson(ctx, -2, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}
