package service

import (
	"encoding/json"
	"fmt"
	"gin-synolux/common"
	"gin-synolux/db"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/repository"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/thedevsaddam/govalidator"
)

type AdService struct {
	repo *repository.AdRepo
}

func NewAdService(db *gorm.DB) *AdService {
	if db == nil {
		panic("db is nil (service)")
	}

	return &AdService{
		repo: repository.NewAdRepo(db),
	}
}

// 列表
func (s *AdService) List(query dto.AdQuery) (*dto.AdListResp, error) {
	list, count, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	return &dto.AdListResp{
		List:  list,
		Count: count,
	}, nil
}

// 获取详情
func (s *AdService) GetById(id int) (*models.Ad, error) {
	// 参数验证
	entity := models.Ad{Id: id}
	rules := govalidator.MapData{
		"id": []string{"required"},
	}
	messages := govalidator.MapData{
		"id": []string{"required:id 不能为空"},
	}
	if err := common.ValidateStruct(&entity, rules, messages); err != nil {
		return nil, common.NewError(-1, err.Error())
	}

	//先查缓存
	cache_key := fmt.Sprintf("ad:detail:%d", id)
	val, err := common.Redis.Get(cache_key).Result()
	if err == nil {
		var info models.Ad
		if jsonErr := json.Unmarshal([]byte(val), &info); jsonErr == nil {
			return &info, nil
		}
	}

	//缓存未命中，查库
	info, err := s.repo.GetByID(id)
	if err != nil {
		log.Error("广告不存在 "+strconv.Itoa(id), err)
		return nil, common.NewError(-1, "广告不存在")
	}

	//写缓存
	bytes, err := json.Marshal(info)
	if err == nil {
		_ = common.Redis.Set(cache_key, bytes, time.Hour).Err()
	}

	return info, nil
}

// 保存
func (s *AdService) Save(data models.Ad, isAdmin bool) error {
	rules := govalidator.MapData{
		"title": []string{"required"},
	}
	messages := govalidator.MapData{
		"title": []string{"required:title 不能为空"},
	}

	isUpdate := data.Id > 0

	// 更新操作必须验证 ID
	if isUpdate {
		rules["id"] = []string{"required"}
		messages["id"] = []string{"required:id 不能为空"}
	}

	if err := common.ValidateStruct(&data, rules, messages); err != nil {
		return common.NewError(-1, err.Error())
	}

	//开启事务
	tx := db.DB.Self.Begin()
	var err error
	defer func() {
		//防止紧急停止
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 用事务 repo
	repo := s.repo.WithTx(tx)
	now := common.Timestamp()

	if isUpdate {
		//检测数据是否存在
		exists, err := repo.ExistsByID(data.Id)
		if err != nil {
			tx.Rollback()
			log.Error("查询广告失败 id="+strconv.Itoa(data.Id), err)
			return common.NewError(-2, "查询失败")
		}
		if !exists {
			tx.Rollback()
			return common.NewError(-2, "广告不存在")
		}

		// ===== 更新字段 =====
		updateData := map[string]interface{}{
			"catid":       data.Catid,
			"title":       data.Title,
			"img":         data.Img,
			"url":         data.Url,
			"sort":        data.Sort,
			"status":      data.Status,
			"update_user": "1",
			"update_time": now,
		}

		// ===== 更新（走 repo）=====
		if err = repo.Update(data.Id, updateData); err != nil {
			tx.Rollback()
			log.Error("广告更新 "+strconv.Itoa(data.Id), err)
			return common.NewError(-3, "广告更新失败")
		}
	} else {
		data.Status = 1
		data.CreateUser = "1"
		data.CreateTime = common.Timestamp()

		// ===== 创建（走 repo）=====
		if err = repo.Create(&data); err != nil {
			tx.Rollback()
			log.Error("广告添加失败", err)
			return common.NewError(-4, "广告添加失败")
		}
	}

	// ===== 提交事务 =====
	if err = tx.Commit().Error; err != nil {
		log.Error("事务提交失败", err)
		return common.NewError(-5, err.Error())
	}

	// ===== 清缓存 =====
	if isUpdate {
		cacheKey := fmt.Sprintf("ad:detail:%d", data.Id)
		_ = common.Redis.Del(cacheKey).Err()
	}

	// ===== 日志 =====
	if isAdmin {
		if isUpdate {
			log.Info(fmt.Sprintf("更新广告 id=%d", data.Id))
		} else {
			log.Info("添加广告")
		}
	}

	return nil
}

// 软删除
func (s *AdService) DeleteById(id int, isAdmin bool) error {
	// 参数验证
	entity := models.Ad{Id: id}
	rules := govalidator.MapData{
		"id": []string{"required"},
	}
	messages := govalidator.MapData{
		"id": []string{"required:id 不能为空"},
	}
	if err := common.ValidateStruct(&entity, rules, messages); err != nil {
		return common.NewError(-1, err.Error())
	}

	tx := db.DB.Self.Begin() //开启事务
	defer func() {
		//防止紧急停止
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 用事务 repo
	repo := s.repo.WithTx(tx)

	exists, err := repo.ExistsByID(id)
	if err != nil {
		tx.Rollback()
		log.Error("查询广告失败 id="+strconv.Itoa(id), err)
		return common.NewError(-2, "查询失败")
	}
	if !exists {
		tx.Rollback()
		return common.NewError(-2, "广告不存在")
	}

	//软删除
	now := common.Timestamp()
	data := map[string]interface{}{
		"delete_user": "1",
		"delete_time": now,
	}

	//更新（删除）
	if err := repo.Update(id, data); err != nil {
		tx.Rollback()
		log.Error("广告删除失败 id="+strconv.Itoa(id), err)
		return common.NewError(-3, "删除失败")
	}

	//提交事务
	if err := tx.Commit().Error; err != nil {
		log.Error("事务提交失败", err)
		return common.NewError(-4, err.Error())
	}

	//删除缓存（事务成功后）
	cacheKey := fmt.Sprintf("ad:detail:%d", id)
	if err := common.Redis.Del(cacheKey).Err(); err != nil {
		log.Error("删除缓存失败", err)
	}

	//日志
	if isAdmin {
		log.Infof("删除广告 id=%d", id)
	}

	return nil
}

// 变更状态
func (s *AdService) ChangeStatus(id int, status int, isAdmin bool) error {
	// 参数验证
	entity := models.Ad{Id: id}
	rules := govalidator.MapData{
		"id": []string{"required"},
	}
	messages := govalidator.MapData{
		"id": []string{"required:id 不能为空"},
	}
	if err := common.ValidateStruct(&entity, rules, messages); err != nil {
		return common.NewError(-1, err.Error())
	}

	tx := db.DB.Self.Begin() //开启事务
	defer func() {
		//防止紧急停止
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 用事务 repo
	repo := s.repo.WithTx(tx)

	//检测数据是否存在
	exists, err := repo.ExistsByID(id)
	if err != nil {
		tx.Rollback()
		log.Error("查询广告失败 id="+strconv.Itoa(id), err)
		return common.NewError(-2, "查询失败")
	}
	if !exists {
		tx.Rollback()
		return common.NewError(-2, "广告不存在")
	}

	now := common.Timestamp()
	data := map[string]interface{}{
		"status":      int8(status),
		"update_user": "1",
		"update_time": now,
	}

	if err := repo.Update(id, data); err != nil {
		tx.Rollback()
		log.Error("修改广告状态失败 id="+strconv.Itoa(id), err)
		return common.NewError(-3, "状态修改失败")
	}

	//手动提交事务
	if err := tx.Commit().Error; err != nil {
		log.Error("事务提交失败", err)
		return common.NewError(-4, err.Error())
	}

	cacheKey := fmt.Sprintf("ad:detail:%d", id)
	if err := common.Redis.Del(cacheKey).Err(); err != nil {
		log.Error("删除缓存失败", err)
	}

	// 后台操作日志
	if isAdmin {
		log.Infof("修改广告状态 id=%d status=%d", id, status)
	}

	return nil
}