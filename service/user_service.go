package service

import (
	"gin-synolux/common"
	"gin-synolux/db"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/repository"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/govalidator"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(db *gorm.DB) *UserService {
	if db == nil {
		panic("db is nil (service)")
	}

	return &UserService{
		repo: repository.NewUserRepo(db),
	}
}

// Login 登录处理
func (s *UserService) Login(username, password, key, code, ip string) (*dto.UserLoginResp, error) {

	enableCaptcha := viper.GetBool("enable_captcha")

	req := dto.UserLoginReq{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}
	rules := govalidator.MapData{
		"username": {"required"},
		"password": {"required"},
	}
	messages := govalidator.MapData{
		"username": {"required:username 不能为空"},
		"password": {"required:password 不能为空"},
	}

	if err := common.ValidateStruct(&req, rules, messages); err != nil {
		return nil, common.NewError(-1, err.Error())
	}

	// 验证码
	if enableCaptcha && !common.Store.Verify(key, code, true) {
		return nil, common.NewError(-2, "验证码错误")
	}

	// 用 repo
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil || user == nil {
		return nil, common.NewError(-3, "用户名或密码无效")
	}

	if !common.PasswordVerify(req.Password, user.Password) {
		return nil, common.NewError(-3, "用户名或密码无效")
	}

	if user.Status == 0 {
		return nil, common.NewError(-4, "账号已被禁用")
	}

	// 更新登录信息
	updateData := map[string]interface{}{
		"login_ip":   ip,
		"login_time": common.Timestamp(),
	}

	if err := s.repo.Update(user.Id, updateData); err != nil {
		log.Error("更新登录信息失败", err)
		return nil, common.NewError(-5, "登录异常")
	}

	// 用 pkg/jwt
	token, err := common.GenerateToken(user.Id, common.RoleUser, 0)
	if err != nil {
		log.Error("生成token失败", err)
		return nil, common.NewError(-6, "生成token失败")
	}

	return &dto.UserLoginResp{
		Id:        user.Id,
		Username:  user.Username,
		Realname:  user.Realname,
		Email:     user.Email,
		PhoneCode: user.PhoneCode,
		Phone:     user.Phone,
		Avatar:    common.DisplayImg(user.Avatar),
		Language:  user.Language,
		Token:     token,
	}, nil
}

// 修改密码
func (s *UserService) SetPassword(password, newPassword, rePassword, uid string) error {

	// ===== DTO 参数承载 =====
	req := dto.UserPasswordReq{
		Password:    strings.TrimSpace(password),
		NewPassword: strings.TrimSpace(newPassword),
		RePassword:  strings.TrimSpace(rePassword),
		Uid:         uid,
	}

	rules := govalidator.MapData{
		"password":     {"required"},
		"new_password": {"required"},
		"re_password":  {"required"},
	}
	messages := govalidator.MapData{
		"password":     {"required:password 不能为空"},
		"new_password": {"required:new_password 不能为空"},
		"re_password":  {"required:re_password 不能为空"},
	}

	if err := common.ValidateStruct(&req, rules, messages); err != nil {
		return common.NewError(-1, err.Error())
	}

	// ===== 业务校验 =====
	if req.NewPassword != req.RePassword {
		return common.NewError(-2, "确认密码不一致")
	}

	// ===== 查用户（走 repo）=====
	user, err := s.repo.GetByID(req.Uid)
	if err != nil {
		log.Error("查询用户失败", err)
		return common.NewError(-3, "系统错误")
	}
	if user == nil {
		return common.NewError(-3, "用户不存在")
	}

	// ===== 校验旧密码 =====
	if !common.PasswordVerify(req.Password, user.Password) {
		log.Warn("原始密码错误 uid=" + req.Uid)
		return common.NewError(-4, "原始密码错误")
	}

	// ===== 新密码加密 =====
	hash, err := common.PasswordHash(req.NewPassword)
	if err != nil {
		log.Error("密码加密失败", err)
		return common.NewError(-5, "系统错误")
	}

	// ===== 更新（map方式，避免零值问题）=====
	updateData := map[string]interface{}{
		"password": hash,
	}

	if err := s.repo.Update(req.Uid, updateData); err != nil {
		log.Error("密码更新失败", err)
		return common.NewError(-6, "密码更新失败")
	}

	return nil
}

// 获取详情
func (s *UserService) GetById(id string) (*models.User, error) {

	// 参数校验（建议加上）
	if id == "" {
		return nil, common.NewError(-1, "id 不能为空")
	}

	// 👉 走 repo
	user, err := s.repo.GetByID(id)
	if err != nil {
		log.Error("查询用户失败 id="+id, err)
		return nil, common.NewError(-2, "查询失败")
	}

	if user == nil {
		return nil, common.NewError(-3, "用户不存在")
	}

	return user, nil
}

// 保存
func (s *UserService) Save(data models.User) error {

	isUpdate := data.Id != ""

	// ===== 参数校验 =====
	rules := govalidator.MapData{
		"realname":   {"required"},
		"email":      {"required"},
		"phone_code": {"required"},
		"phone":      {"required"},
	}
	messages := govalidator.MapData{
		"realname":   {"required:realname 不能为空"},
		"email":      {"required:email 不能为空"},
		"phone_code": {"required:phone_code 不能为空"},
		"phone":      {"required:phone 不能为空"},
	}

	if isUpdate {
		rules["id"] = []string{"required"}
		messages["id"] = []string{"required:id 不能为空"}
	} else {
		rules["username"] = []string{"required"}
		rules["password"] = []string{"required"}

		messages["username"] = []string{"required:username 不能为空"}
		messages["password"] = []string{"required:password 不能为空"}
	}

	if err := common.ValidateStruct(&data, rules, messages); err != nil {
		return common.NewError(-1, err.Error())
	}

	// ===== 开启事务 =====
	tx := db.DB.Self.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	repo := s.repo.WithTx(tx)
	now := common.Timestamp()

	if isUpdate {

		// ===== 是否存在 =====
		exists, err := repo.ExistsByID(data.Id)
		if err != nil {
			tx.Rollback()
			return common.NewError(-2, "查询失败")
		}
		if !exists {
			tx.Rollback()
			return common.NewError(-2, "用户不存在")
		}

		// ===== 更新字段 =====
		updateData := map[string]interface{}{
			"realname":    data.Realname,
			"email":       data.Email,
			"phone_code":  data.PhoneCode,
			"phone":       data.Phone,
			"avatar":      data.Avatar,
			"sex":         data.Sex,
			"update_time": now,
			"update_user": data.UpdateUser,
		}

		// 👉 密码非空才更新
		if data.Password != "" {
			hash, err := common.PasswordHash(data.Password)
			if err != nil {
				tx.Rollback()
				return common.NewError(-3, "系统错误")
			}
			updateData["password"] = hash
		}

		if err := repo.Update(data.Id, updateData); err != nil {
			tx.Rollback()
			return common.NewError(-4, "用户更新失败")
		}

	} else {

		// ===== 用户名是否存在 =====
		exists, err := repo.ExistsByUsername(data.Username)
		if err != nil {
			tx.Rollback()
			return common.NewError(-5, "系统错误")
		}
		if exists {
			tx.Rollback()
			return common.NewError(-6, "用户名已存在")
		}

		// ===== 密码加密 =====
		hash, err := common.PasswordHash(data.Password)
		if err != nil {
			tx.Rollback()
			return common.NewError(-7, "系统错误")
		}

		data.Id = common.UniqueId()
		data.Password = hash
		data.Language = "cn"
		data.Status = 1
		data.CreateTime = now

		if err := repo.Create(&data); err != nil {
			tx.Rollback()
			return common.NewError(-8, "用户添加失败")
		}
	}

	// ===== 提交事务 =====
	if err := tx.Commit().Error; err != nil {
		return common.NewError(-9, "事务失败")
	}

	return nil
}