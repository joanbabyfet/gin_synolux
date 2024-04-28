package admin

import (
	"gin-synolux/utils"

	"github.com/gin-gonic/gin"
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
