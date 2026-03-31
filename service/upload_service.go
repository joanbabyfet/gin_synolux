package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gin-synolux/dto"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type UploadService struct {
}

func NewUploadService(db *gorm.DB) *UploadService {
	if db == nil {
		panic("db is nil (service)")
	}

	return &UploadService{}
}

func (s *UploadService) Upload(f *multipart.FileHeader, dir string, thumbW, thumbH int) (*dto.UploadResp, error) {

	fileURL := viper.GetString("file_url")
	uploadDir := viper.GetString("upload_dir")

	// ===== 文件大小校验 =====
	maxSizeStr := viper.GetString("upload_max_size")
	maxSize, _ := strconv.ParseInt(maxSizeStr, 10, 64)
	if f.Size > maxSize*1024*1024 {
		return nil, errors.New("文件过大，最大 " + maxSizeStr + "MB")
	}

	// ===== 后缀校验 =====
	ext := strings.ToLower(path.Ext(f.Filename))
	allowExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	}
	if !allowExt[ext] {
		return nil, errors.New("文件格式不支持")
	}

	// ===== 创建目录 =====
	dirDate := time.Now().Format("20060102")
	fullDir := filepath.Join(uploadDir, dir, dirDate)

	if err := os.MkdirAll(fullDir, 0775); err != nil {
		log.Error("创建目录失败", err)
		return nil, errors.New("创建目录失败")
	}

	// ===== 原图保存 =====
	fileName := genFileName(ext)
	savePath := filepath.Join(fullDir, fileName)

	// 👉 这里不再用 ctx
	src, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	fileLink := fmt.Sprintf("%s/%s/%s/%s", fileURL, dir, dirDate, fileName)

	// ===== 缩略图 =====
	if thumbW > 0 || thumbH > 0 {
		img, err := imaging.Open(savePath)
		if err != nil {
			return nil, errors.New("读取图片失败")
		}

		thumb := imaging.Resize(img, thumbW, thumbH, imaging.Lanczos)

		fileName = genFileName(ext)
		savePath = filepath.Join(fullDir, fileName)
		fileLink = fmt.Sprintf("%s/%s/%s/%s", fileURL, dir, dirDate, fileName)

		if err = imaging.Save(thumb, savePath); err != nil {
			return nil, errors.New("生成缩略图失败")
		}
	}

	// ===== 返回 DTO =====
	return &dto.UploadResp{
		RealName: f.Filename,
		FileName: dirDate + "/" + fileName,
		FileLink: fileLink,
	}, nil
}

func (s *UploadService) GetFilePath(filename string) (string, string, error) {
	if filename == "" {
		return "", "", errors.New("文件名不能为空")
	}

	// 防止路径穿越攻击（非常重要）
	if strings.Contains(filename, "..") {
		return "", "", errors.New("非法文件路径")
	}

	uploadDir := viper.GetString("upload_dir")

	// 默认 image 目录（你也可以做成参数）
	fullPath := filepath.Join(uploadDir, "image", filename)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", "", errors.New("文件不存在")
	}

	// 取真实文件名（用于下载显示）
	realName := filepath.Base(filename)

	return fullPath, realName, nil
}

func genFileName(ext string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := fmt.Sprintf("%d", r.Intn(9000)+1000)

	hash := md5.Sum([]byte(time.Now().Format("20060102150405") + randNum))
	return fmt.Sprintf("%x", hash) + ext
}