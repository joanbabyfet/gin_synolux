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
func (s *AdService) GetById(req dto.AdDetailReq) (*models.Ad, error) {
	//先查缓存
	cache_key := fmt.Sprintf("ad:detail:%d", req.ID)
	val, err := common.Redis.Get(cache_key).Result()
	if err == nil {
		var info models.Ad
		if jsonErr := json.Unmarshal([]byte(val), &info); jsonErr == nil {
			return &info, nil
		}
	}

	//缓存未命中，查库
	info, err := s.repo.GetByID(req.ID)
	if err != nil {
		common.Log.Error("广告不存在 "+strconv.Itoa(req.ID), err)
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
func (s *AdService) Save(req *dto.AdSaveReq, isAdmin bool) error {
	isUpdate := req.ID > 0

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
		exists, err := repo.ExistsByID(req.ID)
		if err != nil {
			tx.Rollback()
			common.Log.Error("查询广告失败 id="+strconv.Itoa(req.ID), err)
			return common.NewError(-2, "查询失败")
		}
		if !exists {
			tx.Rollback()
			return common.NewError(-2, "广告不存在")
		}

		// ===== 更新字段 =====
		updateData := map[string]interface{}{
			"catid":       req.Catid,
			"title":       req.Title,
			"img":         req.Img,
			"url":         req.Url,
			"sort":        req.Sort,
			"status":      req.Status,
			"update_user": req.UpdateUser,
			"update_time": now,
		}

		// ===== 更新（走 repo）=====
		if err = repo.Update(req.ID, updateData); err != nil {
			tx.Rollback()
			common.Log.Error("广告更新 "+strconv.Itoa(req.ID), err)
			return common.NewError(-3, "广告更新失败")
		}
	} else {
		data := models.Ad{
			Catid:      req.Catid,
			Title:      req.Title,
			Img:        req.Img,
			Url:        req.Url,
			Sort:       int16(req.Sort),
			Status:     1, // 默认启用
			CreateUser: req.CreateUser,
			CreateTime: now,
		}

		// ===== 创建（走 repo）=====
		if err = repo.Create(&data); err != nil {
			tx.Rollback()
			common.Log.Error("广告添加失败", err)
			return common.NewError(-4, "广告添加失败")
		}
		req.ID = data.Id // 可选：回写 ID
	}

	// ===== 提交事务 =====
	if err = tx.Commit().Error; err != nil {
		common.Log.Error("事务提交失败", err)
		return common.NewError(-5, err.Error())
	}

	// ===== 清缓存 =====
	if isUpdate {
		cacheKey := fmt.Sprintf("ad:detail:%d", req.ID)
		_ = common.Redis.Del(cacheKey).Err()
	}

	// ===== 日志 =====
	if isAdmin {
		if isUpdate {
			common.Log.Info(fmt.Sprintf("更新广告 id=%d", req.ID))
		} else {
			common.Log.Info("添加广告")
		}
	}

	return nil
}

// 软删除
func (s *AdService) DeleteById(req dto.AdDeleteReq, isAdmin bool) error {
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

	exists, err := repo.ExistsByID(req.ID)
	if err != nil {
		tx.Rollback()
		common.Log.Error("查询广告失败 id="+strconv.Itoa(req.ID), err)
		return common.NewError(-2, "查询失败")
	}
	if !exists {
		tx.Rollback()
		return common.NewError(-2, "广告不存在")
	}

	//软删除
	now := common.Timestamp()
	data := map[string]interface{}{
		"delete_user": req.UserID,
		"delete_time": now,
	}

	//更新（删除）
	if err := repo.Update(req.ID, data); err != nil {
		tx.Rollback()
		common.Log.Error("广告删除失败 id="+strconv.Itoa(req.ID), err)
		return common.NewError(-3, "删除失败")
	}

	//提交事务
	if err := tx.Commit().Error; err != nil {
		common.Log.Error("事务提交失败", err)
		return common.NewError(-4, err.Error())
	}

	//删除缓存（事务成功后）
	cacheKey := fmt.Sprintf("ad:detail:%d", req.ID)
	if err := common.Redis.Del(cacheKey).Err(); err != nil {
		common.Log.Error("删除缓存失败", err)
	}

	//日志
	if isAdmin {
		common.Log.Infof("删除广告 id=%d", req.ID)
	}

	return nil
}

// 变更状态
func (s *AdService) ChangeStatus(req dto.AdChangeStatusReq, isAdmin bool) error {
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
	exists, err := repo.ExistsByID(req.ID)
	if err != nil {
		tx.Rollback()
		common.Log.Error("查询广告失败 id="+strconv.Itoa(req.ID), err)
		return common.NewError(-2, "查询失败")
	}
	if !exists {
		tx.Rollback()
		return common.NewError(-2, "广告不存在")
	}

	now := common.Timestamp()
	data := map[string]interface{}{
		"status":      int8(req.Status),
		"update_user": req.UserID,
		"update_time": now,
	}

	if err := repo.Update(req.ID, data); err != nil {
		tx.Rollback()
		common.Log.Error("修改广告状态失败 id="+strconv.Itoa(req.ID), err)
		return common.NewError(-3, "状态修改失败")
	}

	//手动提交事务
	if err := tx.Commit().Error; err != nil {
		common.Log.Error("事务提交失败", err)
		return common.NewError(-4, err.Error())
	}

	cacheKey := fmt.Sprintf("ad:detail:%d", req.ID)
	if err := common.Redis.Del(cacheKey).Err(); err != nil {
		common.Log.Error("删除缓存失败", err)
	}

	// 后台操作日志
	if isAdmin {
		common.Log.Infof("修改广告状态 id=%d status=%d", req.ID, req.Status)
	}

	return nil
}