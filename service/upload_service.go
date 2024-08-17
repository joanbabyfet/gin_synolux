package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type UploadService struct {
}

// 获取全部列表
func (s *UploadService) Upload(ctx *gin.Context, f *multipart.FileHeader, dir string, thumb_w int, thumb_h int) (int, interface{}, error) {
	file_url := viper.GetString("file_url")
	stat := 1

	//文件大小校验
	upload_max_size := viper.GetString("upload_max_size")
	max_size, _ := strconv.ParseInt(upload_max_size, 10, 64)
	if f.Size > max_size*1024*1024 {
		log.Error("您上传的文件过大,最大值为"+upload_max_size+"MB", nil)
		return -2, nil, errors.New("您上传的文件过大,最大值为" + upload_max_size + "MB")
	}

	//文件后缀过滤
	ext := path.Ext(f.Filename) //输出.jpg
	allow_ext_map := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	if _, ok := allow_ext_map[ext]; !ok {
		log.Error("文件格式不正确", nil)
		return -3, nil, errors.New("文件格式不正确")
	}

	//创建目录
	upload_dir := viper.GetString("upload_dir")
	dir_num := time.Now().Format("20060102") //输出 20240404/
	err := os.MkdirAll(upload_dir+"/"+dir+"/"+dir_num, os.FileMode(0775))
	if err != nil {
		log.Error("创建目录失败", err)
		return -4, nil, errors.New("创建目录失败")
	}

	//构造文件名
	source := rand.NewSource(time.Now().UnixNano()) //这里用系统时间毫秒值当种子值
	r := rand.New(source)
	rand_num := fmt.Sprintf("%d", r.Intn(9999)+1000) //获取1000-9999随机数
	hash_name := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + rand_num))
	file_name := fmt.Sprintf("%x", hash_name) + ext //文件名 例 cf386af3f37962ad3769054f68d7a049.jpg

	path := upload_dir + "/" + dir + "/" + dir_num + "/" + file_name
	filelink := file_url + "/" + dir + "/" + dir_num + "/" + file_name
	ctx.SaveUploadedFile(f, path) //名称与 c.GetFile("xxx") 一致

	//生成缩略图
	if thumb_w > 0 || thumb_h > 0 {
		src, err := imaging.Open(path)
		if err != nil {
			log.Error("开启缩略图失败", err)
			return -5, nil, errors.New("开启缩略图失败")
		}
		dsc := imaging.Resize(src, thumb_w, thumb_h, imaging.Lanczos)

		//构造缩略图文件名
		source := rand.NewSource(time.Now().UnixNano()) //这里用系统时间毫秒值当种子值
		r := rand.New(source)
		rand_num := fmt.Sprintf("%d", r.Intn(9999)+1000) //获取1000-9999随机数
		hash_name := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + rand_num))
		file_name = fmt.Sprintf("%x", hash_name) + ext //文件名 例 cf386af3f37962ad3769054f68d7a049.jpg
		path = upload_dir + "/" + dir + "/" + dir_num + "/" + file_name
		filelink = file_url + "/" + dir + "/" + dir_num + "/" + file_name

		err = imaging.Save(dsc, path)
		if err != nil {
			log.Error("生成缩略图失败", err)
			return -6, nil, errors.New("生成缩略图失败")
		}
	}

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["realname"] = f.Filename
	resp["filename"] = dir_num + "/" + file_name
	resp["filelink"] = filelink
	return stat, resp, nil
}
