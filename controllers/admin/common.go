package admin

import (
	"gin-synolux/common"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type CommonController struct {
	AdminBaseController
}

// 获取列表
func (c *CommonController) ChatGPT(ctx *gin.Context) {
	keyword := ctx.Query("keyword")

	stat, content := common.ChatGPT(keyword)
	if !stat {
		common.Fail(ctx, -1, "发送错误", nil)
		return
	}
	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["content"] = content

	common.Success(ctx, resp)
}

// 返回客户端ip
func (c *CommonController) Ip(ctx *gin.Context) {
	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["ip"] = ctx.ClientIP()
	common.Success(ctx, resp)
}

// 检测用,可查看是否返回信息及时间戳
func (c *CommonController) Ping(ctx *gin.Context) {
	common.Success(ctx, nil)
}

// 获取图形验证码
func (c *CommonController) Captcha(ctx *gin.Context) {
	id, b64s, _, err := common.GetCaptcha()
	if err != nil {
		common.Fail(ctx, -1, "生成验证码错误", nil)
		return
	}

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["key"] = id
	resp["img"] = b64s
	common.Success(ctx, resp)
}

// 发送消息
func (c *CommonController) SendMsg(ctx *gin.Context) {
	message := ctx.PostForm("message")

	url := viper.GetString("rabbitmq_host")
	rabbitmq, err := common.NewRabbitMQ("queue1", "", "", url)
	defer rabbitmq.Destroy()
	if err != nil {
		common.Fail(ctx, -1, err.Error(), nil)
		return
	}
	//发送消息
	err = rabbitmq.Publish(message)
	if err != nil {
		common.Fail(ctx, -2, err.Error(), nil)
		return
	}
	common.Success(ctx, nil)
}
