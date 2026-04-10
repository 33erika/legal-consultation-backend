package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var allowedExtensions = map[string]bool{
	"pdf":  true,
	"doc":  true,
	"docx": true,
	"xlsx": true,
	"xls":  true,
	"jpg":  true,
	"jpeg": true,
	"png":  true,
}

// UploadFile 保存上传的文件
func UploadFile(file *multipart.FileHeader, entityType, entityID string) (string, error) {
	// 检查文件大小
	maxSize := viper.GetInt64("upload.max_size")
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds limit: %d bytes", maxSize)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("file type not allowed: %s", ext)
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s_%s%s", entityType, uuid.New().String(), filepath.Ext(file.Filename))

	// 创建上传目录
	uploadPath := viper.GetString("upload.path")
	dir := filepath.Join(uploadPath, entityType, entityID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 保存文件
	dst := filepath.Join(dir, filename)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return dst, nil
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	return os.Remove(filePath)
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
