package controllers

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
	Service *service.UserService //依赖注入
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{Service: s}
}

// 登录
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}

	req.Ip = c.getClientIp(ctx)

	info, err := c.Service.Login(&req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, info)
}

// 登录退出 (JWT + Redis 黑名单)
func (c *UserController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		common.Fail(ctx, -1, "未登录", nil)
		return
	}

	// 去掉 Bearer
	parts := strings.Split(token, " ")
	if len(parts) == 2 {
		token = parts[1]
	}

	// 写入 Redis（设置过期时间 = token 剩余时间）
	err := common.Redis.Set("jwt:blacklist:"+token, 1, time.Hour*24).Err()
	if err != nil {
		common.Fail(ctx, -2, "退出失败", nil)
		return
	}

	common.Success(ctx, nil)
}

// 修改密码
func (c *UserController) SetPassword(ctx *gin.Context) {
	var req dto.UserSetPasswordReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}

	uid := common.GetUserID(ctx)
	if uid == "" {
		common.Fail(ctx, -1, "未登录", nil)
		return
	}

	req.UID = uid

	if err := c.Service.SetPassword(&req); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}

// 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	if uid == "" {
		common.Fail(ctx, -1, "未登录", nil)
		return
	}

	var req dto.UserDetailReq
	req.UID = uid

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}

	//获取用户信息
	info, err := c.Service.GetById(req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, info)
}

// 注册
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.UserRegisterReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}

	req.RegIp = c.getClientIp(ctx)

	if err := c.Service.Register(req, false); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}

// 修改用户信息
func (c *UserController) Profile(ctx *gin.Context) {
	var req dto.UserProfileReq

	// 统一绑定（支持 form / json）
	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}

	// 登录校验
	uid := common.GetUserID(ctx)
	if uid == "" {
		common.Fail(ctx, -1, "未登录", nil)
		return
	}

	// 注入JWT信息（非常关键）
	req.ID = uid

	// 调用 service
	err := c.Service.UpdateProfile(req, false)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}