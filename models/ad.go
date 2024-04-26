package models

import (
	"gin-synolux/dto"
	"reflect"

	"github.com/spf13/viper"
)

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user
type Ad struct {
	Id         int    `gorm:"primary_key;auto_increment;default();description(ID)" json:"id"`
	Catid      int    `gorm:"default(0);null;description(分類id)" json:"catid"`
	Title      string `gorm:"size(50);default();null;index;description(标题)" json:"title"`
	Img        string `gorm:"size(100);default();null;description(图片)" json:"img"`
	Url        string `gorm:"size(100);default();null;description(链接)" json:"url"`
	Sort       int16  `gorm:"default(0);null;description(排序: 数字小的排前面)" json:"sort"`
	Status     int8   `gorm:"default(1);null;description(状态: 0=禁用 1=启用)" json:"status"`
	CreateTime int    `gorm:"default(0);null;description(創建時間)" json:"create_time"`
	CreateUser string `gorm:"size(32);default(0);null;description(創建人)" json:"create_user"`
	UpdateTime int    `gorm:"default(0);null;description(修改時間)" json:"update_time"`
	UpdateUser string `gorm:"size(32);default(0);null;description(修改人)" json:"update_user"`
	DeleteTime int    `gorm:"default(0);null;description(刪除時間)" json:"delete_time"`
	DeleteUser string `gorm:"size(32);default(0);null;description(刪除人)" json:"delete_user"`
}

func (m *Ad) TableName() string {
	return viper.GetString("db.prefix") + "ad"
}

// 获取全部列表
func (m *Ad) All(query dto.AdQuery) (list []*Ad) {
	qs := DB.Self.Model(new(Ad))
	qs = qs.Where("delete_time = ?", 0) //未删除
	if !reflect.ValueOf(&query.Status).IsNil() {
		qs = qs.Where("status = ?", query.Status)
	}
	if query.Catid != 0 {
		qs = qs.Where("catid = ?", query.Catid)
	}
	qs.Order("create_time desc").Find(&list)
	return list
}

// 获取分页列表
func (m *Ad) PageList(query dto.AdQuery) ([]*Ad, int64) {
	qs := DB.Self.Model(new(Ad))
	qs = qs.Where("delete_time = ?", 0) //未删除
	if !reflect.ValueOf(&query.Status).IsNil() {
		qs = qs.Where("status = ?", query.Status)
	}
	if query.Catid != 0 {
		qs = qs.Where("catid = ?", query.Catid)
	}
	//总条数
	var count int64
	qs.Count(&count)
	var list []*Ad
	if count > 0 {
		offset := (query.Page - 1) * query.PageSize
		qs.Order("create_time desc").Limit(query.PageSize).Offset(offset).Find(&list)
	}
	if reflect.ValueOf(list).IsNil() {
		list = make([]*Ad, 0) //赋值为空切片[]
	}
	return list, count
}

// 获取单条
func (m *Ad) GetById(id int) (v *Ad, err error) {
	v = &Ad{}
	d := DB.Self.Where("delete_time = ?", 0).Where("id = ?", id).First(&v)
	if d.Error != nil {
		return nil, d.Error
	}
	return v, nil
}

// 单条添加
func (m *Ad) Add() error {
	return DB.Self.Create(&m).Error
}

// 更新
func (m *Ad) UpdateById() error {
	return DB.Self.Save(m).Error
}

// 删除
func (m *Ad) DeleteById(id int) error {
	m.Id = id
	return DB.Self.Delete(m).Error
}
