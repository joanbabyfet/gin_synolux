package repository

import (
	"gin-synolux/dto"
	"gin-synolux/models"

	"github.com/jinzhu/gorm"
)

//只管 DB
type ArticleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

// 列表
func (r *ArticleRepo) List(query dto.ArticleQuery) ([]*models.Article, int64, error) {
	qs := r.db.Model(&models.Article{}).Where("delete_time = ?", 0)

	if query.Status != nil {
		qs = qs.Where("status = ?", *query.Status)
	}
	if query.Catid != 0 {
		qs = qs.Where("catid = ?", query.Catid)
	}
	if len(query.Catids) > 0 {
		qs = qs.Where("catid IN (?)", query.Catids)
	}
	if query.Title != "" {
		qs = qs.Where("title LIKE ?", "%"+query.Title+"%")
	}

	var list []*models.Article
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
func (r *ArticleRepo) GetByID(id int) (*models.Article, error) {
	var article models.Article
	err := r.db.Where("id = ? AND delete_time = ?", id, 0).First(&article).Error
	return &article, err
}

func (r *ArticleRepo) WithTx(tx *gorm.DB) *ArticleRepo {
	return &ArticleRepo{db: tx}
}

// 创建
func (r *ArticleRepo) Create(m *models.Article) error {
	return r.db.Create(m).Error
}

// 更新 (会忽略 0)
// func (r *ArticleRepo) Update(m *models.Article) error {
// 	//GORM 在 Updates(struct) 时会忽略 0 / "" / false 这些零值字段
// 	return r.db.Debug().Model(m).Updates(m).Error
// }

// 更新 (不再用 Updates(struct))
func (r *ArticleRepo) Update(id int, data map[string]interface{}) error {
	return r.db.Model(&models.Article{}).
		Where("id = ?", id).
		Updates(data).Error
}

// 软删除
func (r *ArticleRepo) Delete(tx *gorm.DB, id int) error {
	return tx.Model(&models.Article{}).Where("id = ?", id).Update("delete_time", gorm.Expr("UNIX_TIMESTAMP()")).Error
}

// 是否存在
func (r *ArticleRepo) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Article{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}