package service

import (
	"gin-synolux/common"
	"gin-synolux/db"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/repository"

	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
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
func (s *UserService) Login(req *dto.UserLoginReq) (*dto.UserLoginResp, error) {

	enableCaptcha := viper.GetBool("enable_captcha")

	// 验证码
	if enableCaptcha && !common.Store.Verify(req.Key, req.Code, true) {
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
		"login_ip":   req.Ip,
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
func (s *UserService) SetPassword(req *dto.UserSetPasswordReq) error {

	// ===== 业务校验 =====
	if req.NewPassword != req.RePassword {
		return common.NewError(-2, "确认密码不一致")
	}

	// ===== 查用户（走 repo）=====
	user, err := s.repo.GetByID(req.UID)
	if err != nil {
		log.Error("查询用户失败", err)
		return common.NewError(-3, "系统错误")
	}
	if user == nil {
		return common.NewError(-3, "用户不存在")
	}

	// ===== 校验旧密码 =====
	if !common.PasswordVerify(req.Password, user.Password) {
		log.Warn("原始密码错误 uid=" + req.UID)
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

	if err := s.repo.Update(req.UID, updateData); err != nil {
		log.Error("密码更新失败", err)
		return common.NewError(-6, "密码更新失败")
	}

	return nil
}

// 获取详情
func (s *UserService) GetById(req dto.UserDetailReq) (*models.User, error) {
	// 👉 走 repo
	user, err := s.repo.GetByID(req.UID)
	if err != nil {
		log.Error("查询用户失败 id="+req.UID, err)
		return nil, common.NewError(-2, "查询失败")
	}

	if user == nil {
		return nil, common.NewError(-3, "用户不存在")
	}

	return user, nil
}

// 注册
func (s *UserService) Register(req dto.UserRegisterReq, isAdmin bool) error {

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

	// ===== 用户名是否存在 =====
	exists, err := repo.ExistsByUsername(req.Username)
	if err != nil {
		tx.Rollback()
		return common.NewError(-5, "系统错误")
	}
	if exists {
		tx.Rollback()
		return common.NewError(-6, "用户名已存在")
	}

	// ===== 密码加密 =====
	hash, err := common.PasswordHash(req.Password)
	if err != nil {
		tx.Rollback()
		return common.NewError(-7, "系统错误")
	}
	
	user := models.User{
		Id:         common.UniqueId(),
		Username:   req.Username,
		Password:   hash,
		Realname:   req.Realname,
		Email:      req.Email,
		PhoneCode:  req.PhoneCode,
		Phone:      req.Phone,
		Avatar:     req.Avatar,
		Sex:        req.Sex,
		Status:     1,
		Language:   "cn",
		RegIp:		req.RegIp,
		CreateTime: now,
	}

	if err := repo.Create(&user); err != nil {
		tx.Rollback()
		return common.NewError(-8, "用户添加失败")
	}

	// ===== 提交事务 =====
	if err := tx.Commit().Error; err != nil {
		return common.NewError(-9, "事务失败")
	}

	return nil
}

//修改用户信息
func (s *UserService) UpdateProfile(req dto.UserProfileReq, isAdmin bool) error {

	tx := db.DB.Self.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	repo := s.repo.WithTx(tx)
	now := common.Timestamp()

	// ===== 是否存在 =====
	exists, err := repo.ExistsByID(req.ID)
	if err != nil {
		tx.Rollback()
		return common.NewError(-2, "查询失败")
	}
	if !exists {
		tx.Rollback()
		return common.NewError(-2, "用户不存在")
	}

	updateData := map[string]interface{}{
		"realname":    req.Realname,
		"email":       req.Email,
		"phone_code":  req.PhoneCode,
		"phone":       req.Phone,
		"avatar":      req.Avatar,
		"sex":         req.Sex,
		"update_time": now,
	}

	// 👉 可选密码更新
	if req.Password != "" {
		hash, err := common.PasswordHash(req.Password)
		if err != nil {
			tx.Rollback()
			return common.NewError(-3, "系统错误")
		}
		updateData["password"] = hash
	}

	if err := repo.Update(req.ID, updateData); err != nil {
		tx.Rollback()
		return common.NewError(-4, "用户更新失败")
	}

	if err := tx.Commit().Error; err != nil {
		return common.NewError(-9, "事务失败")
	}

	return nil
}