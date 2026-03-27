package service

import (
	"fmt"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/utils"
	"strconv"

	"github.com/lexkong/log"
	"github.com/thedevsaddam/govalidator"
)

type AdService struct {
}

// 获取分页列表
func (s *AdService) List(query dto.AdQuery) ([]*models.Ad, int64) {
	entity := new(models.Ad) //new实例化
	return entity.List(query)
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
    if err := utils.ValidateStruct(&entity, rules, messages); err != nil {
        return nil, NewServiceError(-1, err.Error())
    }

	info, err := entity.GetById(id)
	if err != nil {
		log.Error("广告不存在 "+strconv.Itoa(id), err)
		return nil, NewServiceError(-1, "广告不存在")
	}
	return info, nil
}

// 保存
func (s *AdService) Save(data models.Ad, isAdmin bool) (error) {
	rules := govalidator.MapData{
        "title": []string{"required"},
    }
    messages := govalidator.MapData{
        "title": []string{"required:title 不能为空"},
    }
    // 更新操作必须验证 ID
    if data.Id > 0 {
        rules["id"] = []string{"required"}
        messages["id"] = []string{"required:id 不能为空"}
    }

	if err := utils.ValidateStruct(&data, rules, messages); err != nil {
        return NewServiceError(-1, err.Error())
    }

	if data.Id > 0 {
		//检测数据是否存在
		entity := new(models.Ad)
		info, err := entity.GetById(data.Id)
		if err != nil {
			log.Error("广告不存在 "+strconv.Itoa(data.Id), err)
			return NewServiceError(-2, "广告不存在")
		}
		info.Catid = data.Catid
		info.Title = data.Title
		info.Status = data.Status
		info.UpdateUser = "1"               //修改人
		info.UpdateTime = utils.Timestamp() //修改时间
		err = info.UpdateById()
		if err != nil {
			log.Error("广告更新 "+strconv.Itoa(data.Id), err)
			return NewServiceError(-3, "广告更新")
		}
	} else {
		data.Status = 1
		data.CreateUser = "1"               //添加人
		data.CreateTime = utils.Timestamp() //添加时间
		err := data.Add()
		if err != nil {
			log.Error("广告添加失败", err)
			return NewServiceError(-4, "广告添加失败")
		}
	}
	// 后台操作日志
	if isAdmin {
		if data.Id > 0 {
			log.Info(fmt.Sprintf("更新广告 id=%d", data.Id))
		} else {
			log.Info("添加广告")
		}
	}

	return nil
}

// 软删除
func (s *AdService) DeleteById(id int, isAdmin bool) (error) {
	// 参数验证
	entity := models.Ad{Id: id}
    rules := govalidator.MapData{
        "id": []string{"required"},
    }
    messages := govalidator.MapData{
        "id": []string{"required:id 不能为空"},
    }
    if err := utils.ValidateStruct(&entity, rules, messages); err != nil {
        return NewServiceError(-1, err.Error())
    }

	tx := models.DB.Self.Begin() //开启事务
	defer func() {
		//防止紧急停止
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	//检测数据是否存在
	info, err := entity.GetByIdTx(tx, id)
	if err != nil {
		tx.Rollback() //手动回滚事务
		log.Error("广告不存在 "+strconv.Itoa(id), err)
		return NewServiceError(-2, "广告不存在")
	}

	info.DeleteUser = "1"               //修改人
	info.DeleteTime = utils.Timestamp() //修改时间
	err = info.UpdateById(tx)
	if err != nil {
		tx.Rollback() //手动回滚事务
		log.Error("广告删除 "+strconv.Itoa(id), err)
		return NewServiceError(-3, "广告删除失败")
	}

	//手动提交事务
	if err := tx.Commit().Error; err != nil {
		return NewServiceError(-4, err.Error())
    }

	// 后台操作日志
	if isAdmin {
		log.Info(fmt.Sprintf("删除广告 id=%d", id))
	}

	return nil
}

//变更状态
func (s *AdService) ChangeStatus(id int, status int, isAdmin bool) (error) {
	// 参数验证
	entity := models.Ad{Id: id}
    rules := govalidator.MapData{
        "id": []string{"required"},
    }
    messages := govalidator.MapData{
        "id": []string{"required:id 不能为空"},
    }
    if err := utils.ValidateStruct(&entity, rules, messages); err != nil {
        return NewServiceError(-1, err.Error())
    }

	tx := models.DB.Self.Begin() //开启事务
	defer func() {
		//防止紧急停止
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	//检测数据是否存在
	info, err := entity.GetById(id)
	if err != nil {
		tx.Rollback() //手动回滚事务
		log.Error("广告不存在 "+strconv.Itoa(id), err)
		return NewServiceError(-2, "广告不存在")
	}

	info.Status = int8(status)
	info.UpdateUser = "1"               //修改人
	info.UpdateTime = utils.Timestamp() //修改时间
	err = info.UpdateById(tx)
	if err != nil {
		tx.Rollback() //手动回滚事务
		log.Error("广告禁用 "+strconv.Itoa(id), err)
		return NewServiceError(-3, "广告禁用失败")
	}

	//手动提交事务
	if err := tx.Commit().Error; err != nil {
		return NewServiceError(-4, err.Error())
    }

	cache_key := fmt.Sprintf("ad:detail:%d", id)
	_ = utils.Redis.Del(cache_key).Err()

	// 后台操作日志
	if isAdmin {
		log.Info(fmt.Sprintf("修改广告状态 id=%d status=%d", id, status))
	}

	return nil
}