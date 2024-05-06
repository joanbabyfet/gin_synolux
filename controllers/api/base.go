// 父控制器
package controllers

import (
	"gin-synolux/consts"
	"gin-synolux/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

// 封装接口统一返回json格式
type ReturnMsg struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Timestamp int         `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// 初始化, 先于Prepare函数
func init() {
}

// 定义prepare方法, 用户扩展用
func (c *BaseController) Prepare() {

}

// @Title API成功响应
// @Description API成功响应
// @Param msg 成功消息
// @Param data 成功返回信息
func (c *BaseController) SuccessJson(ctx *gin.Context, msg string, data interface{}) {
	if msg == "" {
		msg = "success"
	}
	if data == nil || data == "" {
		data = struct{}{}
	}
	timestamp := utils.Timestamp()
	res := &ReturnMsg{
		consts.SUCCESS, msg, timestamp, data, //0=成功
	}
	ctx.JSON(http.StatusOK, res)
}

// @Title API失败响应
// @Description API失败响应
// @Param code 错误码
// @Param msg 异常消息
// @Param data 异常返回信息
func (c *BaseController) ErrorJson(ctx *gin.Context, code int, msg string, data interface{}) {
	if code >= 0 {
		code = consts.UNKNOWN_ERROR_STATUS
	}
	if msg == "" {
		msg = "error"
	}
	if data == nil || data == "" {
		data = struct{}{}
	}
	timestamp := utils.Timestamp()
	res := &ReturnMsg{
		code, msg, timestamp, data,
	}
	ctx.JSON(http.StatusOK, res)
}

// 获取客户端ip
func (c *BaseController) getClientIp(ctx *gin.Context) string {
	return ctx.ClientIP()
}
