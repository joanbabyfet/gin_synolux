package controllers

import (
	"gin-synolux/models"
	"gin-synolux/service"
	"gin-synolux/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
	Service *service.UserService //依赖注入
}

// 登录
func (c *UserController) Login(ctx *gin.Context) {
	username := ctx.PostForm("username") //帐号
	password := ctx.PostForm("password") //密码
	code := ctx.PostForm("code")         //验证码
	key := ctx.PostForm("key")           //验证码key
	ip := c.getClientIp(ctx)

	info, err := c.Service.Login(username, password, key, code, ip)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", info)
}

// 登录退出 (JWT + Redis 黑名单)
func (c *UserController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		c.ErrorJson(ctx, -1, "未登录", nil)
		return
	}

	// 去掉 Bearer
	parts := strings.Split(token, " ")
	if len(parts) == 2 {
		token = parts[1]
	}

	// 写入 Redis（设置过期时间 = token 剩余时间）
	err := utils.Redis.Set("jwt:blacklist:"+token, 1, time.Hour*24).Err()
	if err != nil {
		c.ErrorJson(ctx, -2, "退出失败", nil)
		return
	}

	c.SuccessJson(ctx, "success", nil)
}

// 修改密码
func (c *UserController) SetPassword(ctx *gin.Context) {
	password := ctx.PostForm("password")         //原始密码
	new_password := ctx.PostForm("new_password") //新密码
	re_password := ctx.PostForm("re_password")   //确认密码

	uid := ctx.GetString("userID")
	if uid == "" {
		c.ErrorJson(ctx, -1, "未登录", nil)
		return
	}

	//修改密码
	err := c.Service.SetPassword(password, new_password, re_password, uid)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	uid := ctx.GetString("userID")
	if uid == "" {
		c.ErrorJson(ctx, -1, "未登录", nil)
		return
	}

	//获取用户信息
	info, err := c.Service.GetById(uid)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", info)
}

// 注册
func (c *UserController) Register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	realname := ctx.PostForm("realname")
	email := ctx.PostForm("email")
	phone_code := ctx.PostForm("phone_code")
	phone := ctx.PostForm("phone")
	avatar := ctx.PostForm("avatar")

	sexVal := ctx.PostForm("sex")
	sex, err := strconv.Atoi(sexVal)
	if err != nil {
		sex = 0 // 默认值
	}

	// 构造实体（建议后面可换 DTO）
	entity := models.User{
		Username:  username,
		Password:  password,
		Realname:  realname,
		Email:     email,
		PhoneCode: phone_code,
		Phone:     phone,
		Avatar:    avatar,
		Sex:       int8(sex),
		RegIp:     c.getClientIp(ctx),
	}
	
	// 构造实体（建议后面可换 DTO）
	err = c.Service.Save(entity, false)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 修改用户信息
func (c *UserController) Profile(ctx *gin.Context) {
	realname := ctx.PostForm("realname")
	email := ctx.PostForm("email")
	phone_code := ctx.PostForm("phone_code")
	phone := ctx.PostForm("phone")
	avatar := ctx.PostForm("avatar")

	uid := ctx.GetString("userID")
	if uid == "" {
		c.ErrorJson(ctx, -1, "未登录", nil)
		return
	}

	sexVal := ctx.PostForm("sex")
	sex, err := strconv.Atoi(sexVal)
	if err != nil {
		sex = 0
	}

	entity := models.User{
		Id:        uid,
		Realname:  realname,
		Email:     email,
		PhoneCode: phone_code,
		Phone:     phone,
		Avatar:    avatar,
		Sex:       int8(sex),
	}

	//保存
	err = c.Service.Save(entity, false)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}