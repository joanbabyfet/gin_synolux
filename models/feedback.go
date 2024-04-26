package models

import (
	"gin-synolux/dto"
	"reflect"

	"github.com/spf13/viper"
)

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user
type Feedback struct {
	Id         int    `gorm:"primary_key;auto_increment;default();description(ID)" json:"id"`
	Name       string `gorm:"size(50);default();null;index;description(姓名)" json:"name"`
	Mobile     string `gorm:"size(20);default();null;index;description(手机号)" json:"mobile"`
	Email      string `gorm:"size(100);default();null;index;description(信箱)" json:"email"`
	Content    string `gorm:"type(text);null;description(內容)" json:"content"`
	CreateTime int    `gorm:"default(0);null;description(創建時間)" json:"create_time"`
	CreateUser string `gorm:"size(32);default(0);null;description(創建人)" json:"create_user"`
	UpdateTime int    `gorm:"default(0);null;description(修改時間)" json:"update_time"`
	UpdateUser string `gorm:"size(32);default(0);null;description(修改人)" json:"update_user"`
	DeleteTime int    `gorm:"default(0);null;description(刪除時間)" json:"delete_time"`
	DeleteUser string `gorm:"size(32);default(0);null;description(刪除人)" json:"delete_user"`
}

func (m *Feedback) TableName() string {
	return viper.GetString("db.prefix") + "feedback"
}

// 获取全部列表
func (m *Feedback) All(query dto.FeedbackQuery) (list []*Feedback) {
	qs := DB.Self.Model(new(Feedback))
	qs = qs.Where("delete_time = ?", 0) //未删除
	qs.Order("create_time desc").Find(&list)
	return list
}

// 获取分页列表
func (m *Feedback) PageList(query dto.FeedbackQuery) ([]*Feedback, int64) {
	qs := DB.Self.Model(new(Feedback))
	qs = qs.Where("delete_time = ?", 0) //未删除
	//总条数
	var count int64
	qs.Count(&count)
	var list []*Feedback
	if count > 0 {
		offset := (query.Page - 1) * query.PageSize
		qs.Order("create_time desc").Limit(query.PageSize).Offset(offset).Find(&list)
	}
	if reflect.ValueOf(list).IsNil() {
		list = make([]*Feedback, 0) //赋值为空切片[]
	}
	return list, count
}

// 获取单条
func (m *Feedback) GetById(id int) (v *Feedback, err error) {
	v = &Feedback{}
	d := DB.Self.Where("delete_time = ?", 0).Where("id = ?", id).First(&v)
	if d.Error != nil {
		return nil, d.Error
	}
	return v, nil
}

// 单条添加
func (m *Feedback) Add() error {
	return DB.Self.Create(&m).Error
}

// 更新
func (m *Feedback) UpdateById() error {
	return DB.Self.Save(m).Error
}

// 删除
func (m *Feedback) DeleteById(id int) error {
	m.Id = id
	return DB.Self.Delete(m).Error
}
