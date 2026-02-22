package service

import (
	"encoding/json"
	"fmt"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/utils"
	"strconv"
	"time"

	"github.com/lexkong/log"
	"github.com/thedevsaddam/govalidator"
)

type ArticleService struct {
}

// 获取分页列表
func (s *ArticleService) List(query dto.ArticleQuery) ([]*models.Article, int64) {
	entity := new(models.Article) //new实例化
	return entity.List(query)
}

// 获取详情
func (s *ArticleService) GetById(id int) (*models.Article, error) {
	// 参数验证
	entity := models.Article{Id: id}
    rules := govalidator.MapData{
        "id": []string{"required"},
    }
    messages := govalidator.MapData{
        "id": []string{"required:id 不能为空"},
    }
    if err := utils.ValidateStruct(&entity, rules, messages); err != nil {
        return nil, NewServiceError(-1, err.Error())
    }

	//先查缓存
	cache_key := fmt.Sprintf("article:detail:%d", id)
	val, err := utils.Redis.Get(cache_key).Result()
	if err == nil {
		var info models.Article
		if jsonErr := json.Unmarshal([]byte(val), &info); jsonErr == nil {
			return &info, nil
		}
	}

	//缓存未命中，查库
	info, err := entity.GetById(id)
	if err != nil {
		log.Error("文章不存在 "+strconv.Itoa(id), err)
		return nil, NewServiceError(-1, "文章不存在")
	}

	bytes, err := json.Marshal(info)
	if err == nil {
		_ = utils.Redis.Set(cache_key, bytes, time.Hour).Err()
	}

	return info, nil
}

// 保存
func (s *ArticleService) Save(data models.Article, isAdmin bool) (error) {
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

	tx := models.DB.Self.Begin() //开启事务
	defer func() {
		//防止紧急停止
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if data.Id > 0 {
		//检测数据是否存在
		entity := new(models.Article)
		info, err := entity.GetByIdTx(tx, data.Id)
		if err != nil {
			tx.Rollback() //手动回滚事务
			log.Error("文章不存在 "+strconv.Itoa(data.Id), err)
			return NewServiceError(-2, "文章不存在")
		}
		info.Catid = data.Catid
		info.Title = data.Title
		info.Info = data.Info
		info.Content = data.Content
		info.Author = data.Author
		info.Status = data.Status
		info.UpdateUser = "1"               //修改人
		info.UpdateTime = utils.Timestamp() //修改时间
		if err := info.UpdateById(tx); err != nil {
			tx.Rollback() //手动回滚事务
			log.Error("文章更新 "+strconv.Itoa(data.Id), err)
			return NewServiceError(-3, "文章更新失败")
		}
	} else {
		data.Status = 1
		data.CreateUser = "1"               //添加人
		data.CreateTime = utils.Timestamp() //添加时间
		err := data.Add(tx)
		if err != nil {
			tx.Rollback() //手动回滚事务
			log.Error("文章添加", err)
			return NewServiceError(-4, "文章添加失败")
		}
	}
	
	//手动提交事务
	if err := tx.Commit().Error; err != nil {
		return NewServiceError(-5, err.Error())
    }
	
	if data.Id > 0 {
		cache_key := fmt.Sprintf("article:detail:%d", data.Id)
		_ = utils.Redis.Del(cache_key).Err()
	}

	// 后台操作日志
	if isAdmin {
		if data.Id > 0 {
			log.Info(fmt.Sprintf("更新文章 id=%d", data.Id))
		} else {
			log.Info("添加文章")
		}
	}

	return nil
}

// 软删除
func (s *ArticleService) DeleteById(id int, isAdmin bool) (error) {
	// 参数验证
	entity := models.Article{Id: id}
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
		log.Error("文章不存在 "+strconv.Itoa(id), err)
		return NewServiceError(-2, "文章不存在")
	}

	info.DeleteUser = "1"               //修改人
	info.DeleteTime = utils.Timestamp() //修改时间
	err = info.UpdateById(tx)
	if err != nil {
		tx.Rollback() //手动回滚事务
		log.Error("文章删除 "+strconv.Itoa(id), err)
		return NewServiceError(-3, "文章删除失败")
	}

	//手动提交事务
	if err := tx.Commit().Error; err != nil {
		return NewServiceError(-4, err.Error())
    }

	cache_key := fmt.Sprintf("article:detail:%d", id)
	_ = utils.Redis.Del(cache_key).Err()

	// 后台操作日志
	if isAdmin {
		log.Info(fmt.Sprintf("删除文章 id=%d", id))
	}

	return nil
}

//变更状态
func (s *ArticleService) ChangeStatus(id int, status int, isAdmin bool) (error) {
	// 参数验证
	entity := models.Article{Id: id}
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
		log.Error("文章不存在 "+strconv.Itoa(id), err)
		return NewServiceError(-2, "文章不存在")
	}

	info.Status = int8(status)
	info.UpdateUser = "1"               //修改人
	info.UpdateTime = utils.Timestamp() //修改时间
	err = info.UpdateById(tx)
	if err != nil {
		tx.Rollback() //手动回滚事务
		log.Error("文章禁用 "+strconv.Itoa(id), err)
		return NewServiceError(-3, "文章禁用失败")
	}

	//手动提交事务
	if err := tx.Commit().Error; err != nil {
		return NewServiceError(-4, err.Error())
    }

	cache_key := fmt.Sprintf("article:detail:%d", id)
	_ = utils.Redis.Del(cache_key).Err()

	// 后台操作日志
	if isAdmin {
		log.Info(fmt.Sprintf("修改文章状态 id=%d status=%d", id, status))
	}

	return nil
}