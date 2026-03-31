package repository

import (
	"errors"
	"gin-synolux/dto"
	"gin-synolux/models"

	"github.com/jinzhu/gorm"
)

//只管 DB
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// 列表
func (r *UserRepo) List(query dto.UserQuery) ([]*models.User, int64, error) {
	qs := r.db.Model(&models.User{}).Where("delete_time = ?", 0)

	if query.Status != nil {
		qs = qs.Where("status = ?", *query.Status)
	}

	var list []*models.User
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
func (r *UserRepo) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND delete_time = ?", id, 0).First(&user).Error
	return &user, err
}

func (r *UserRepo) WithTx(tx *gorm.DB) *UserRepo {
	return &UserRepo{db: tx}
}

// 创建
func (r *UserRepo) Create(m *models.User) error {
	return r.db.Create(m).Error
}

// 更新 (不再用 Updates(struct))
func (r *UserRepo) Update(id string, data map[string]interface{}) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(data).Error
}

// 软删除
func (r *UserRepo) Delete(tx *gorm.DB, id string) error {
	return tx.Model(&models.User{}).Where("id = ?", id).Update("delete_time", gorm.Expr("UNIX_TIMESTAMP()")).Error
}

// 是否存在
func (r *UserRepo) ExistsByID(id string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User

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

func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
	var count int64

	err := r.db.Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}