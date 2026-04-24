package service

import (
	"gin-synolux/common"
	"gin-synolux/db"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/repository"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type AdminService struct {
	repo *repository.AdminRepo
}

func NewAdminService(db *gorm.DB) *AdminService {
	if db == nil {
		panic("db is nil (service)")
	}

	return &AdminService{
		repo: repository.NewAdminRepo(db),
	}
}

// Login 登录处理
func (s *AdminService) Login(req *dto.AdminLoginReq) (*dto.AdminLoginResp, error) {

	enableCaptcha := viper.GetBool("enable_captcha")

	// 验证码
	if enableCaptcha && !common.Store.Verify(req.Key, req.Code, true) {
		return nil, common.NewError(-2, "验证码错误")
	}

	// 用 repo
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil || user == nil {
		return nil, common.NewError(-3, "管理员名或密码无效")
	}

	if !common.PasswordVerify(req.Password, user.Password) {
		return nil, common.NewError(-3, "管理员名或密码无效")
	}

	if user.Status == 0 {
		return nil, common.NewError(-4, "账号已被禁用")
	}

	// 更新登录信息
	updateData := map[string]interface{}{
		"login_ip":   req.LoginIp,
		"login_time": common.Timestamp(),
	}

	if err := s.repo.Update(user.Id, updateData); err != nil {
		common.Log.Error("更新登录信息失败", err)
		return nil, common.NewError(-5, "登录异常")
	}

	// 用 pkg/jwt
	token, err := common.GenerateToken(user.Id, common.RoleAdmin, 0)
	if err != nil {
		common.Log.Error("生成token失败", err)
		return nil, common.NewError(-6, "生成token失败")
	}

	return &dto.AdminLoginResp{
		Id:        user.Id,
		Username:  user.Username,
		Realname:  user.Realname,
		Email:     user.Email,
		Token:     token,
	}, nil
}

// 修改密码
func (s *AdminService) SetPassword(req *dto.AdminSetPasswordReq) error {

	// ===== 业务校验 =====
	if req.NewPassword != req.RePassword {
		return common.NewError(-2, "确认密码不一致")
	}

	// ===== 查管理员（走 repo）=====
	user, err := s.repo.GetByID(req.UID)
	if err != nil {
		common.Log.Error("查询管理员失败", err)
		return common.NewError(-3, "系统错误")
	}
	if user == nil {
		return common.NewError(-3, "管理员不存在")
	}

	// ===== 校验旧密码 =====
	if !common.PasswordVerify(req.Password, user.Password) {
		common.Log.Warn("原始密码错误 uid=" + req.UID)
		return common.NewError(-4, "原始密码错误")
	}

	// ===== 新密码加密 =====
	hash, err := common.PasswordHash(req.NewPassword)
	if err != nil {
		common.Log.Error("密码加密失败", err)
		return common.NewError(-5, "系统错误")
	}

	// ===== 更新（map方式，避免零值问题）=====
	updateData := map[string]interface{}{
		"password": hash,
	}

	if err := s.repo.Update(req.UID, updateData); err != nil {
		common.Log.Error("密码更新失败", err)
		return common.NewError(-6, "密码更新失败")
	}

	return nil
}

// 获取详情
func (s *AdminService) GetById(req dto.AdminDetailReq) (*models.Admin, error) {

	// 👉 走 repo
	user, err := s.repo.GetByID(req.UID)
	if err != nil {
		common.Log.Error("查询管理员失败 id="+req.UID, err)
		return nil, common.NewError(-2, "查询失败")
	}

	if user == nil {
		return nil, common.NewError(-3, "管理员不存在")
	}

	return user, nil
}

// 注册
func (s *AdminService) Register(req dto.AdminRegisterReq, isAdmin bool) error {

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

	// ===== 管理员名是否存在 =====
	exists, err := repo.ExistsByUsername(req.Username)
	if err != nil {
		tx.Rollback()
		return common.NewError(-5, "系统错误")
	}
	if exists {
		tx.Rollback()
		return common.NewError(-6, "管理员名已存在")
	}

	// ===== 密码加密 =====
	hash, err := common.PasswordHash(req.Password)
	if err != nil {
		tx.Rollback()
		return common.NewError(-7, "系统错误")
	}
	
	user := models.Admin{
		Id:         common.UniqueId(),
		Username:   req.Username,
		Password:   hash,
		Realname:   req.Realname,
		Email:      req.Email,
		Sex:        req.Sex,
		Status:     1,
		RegIp:		req.RegIp,
		CreateTime: now,
	}

	if err := repo.Create(&user); err != nil {
		tx.Rollback()
		return common.NewError(-8, "管理员添加失败")
	}

	// ===== 提交事务 =====
	if err := tx.Commit().Error; err != nil {
		return common.NewError(-9, "事务失败")
	}

	return nil
}

//修改管理员信息
func (s *AdminService) UpdateProfile(req dto.AdminProfileReq, isAdmin bool) error {

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
		return common.NewError(-2, "管理员不存在")
	}

	updateData := map[string]interface{}{
		"realname":    req.Realname,
		"email":       req.Email,
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
		return common.NewError(-4, "管理员更新失败")
	}

	if err := tx.Commit().Error; err != nil {
		return common.NewError(-9, "事务失败")
	}

	return nil
}