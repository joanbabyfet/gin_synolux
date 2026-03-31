package repository

import (
	"gin-synolux/dto"
	"gin-synolux/models"

	"github.com/jinzhu/gorm"
)

//只管 DB
type AdRepo struct {
	db *gorm.DB
}

func NewAdRepo(db *gorm.DB) *AdRepo {
	return &AdRepo{db: db}
}

// 列表
func (r *AdRepo) List(query dto.AdQuery) ([]*models.Ad, int64, error) {
	qs := r.db.Model(&models.Ad{}).Where("delete_time = ?", 0)

	if query.Status != nil {
		qs = qs.Where("status = ?", *query.Status)
	}
	if query.Catid != 0 {
		qs = qs.Where("catid = ?", query.Catid)
	}
	if len(query.Catids) > 0 {
		qs = qs.Where("catid IN (?)", query.Catids)
	}

	var list []*models.Ad
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

// 单条
func (r *AdRepo) GetByID(id int) (*models.Ad, error) {
	var Ad models.Ad
	err := r.db.Where("id = ? AND delete_time = ?", id, 0).First(&Ad).Error
	return &Ad, err
}

func (r *AdRepo) WithTx(tx *gorm.DB) *AdRepo {
	return &AdRepo{db: tx}
}

// 创建
func (r *AdRepo) Create(m *models.Ad) error {
	return r.db.Create(m).Error
}

// 更新 (不再用 Updates(struct))
func (r *AdRepo) Update(id int, data map[string]interface{}) error {
	return r.db.Model(&models.Ad{}).
		Where("id = ?", id).
		Updates(data).Error
}

// 软删除
func (r *AdRepo) Delete(tx *gorm.DB, id int) error {
	return tx.Model(&models.Ad{}).Where("id = ?", id).Update("delete_time", gorm.Expr("UNIX_TIMESTAMP()")).Error
}

// 是否存在
func (r *AdRepo) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Ad{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}