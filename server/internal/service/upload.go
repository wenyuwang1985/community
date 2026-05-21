package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/wenyuwang1985/community/pkg/snowflake"
)

const (
	uploadDir     = "assets/uploads"
	maxFileSize   = 5 * 1024 * 1024 // 5MB
	maxFileCount  = 9
)

var allowedExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

type UploadService struct {
	baseURL string
}

func NewUploadService(baseURL string) *UploadService {
	// 确保上传目录存在
	_ = os.MkdirAll(uploadDir, 0755)
	return &UploadService{baseURL: baseURL}
}

func (s *UploadService) SaveFiles(files []*multipart.FileHeader) ([]string, error) {
	if len(files) > maxFileCount {
		return nil, fmt.Errorf("一次最多上传 %d 张图片", maxFileCount)
	}

	var urls []string
	for _, file := range files {
		if file.Size > maxFileSize {
			return nil, fmt.Errorf("文件 %s 超过 5MB 限制", file.Filename)
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			return nil, fmt.Errorf("不支持的文件格式: %s", ext)
		}

		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("打开文件失败: %w", err)
		}
		defer src.Close()

		filename := fmt.Sprintf("%d%s", snowflake.Generate(), ext)
		dstPath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return nil, fmt.Errorf("创建文件失败: %w", err)
		}

		_, err = io.Copy(dst, src)
		dst.Close()
		if err != nil {
			return nil, fmt.Errorf("保存文件失败: %w", err)
		}

		urls = append(urls, s.baseURL+"/uploads/"+filename)
	}

	return urls, nil
}
