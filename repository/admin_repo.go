package repository

import (
	"errors"
	"gin-synolux/dto"
	"gin-synolux/models"

	"github.com/jinzhu/gorm"
)

//只管 DB
type AdminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

// 列表
func (r *AdminRepo) List(query dto.UserQuery) ([]*models.Admin, int64, error) {
	qs := r.db.Model(&models.Admin{}).Where("delete_time = ?", 0)

	if query.Status != nil {
		qs = qs.Where("status = ?", *query.Status)
	}

	var list []*models.Admin
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
func (r *AdminRepo) GetByID(id string) (*models.Admin, error) {
	var user models.Admin
	err := r.db.Where("id = ? AND delete_time = ?", id, 0).First(&user).Error
	return &user, err
}

func (r *AdminRepo) WithTx(tx *gorm.DB) *AdminRepo {
	return &AdminRepo{db: tx}
}

// 创建
func (r *AdminRepo) Create(m *models.Admin) error {
	return r.db.Create(m).Error
}

// 更新 (不再用 Updates(struct))
func (r *AdminRepo) Update(id string, data map[string]interface{}) error {
	return r.db.Model(&models.Admin{}).
		Where("id = ?", id).
		Updates(data).Error
}

// 软删除
func (r *AdminRepo) Delete(tx *gorm.DB, id string) error {
	return tx.Model(&models.Admin{}).Where("id = ?", id).Update("delete_time", gorm.Expr("UNIX_TIMESTAMP()")).Error
}

// 是否存在
func (r *AdminRepo) ExistsByID(id string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Admin{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *AdminRepo) GetByUsername(username string) (*models.Admin, error) {
	var user models.Admin

	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		// 关键：区分“查不到”和“真正错误”
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没找到不算错误
		}
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) ExistsByUsername(username string) (bool, error) {
	var count int64

	err := r.db.Model(&models.Admin{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}