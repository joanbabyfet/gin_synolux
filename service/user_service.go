package service

import (
	"errors"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/utils"
	"strings"

	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/govalidator"
)

type UserService struct {
}

// Login 登录处理
func (s *UserService) Login(username string, password string, key string, code string, ip string) (*dto.UserLoginResp, error) {
	enable_captcha := viper.GetBool("enable_captcha") //是否启用验证码
	
	// 参数验证, 使用 DTO 承载参数
	entity := dto.UserLogin{
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
    if err := utils.ValidateStruct(&entity, rules, messages); err != nil {
        return nil, NewServiceError(-1, err.Error())
    }
	
	//检测验证码
	if enable_captcha && !utils.Store.Verify(key, code, true) {
		return nil, NewServiceError(-2, "验证码错误")
	}

	var u models.User
	err := models.DB.Self.Where("username = ?", entity.Username).First(&u).Error
	if err != nil {
		log.Error("用户名或密码无效", err)
		return nil, NewServiceError(-3, "用户名或密码无效")
	}
	if !utils.PasswordVerify(entity.Password, u.Password) {
		log.Error("用户名或密码无效", err)
		return nil, NewServiceError(-3, "用户名或密码无效")
	}
	if u.Status == 0 {
		return nil, NewServiceError(-4, "账号已被禁用")
	}

	//更新用户信息
	u.LoginIp = ip
	u.LoginTime = utils.Timestamp()
	if err := models.DB.Self.Save(&u).Error; err != nil {
		log.Error("登录异常", err)
		return nil, NewServiceError(-5, "登录异常")
	}

	// 生成token
	token, err := models.GenerateToken(&u, u.Id, 0)
	if err != nil {
		log.Error("生成token失败", err)
		return nil, NewServiceError(-6, "生成token失败")
	}

	//组装响应数据
	resp := &dto.UserLoginResp{
		Id:        u.Id,
		Username:  u.Username,
		Realname:  u.Realname,
		Email:     u.Email,
		PhoneCode: u.PhoneCode,
		Phone:     u.Phone,
		Avatar:    utils.DisplayImg(u.Avatar),
		Language:  u.Language,
		Token:     token,
	}
	return resp, nil
}

// 修改密码
func (s *UserService) SetPassword(password string, newPassword string, rePassword string, uid string) error {
	// 参数验证, 使用 DTO 承载参数
	data := dto.Password{
		Password: strings.TrimSpace(password),
		NewPassword: strings.TrimSpace(newPassword),
		RePassword: strings.TrimSpace(rePassword),
		Uid: uid,
	}
	rules := govalidator.MapData{
		"password": {"required"},
		"new_password": {"required"},
		"re_password": {"required"},
	}
	messages := govalidator.MapData{
		"password": {"required:password 不能为空"},
		"new_password": {"required:new_password 不能为空"},
		"re_password": {"required:re_password 不能为空"},
	}
    if err := utils.ValidateStruct(&data, rules, messages); err != nil {
        return NewServiceError(-1, err.Error())
    }

	//检测输入密码是否一致
	if data.RePassword != data.NewPassword {
		return NewServiceError(-2, "确认密码不一样")
	}
	
	var user models.User
	//查用户（建议直接查）
	if err := models.DB.Self.First(&user, "id = ?", uid).Error; err != nil {
		log.Error("用户不存在", err)
		return NewServiceError(-3, "用户不存在")
	}
	if !utils.PasswordVerify(data.Password, user.Password) {
		log.Warn("原始密码错误")
		return NewServiceError(-4, "原始密码错误")
	}

	//获取加密后密码
	newPasswordHash, err := utils.PasswordHash(data.NewPassword)
	if err != nil {
		log.Error("密码加密失败", err)
		return NewServiceError(-5, "系统错误")
	}

	//更新（只更新 password）
	err = models.DB.Self.Model(&user).Update("password", newPasswordHash).Error
	if err != nil {
		log.Error("密码更新失败", err)
		return NewServiceError(-6, "密码更新失败")
	}
	return nil
}

// 获取详情
func (s *UserService) GetById(id string) (*models.User, error) {
	entity := new(models.User)
	info, err := entity.GetById(id)
	if err != nil {
		log.Error("用户不存在 "+id, err)
		return nil, errors.New("用户不存在")
	}
	return info, nil
}

// 保存
func (s *UserService) Save(data models.User, isAdmin bool) (error) {
	rules := govalidator.MapData{
		"realname":   []string{"required"},
		"email":      []string{"required"},
		"phone_code": []string{"required"},
		"phone":      []string{"required"},
    }
    messages := govalidator.MapData{
		"realname":   []string{"required:realname 不能为空"},
		"email":      []string{"required:email 不能为空"},
		"phone_code": []string{"required:phone_code 不能为空"},
		"phone":      []string{"required:phone 不能为空"},
    }

	isUpdate := data.Id != ""

    // 更新操作必须验证 ID
    if isUpdate {
        rules["id"] = []string{"required"}
        messages["id"] = []string{"required:id 不能为空"}
    } else {
		// 注册：必须 username + password
		rules["username"] = []string{"required"}
		rules["password"] = []string{"required"}

		messages["username"] = []string{"required:username 不能为空"}
		messages["password"] = []string{"required:password 不能为空"}
	}
	
    if err := utils.ValidateStruct(&data, rules, messages); err != nil {
        return NewServiceError(-1, err.Error())
    }

	//开启事务
	tx := models.DB.Self.Begin() 
	var err error
	defer func() {
		//防止紧急停止
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		//不用再手动 rollback
		if err != nil {
			tx.Rollback()
		}
	}()

	if isUpdate {
		//检测数据是否存在
		var user models.User
		if err = tx.First(&user, "id = ?", data.Id).Error; err != nil {
			log.Error("用户不存在", err)
			return NewServiceError(-2, "用户不存在")
		}

		updateData := map[string]interface{}{
			"realname":    data.Realname,
			"email":       data.Email,
			"phone_code":  data.PhoneCode,
			"phone":       data.Phone,
			"avatar":      data.Avatar,
			"sex":         data.Sex,
			"update_user": "1",
			"update_time": utils.Timestamp(),
		}

		// 密码非空才更新
		if data.Password != "" {
			var hash string
			hash, err = utils.PasswordHash(data.Password)
			if err != nil {
				log.Error("密码加密失败", err)
				return NewServiceError(-3, "系统错误")
			}
			updateData["password"] = hash
		}

		if err = tx.Model(&user).Updates(updateData).Error; err != nil {
			log.Error("用户更新失败", err)
			return NewServiceError(-4, "用户更新失败")
		}
	} else {
		var count int64
		if err = tx.Model(&models.User{}).Where("username = ?", data.Username).Count(&count).Error; err != nil {
			log.Error("检查用户名失败", err)
			return NewServiceError(-5, "系统错误")
		}
		if count > 0 {
			return NewServiceError(-6, "用户名已存在")
		}

		var hash string
		hash, err = utils.PasswordHash(data.Password)
		if err != nil {
			log.Error("密码加密失败", err)
			return NewServiceError(-7, "系统错误")
		}

		data.Id = utils.UniqueId()
		data.Password = hash
		data.Language = "cn"
		data.Status = 1
		data.CreateUser = "1"
		data.CreateTime = utils.Timestamp()

		if err = tx.Create(&data).Error; err != nil {
			log.Error("用户添加失败", err)
			return NewServiceError(-8, "用户添加失败")
		}
	}
	
	if err = tx.Commit().Error; err != nil {
		log.Error("事务提交失败", err)
		return NewServiceError(-9, err.Error())
	}

	if isAdmin {
		if isUpdate {
			log.Info("更新用户 id=" + data.Id)
		} else {
			log.Info("新增用户 id=" + data.Id)
		}
	}

	return nil
}
