package upload

// 上传文件的工具库

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/pkg/util"
)

type FileType int

const TypeImage FileType = iota + 1

func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext) // 去掉后缀
	fileName = util.EncodeMD5(fileName)       // 经过MD5加密后返回新的文件名

	return fileName + ext
}

// 获取文件后缀名，返回的后缀包括点
func GetFileExt(name string) string {
	return path.Ext(name)
}

func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

// 检查路径是否存在
func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst) // 获取文件的描述信息

	return os.IsNotExist(err)
}

// 检查文件权限是否足够
func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}

	}

	return false
}

// 检查文件大小是否超过最大值
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}

	return false
}

// 检查文件权限
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

// 创建保存上传的文件
func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}

	return nil
}

// 保存上传的文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open() // 打开上传的文件
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst) // 创建目标地址的文件对象
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src) // 文件拷贝
	return err
}
