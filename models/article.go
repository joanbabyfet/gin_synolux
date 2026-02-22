package models

import (
	"gin-synolux/dto"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user
type Article struct {
	Id         int    `gorm:"primary_key;auto_increment;default();description(ID)" json:"id"`
	Catid      int    `gorm:"default(0);null;description(分類id)" json:"catid"`
	Title      string `gorm:"size(50);default();null;index;description(标题)" json:"title"`
	Info       string `gorm:"default();null;description(简介)" json:"info"`
	Content    string `gorm:"type(text);null;description(內容)" json:"content"`
	Img        string `gorm:"size(100);default();null;description(图片)" json:"img"`
	Author     string `gorm:"size(30);default();null;index;description(作者)" json:"author"`
	Extra      string `gorm:"default();null;index;description(扩展)" json:"extra"`
	Sort       int16  `gorm:"default(0);null;description(排序: 数字小的排前面)" json:"sort"`
	Status     int8   `gorm:"default(1);null;description(状态: 0=禁用 1=启用)" json:"status"`
	CreateTime int    `gorm:"default(0);null;description(創建時間)" json:"create_time"`
	CreateUser string `gorm:"size(32);default(0);null;description(創建人)" json:"create_user"`
	UpdateTime int    `gorm:"default(0);null;description(修改時間)" json:"update_time"`
	UpdateUser string `gorm:"size(32);default(0);null;description(修改人)" json:"update_user"`
	DeleteTime int    `gorm:"default(0);null;description(刪除時間)" json:"delete_time"`
	DeleteUser string `gorm:"size(32);default(0);null;description(刪除人)" json:"delete_user"`
}

func (m *Article) TableName() string {
	return viper.GetString("db.prefix") + "article"
}

// 获取分页列表
func (m *Article) List(query dto.ArticleQuery) ([]*Article, int64) {
	qs := DB.Self.Model(new(Article))
	qs = qs.Debug().Where("delete_time = ?", 0) //未删除
	if !reflect.ValueOf(&query.Status).IsNil() {
		qs = qs.Where("status = ?", query.Status)
	}
	if query.Catid != 0 {
		qs = qs.Where("catid = ?", query.Catid)
	}
	if len(query.Catids) > 0 {
		qs = qs.Where("catid IN ?", query.Catids)
	}
	if query.Title != "" {
		qs = qs.Where("title LIKE ?", "%"+query.Title+"%")
	}
	
	var list []*Article
	var count int64

	if query.Count { //是否返回总条数
		qs.Count(&count)
	}
	qs = qs.Order("create_time DESC")

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		qs = qs.Limit(query.PageSize).Offset(offset)
	} else if query.Limit > 0 {
		qs = qs.Limit(query.Limit)
	}
	qs.Find(&list)

	return list, count
}

// 获取单条
func (m *Article) GetById(id int) (v *Article, err error) {
	v = &Article{}
	d := DB.Self.Where("delete_time = ?", 0).Where("id = ?", id).First(&v)
	if d.Error != nil {
		return nil, d.Error
	}
	return v, nil
}

// 获取单条(支持事务)
func (m *Article) GetByIdTx(tx *gorm.DB, id int) (v *Article, err error) {
	v = &Article{}
	d := tx.Where("delete_time = ?", 0).Where("id = ?", id).First(&v)
	if d.Error != nil {
		return nil, d.Error
	}
	return v, nil
}

// 单条添加(支持事务)
func (m *Article) Add(tx ...*gorm.DB) error {
	db := DB.Self
    if len(tx) > 0 && tx[0] != nil {
        db = tx[0] // 使用事务对象
    }
    return db.Create(m).Error
}

// 更新(支持事务)
func (m *Article) UpdateById(tx ...*gorm.DB) error {
	db := DB.Self
    if len(tx) > 0 && tx[0] != nil {
        db = tx[0] // 使用事务对象
    }
    return db.Save(m).Error
}

// 删除
func (m *Article) DeleteById(id int, tx ...*gorm.DB) error {
	m.Id = id

    db := DB.Self
    if len(tx) > 0 && tx[0] != nil {
        db = tx[0] // 使用事务对象
    }

    return db.Delete(m).Error
}
