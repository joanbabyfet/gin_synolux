package common

import (
	"gin-synolux/consts"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// 成功
func Success(ctx *gin.Context, data interface{}, msg ...string) {
	m := "success"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}

	if data == nil {
		data = struct{}{}
	}

	ctx.JSON(http.StatusOK, Response{
		Code:      0,
		Msg:       m,
		Timestamp: time.Now().Unix(),
		Data:      data,
	})
}

// 失败
func Fail(ctx *gin.Context, code int, msg string, data interface{}) {
	if code == 0 {
		code = -1
	}
	if msg == "" {
		msg = "error"
	}
	if data == nil {
		data = struct{}{}
	}

	ctx.JSON(http.StatusOK, Response{
		Code:      code,
		Msg:       msg,
		Timestamp: time.Now().Unix(),
		Data:      data,
	})
}

//业务错误（Service 层）, 比如 数据不存在／权限不足／状态不允许删除
func HandleError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	// 业务错误
	if se, ok := err.(*ServiceError); ok {
		Fail(ctx, se.Code, se.Msg, nil)
		return
	}

	// 未知错误
	Fail(ctx, consts.UNKNOWN_ERROR_STATUS, "internal server error", nil)
}