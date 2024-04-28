package controllers

import (
	"encoding/base64"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/service"
	"gin-synolux/utils"
	"strconv"
	"strings"

	"github.com/beego/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UserController struct {
	BaseController
}

// 登录
func (c *UserController) Login(ctx *gin.Context) {
	username := ctx.PostForm("username") //帐号
	password := ctx.PostForm("password") //密码
	code := ctx.PostForm("code")         //验证码
	key := ctx.PostForm("key")           //验证码key
	login := dto.UserLogin{}             //登录请求格式
	login.Username = username
	login.Password = password
	login.Code = code
	login.Key = key
	login.LoginIp = c.getClientIp(ctx)
	enable_captcha := viper.GetBool("enable_captcha") //是否启用验证码

	//用户密码解密
	pwd, _ := base64.StdEncoding.DecodeString(login.Password)
	login.Password = string(pwd)

	//参数验证
	valid := validation.Validation{}
	valid.Required(login.Username, "username")
	valid.Required(login.Password, "password")
	if enable_captcha {
		valid.Required(login.Key, "key")
		valid.Required(login.Code, "code")
	}
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	//检测验证码
	if enable_captcha && !utils.Store.Verify(key, code, true) {
		c.ErrorJson(ctx, -2, "验证码错误", nil)
		return
	}

	//登录
	service_user := new(service.UserService)
	stat, user, err := service_user.Login(login)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", user)
}

// 登录退出
func (c *UserController) Logout(ctx *gin.Context) {
	//c.DestroySession()
	c.SuccessJson(ctx, "success", nil)
}

// 修改密码
func (c *UserController) SetPassword(ctx *gin.Context) {
	password := ctx.PostForm("password")         //原始密码
	new_password := ctx.PostForm("new_password") //新密码
	re_password := ctx.PostForm("re_password")   //确认密码
	auth := ctx.Request.Header.Get("Authorization")
	kv := strings.Split(auth, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.ErrorJson(ctx, -1, "未带token", nil)
		return
	}
	token := kv[1]
	payload, err := models.ValidateToken(token)
	if err != nil {
		c.ErrorJson(ctx, -2, "未登录或登录超时", nil)
		return
	}
	uid := payload.UserID //取得用户id

	//用户密码解密
	pwd, _ := base64.StdEncoding.DecodeString(password)
	new_pwd, _ := base64.StdEncoding.DecodeString(new_password)
	re_pwd, _ := base64.StdEncoding.DecodeString(re_password)

	dto := dto.Password{}
	dto.Password = string(pwd)
	dto.NewPassword = string(new_pwd)
	dto.RePassword = string(re_pwd)
	dto.Uid = uid

	//参数验证
	valid := validation.Validation{}
	valid.Required(dto.Password, "password")
	valid.Required(dto.NewPassword, "new_password")
	valid.Required(dto.RePassword, "re_password")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -3, err.Key+err.Error(), nil)
			return
		}
	}
	//检测输入密码是否一致
	if dto.RePassword != dto.NewPassword {
		c.ErrorJson(ctx, -4, "确认密码不一样", nil)
		return
	}

	//修改密码
	service_user := new(service.UserService)
	stat, err := service_user.SetPassword(dto)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	kv := strings.Split(auth, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.ErrorJson(ctx, -1, "未带token", nil)
		return
	}
	token := kv[1]
	payload, err := models.ValidateToken(token)
	if err != nil {
		c.ErrorJson(ctx, -2, "未登录或登录超时", nil)
		return
	}
	uid := payload.UserID //取得用户id

	//获取用户信息
	service_user := new(service.UserService)
	info, err := service_user.GetById(uid)
	if err != nil {
		c.ErrorJson(ctx, -3, err.Error(), nil)
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
	sex, _ := strconv.Atoi(ctx.PostForm("sex"))
	avatar := ctx.PostForm("avatar")
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

	//用户密码解密
	pwd, _ := base64.StdEncoding.DecodeString(entity.Password)
	entity.Password = string(pwd)

	//参数验证
	valid := validation.Validation{}
	valid.Required(entity.Username, "username")
	valid.Required(entity.Password, "password")
	valid.Required(entity.Realname, "realname")
	valid.Required(entity.Email, "email")
	valid.Required(entity.PhoneCode, "phone_code")
	valid.Required(entity.Phone, "phone")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	//保存
	service_user := new(service.UserService)
	stat, err := service_user.Save(entity)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 修改用户信息
func (c *UserController) Profile(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	kv := strings.Split(auth, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.ErrorJson(ctx, -1, "未带token", nil)
		return
	}
	token := kv[1]
	payload, err := models.ValidateToken(token)
	if err != nil {
		c.ErrorJson(ctx, -2, "未登录或登录超时", nil)
		return
	}
	uid := payload.UserID //取得用户id

	realname := ctx.PostForm("realname")
	email := ctx.PostForm("email")
	phone_code := ctx.PostForm("phone_code")
	phone := ctx.PostForm("phone")
	sex, _ := strconv.Atoi(ctx.PostForm("sex"))
	avatar := ctx.PostForm("avatar")
	entity := models.User{
		Id:        uid,
		Realname:  realname,
		Email:     email,
		PhoneCode: phone_code,
		Phone:     phone,
		Avatar:    avatar,
		Sex:       int8(sex),
	}

	//参数验证
	valid := validation.Validation{}
	valid.Required(entity.Realname, "realname")
	valid.Required(entity.Email, "email")
	valid.Required(entity.PhoneCode, "phone_code")
	valid.Required(entity.Phone, "phone")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -3, err.Key+err.Error(), nil)
			return
		}
	}

	//保存
	service_user := new(service.UserService)
	stat, err := service_user.Save(entity)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}
