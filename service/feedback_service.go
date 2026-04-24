package service

import (
	"gin-synolux/common"
	"gin-synolux/db"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/repository"

	"github.com/jinzhu/gorm"
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
func (s *FeedbackService) Save(req *dto.FeedbackSaveReq, isAdmin bool) (error) {	
	//开启事务
	tx := db.DB.Self.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 用事务 repo
	repo := s.repo.WithTx(tx)
	now := common.Timestamp()

	data := models.Feedback{
		Name:    req.Name,
		Mobile:  req.Mobile,
		Email:   req.Email,
		Content: req.Content,
		CreateUser: req.CreateUser,
		CreateTime: now,
	}

	// 统一走 repo
	if err := repo.Create(&data); err != nil {
		tx.Rollback()
		common.Log.Error("反馈添加失败", err)
		return common.NewError(-4, "反馈添加失败")
	}

	if err := tx.Commit().Error; err != nil {
		common.Log.Error("事务提交失败", err)
		return common.NewError(-5, "系统错误")
	}

	return nil
}
