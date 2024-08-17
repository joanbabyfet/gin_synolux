package service

import (
	"errors"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/utils"
	"strings"

	"github.com/lexkong/log"
)

type UserService struct {
}

// Login 登录处理
func (s *UserService) Login(login dto.UserLogin) (int, *dto.UserLoginResp, error) {
	u := models.User{}
	err := models.DB.Self.Where("username = ?", login.Username).First(&u).Error
	if err != nil {
		log.Error("用户名或密码无效", err)
		return -3, nil, errors.New("用户名或密码无效")
	}
	if !utils.PasswordVerify(strings.Trim(login.Password, " "), u.Password) {
		log.Error("用户名或密码无效", err)
		return -4, nil, errors.New("用户名或密码无效")
	}
	//更新用户信息
	u.LoginIp = login.LoginIp
	u.LoginTime = utils.Timestamp()
	if err := models.DB.Self.Save(&u).Error; err != nil {
		log.Error("登录异常", err)
		return -5, nil, errors.New("登录异常")
	}

	// 生成token
	token, err := models.GenerateToken(&u, u.Id, 0)
	if err != nil {
		log.Error("生成token失败", err)
		return -6, nil, errors.New("生成token失败")
	}

	//组装响应数据
	resp := &dto.UserLoginResp{}
	resp.Id = u.Id
	resp.Username = u.Username
	resp.Realname = u.Realname
	resp.Email = u.Email
	resp.PhoneCode = u.PhoneCode
	resp.Phone = u.Phone
	resp.Avatar = utils.DisplayImg(u.Avatar)
	resp.Language = u.Language
	resp.Token = token

	return 1, resp, nil
}

// 修改密码
func (s *UserService) SetPassword(dto dto.Password) (int, error) {
	entity := new(models.User)
	user, err := entity.GetById(dto.Uid)
	if err != nil {
		log.Error("用户不存在", err)
		return -5, errors.New("用户不存在")
	}
	if !utils.PasswordVerify(strings.Trim(dto.Password, " "), user.Password) {
		log.Error("原始密码错误", err)
		return -6, errors.New("原始密码错误")
	}

	//获取加密后密码
	password, _ := utils.PasswordHash(dto.NewPassword)
	user.Password = password
	err = models.DB.Self.Save(user).Error
	if err != nil {
		log.Error("密码更新失败", err)
		return -7, errors.New("密码更新失败")
	}
	return 1, nil
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
func (s *UserService) Save(entity models.User) (int, error) {
	stat := 1

	if entity.Id != "" {
		//检测数据是否存在
		user, err := entity.GetById(entity.Id)
		if err != nil {
			log.Error("用户不存在 "+entity.Id, err)
			return -4, errors.New("用户不存在")
		}
		//密码为空不做修改
		if entity.Password != "" {
			password, _ := utils.PasswordHash(entity.Password)
			user.Password = password
		}
		user.Realname = entity.Realname
		user.Email = entity.Email
		user.PhoneCode = entity.PhoneCode
		user.Phone = entity.Phone
		user.Avatar = entity.Avatar
		user.Sex = entity.Sex
		user.UpdateUser = "1"               //修改人
		user.UpdateTime = utils.Timestamp() //修改时间
		err = user.UpdateById()
		if err != nil {
			log.Error("用户更新 "+entity.Id, err)
			return -5, errors.New("用户更新失败")
		}
	} else {
		entity.Id = utils.UniqueId()
		password, _ := utils.PasswordHash(entity.Password)
		entity.Password = password
		entity.Language = "cn" //默认中文
		entity.Status = 1
		entity.CreateUser = "1"               //添加人
		entity.CreateTime = utils.Timestamp() //添加时间
		err := entity.Add()
		if err != nil {
			log.Error("用户添加失败", err)
			return -2, errors.New("用户添加失败")
		}
	}
	return stat, nil
}
