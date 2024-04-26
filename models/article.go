package models

import (
	"gin-synolux/dto"
	"reflect"

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

// 获取全部列表
func (m *Article) All(query dto.ArticleQuery) (list []*Article) {
	qs := DB.Self.Model(new(Article))
	qs = qs.Where("delete_time = ?", 0) //未删除
	if query.Limit != 0 {
		qs = qs.Limit(query.Limit)
	}
	if !reflect.ValueOf(&query.Status).IsNil() {
		qs = qs.Where("status = ?", query.Status)
	}
	if query.Catid != 0 {
		qs = qs.Where("catid = ?", query.Catid)
	}
	if len(query.Catids) > 0 {
		qs = qs.Where("catid IN ?", query.Catids)
	}
	if len(query.Title) > 1 {
		qs = qs.Where("title LIKE ?", "%"+query.Title+"%")
	}
	qs.Order("create_time desc").Find(&list)
	return list
}

// 获取分页列表
func (m *Article) PageList(query dto.ArticleQuery) ([]*Article, int64) {
	qs := DB.Self.Model(new(Article))
	qs = qs.Where("delete_time = ?", 0) //未删除
	if !reflect.ValueOf(&query.Status).IsNil() {
		qs = qs.Where("status = ?", query.Status)
	}
	if query.Catid != 0 {
		qs = qs.Where("catid = ?", query.Catid)
	}
	if len(query.Catids) > 0 {
		qs = qs.Where("catid IN ?", query.Catids)
	}
	if len(query.Title) > 1 {
		qs = qs.Where("title LIKE ?", "%"+query.Title+"%")
	}
	//总条数
	var count int64
	qs.Count(&count)
	var list []*Article
	if count > 0 {
		offset := (query.Page - 1) * query.PageSize
		qs.Order("create_time desc").Limit(query.PageSize).Offset(offset).Find(&list)
	}
	if reflect.ValueOf(list).IsNil() {
		list = make([]*Article, 0) //赋值为空切片[]
	}
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

// 单条添加
func (m *Article) Add() error {
	return DB.Self.Create(&m).Error
}

// 更新
func (m *Article) UpdateById() error {
	return DB.Self.Save(m).Error
}

// 删除
func (m *Article) DeleteById(id int) error {
	m.Id = id
	return DB.Self.Delete(m).Error
}
