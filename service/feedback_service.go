package service

import (
	"gin-synolux/common"
	"gin-synolux/db"
	"gin-synolux/models"
	"gin-synolux/repository"

	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/thedevsaddam/govalidator"
)

type FeedbackService struct {
	repo *repository.FeedbackRepo
}

func NewFeedbackService(db *gorm.DB) *FeedbackService {
	if db == nil {
		panic("db is nil (service)")
	}

	return &FeedbackService{
		repo: repository.NewFeedbackRepo(db),
	}
}

// 保存
func (s *FeedbackService) Save(data models.Feedback, isAdmin bool) (error) {
	// 参数校验
	rules := govalidator.MapData{
		"name":    []string{"required"},
		"mobile":  []string{"required"},
		"email":   []string{"required"},
		"content": []string{"required"},
	}
	messages := govalidator.MapData{
		"name":    []string{"required:name 不能为空"},
		"mobile":  []string{"required:mobile 不能为空"},
		"email":   []string{"required:email 不能为空"},
		"content": []string{"required:content 不能为空"},
	}

	if err := common.ValidateStruct(&data, rules, messages); err != nil {
		return common.NewError(-1, err.Error())
	}

	tx := db.DB.Self.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	repo := s.repo.WithTx(tx)

	now := common.Timestamp()
	data.CreateUser = "0"
	data.CreateTime = now

	// 统一走 repo
	if err := repo.Create(&data); err != nil {
		tx.Rollback()
		log.Error("反馈添加失败", err)
		return common.NewError(-4, "反馈添加失败")
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("事务提交失败", err)
		return common.NewError(-5, "系统错误")
	}

	return nil
}
