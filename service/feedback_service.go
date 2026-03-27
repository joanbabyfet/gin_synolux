package service

import (
	"gin-synolux/models"
	"gin-synolux/utils"

	"github.com/lexkong/log"
	"github.com/thedevsaddam/govalidator"
)

type FeedbackService struct {
}

// 保存
func (s *FeedbackService) Save(data models.Feedback, isAdmin bool) (error) {
	// 参数校验
	rules := govalidator.MapData{
		"name": []string{"required"},
		"mobile": []string{"required"},
		"email": []string{"required"},
		"content": []string{"required"},
	}
	messages := govalidator.MapData{
		"name": []string{"required:name 不能为空"},
		"mobile": []string{"required:mobile 不能为空"},
		"email": []string{"required:email 不能为空"},
		"content": []string{"required:content 不能为空"},
	}

	if err := utils.ValidateStruct(&data, rules, messages); err != nil {
		return NewServiceError(-1, err.Error())
	}

	// 开启事务
	tx := models.DB.Self.Begin()
	var err error

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	// 新增
	data.CreateUser = "1"
	data.CreateTime = utils.Timestamp()

	if err = tx.Create(&data).Error; err != nil {
		log.Error("反馈添加失败", err)
		return NewServiceError(-4, "反馈添加失败")
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		log.Error("事务提交失败", err)
		return NewServiceError(-5, "系统错误")
	}

	return nil
}
