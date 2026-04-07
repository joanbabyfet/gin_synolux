package admin

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	AdminBaseController
	Service *service.AdminService //依赖注入
}

func NewAdminController(s *service.AdminService) *AdminController {
	return &AdminController{Service: s}
}

// 登录
func (c *AdminController) Login(ctx *gin.Context) {
	var req dto.AdminLoginReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}

	req.LoginIp = c.getClientIp(ctx)

	info, err := c.Service.Login(&req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, info)
}

// 登录退出 (JWT + Redis 黑名单)
func (c *AdminController) Logout(ctx *gin.Context) {
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
func (c *AdminController) SetPassword(ctx *gin.Context) {
	var req dto.AdminSetPasswordReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
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
func (c *AdminController) GetUserInfo(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	if uid == "" {
		common.Fail(ctx, -1, "未登录", nil)
		return
	}

	var req dto.AdminDetailReq
	req.UID = uid

	//获取用户信息
	info, err := c.Service.GetById(req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, info)

}

// 注册
func (c *AdminController) Register(ctx *gin.Context) {
	var req dto.AdminRegisterReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
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
func (c *AdminController) Profile(ctx *gin.Context) {
	var req dto.AdminProfileReq

	// 统一绑定（支持 form / json）
	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
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