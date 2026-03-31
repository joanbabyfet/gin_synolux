package repository

import (
	"gin-synolux/dto"
	"gin-synolux/models"

	"github.com/jinzhu/gorm"
)

//只管 DB
type FeedbackRepo struct {
	db *gorm.DB
}

func NewFeedbackRepo(db *gorm.DB) *FeedbackRepo {
	return &FeedbackRepo{db: db}
}

// 列表
func (r *FeedbackRepo) List(query dto.FeedbackQuery) ([]*models.Feedback, int64, error) {
	qs := r.db.Model(&models.Feedback{}).Where("delete_time = ?", 0)

	var list []*models.Feedback
	var count int64

	// count
	if query.Count {
		if err := qs.Count(&count).Error; err != nil {
			return nil, 0, err
		}
	}

	qs = qs.Order("create_time DESC")

	// 分页保护
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		qs = qs.Offset(offset).Limit(query.PageSize)
	} else if query.Limit > 0 {
		qs = qs.Limit(query.Limit)
	}

	if err := qs.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	
	return list, count, nil
}

// 创建
func (r *FeedbackRepo) Create(m *models.Feedback) error {
	return r.db.Create(m).Error
}

func (r *FeedbackRepo) WithTx(tx *gorm.DB) *FeedbackRepo {
	return &FeedbackRepo{db: tx}
}