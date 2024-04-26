package models

import (
	"gin-synolux/dto"
	"reflect"

	"github.com/spf13/viper"
)

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user
type User struct {
	Id           string `gorm:"primary_key;size(32);default();description(ID)" json:"id"`
	Origin       int8   `gorm:"default(0);null;description(注册来源 1=H5 2=PC)" json:"origin"`
	Username     string `gorm:"size(40);default();null;index;description(帐号)" json:"username"`
	Password     string `gorm:"size(60);default();null;description(密码)" json:"-"` //密码不输出
	Avatar       string `gorm:"size(100);default();null;description(头像)" json:"avatar"`
	Realname     string `gorm:"size(50);default();null;index;description(姓名)" json:"realname"`
	Sex          int8   `gorm:"default(1);null;description(性别 0=女 1=男)" json:"sex"`
	Email        string `gorm:"size(100);default();null;index;description(信箱)" json:"email"`
	PhoneCode    string `gorm:"size(5);default();null;index;description(手机号国码)" json:"phone_code"`
	Phone        string `gorm:"size(20);default();null;index;description(手机号)" json:"phone"`
	Address      string `gorm:"size(100);default();null;description(地址)" json:"address"`
	Salt         string `gorm:"size(128);default();null;description(加密钥匙)" json:"salt"`
	RoleId       int    `gorm:"default(0);null;description(角色)" json:"role_id"`
	RegIp        string `gorm:"size(15);default();null;description(注册ip)" json:"reg_ip"`
	LoginTime    int    `gorm:"default(0);null;description(最后登录时间)" json:"login_time"`
	LoginIp      string `gorm:"size(15);default();null;description(最后登录IP)" json:"login_ip"`
	LoginCountry string `gorm:"size(2);default();null;description(最后登录国家)" json:"login_country"`
	Language     string `gorm:"size(10);default();null;description(语言)" json:"language"`
	Status       int8   `gorm:"default(1);null;description(状态: 0=禁用 1=启用)" json:"status"`
	CreateTime   int    `gorm:"default(0);null;description(創建時間)" json:"create_time"`
	CreateUser   string `gorm:"size(32);default(0);null;description(創建人)" json:"create_user"`
	UpdateTime   int    `gorm:"default(0);null;description(修改時間)" json:"update_time"`
	UpdateUser   string `gorm:"size(32);default(0);null;description(修改人)" json:"update_user"`
	DeleteTime   int    `gorm:"default(0);null;description(刪除時間)" json:"delete_time"`
	DeleteUser   string `gorm:"size(32);default(0);null;description(刪除人)" json:"delete_user"`
}

func (m *User) TableName() string {
	return viper.GetString("db.prefix") + "user"
}

// 获取全部列表
func (m *User) All(query dto.UserQuery) (list []*User) {
	qs := DB.Self.Model(new(User))
	qs = qs.Where("delete_time = ?", 0) //未删除
	qs.Order("create_time desc").Find(&list)
	return list
}

// 获取分页列表
func (m *User) PageList(query dto.UserQuery) ([]*User, int64) {
	qs := DB.Self.Model(new(User))
	qs = qs.Where("delete_time = ?", 0) //未删除
	//总条数
	var count int64
	qs.Count(&count)
	var list []*User
	if count > 0 {
		offset := (query.Page - 1) * query.PageSize
		qs.Order("create_time desc").Limit(query.PageSize).Offset(offset).Find(&list)
	}
	if reflect.ValueOf(list).IsNil() {
		list = make([]*User, 0) //赋值为空切片[]
	}
	return list, count
}

// 获取单条
func (m *User) GetById(id string) (v *User, err error) {
	v = &User{}
	d := DB.Self.Where("delete_time = ?", 0).Where("id = ?", id).First(&v)
	if d.Error != nil {
		return nil, d.Error
	}
	return v, nil
}

// 单条添加
func (m *User) Add() error {
	return DB.Self.Create(&m).Error
}

// 更新
func (m *User) UpdateById() error {
	return DB.Self.Save(m).Error
}

// 删除
func (m *User) DeleteById(id string) error {
	m.Id = id
	return DB.Self.Delete(m).Error
}
