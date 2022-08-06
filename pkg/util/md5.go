package util

import (
	"crypto/md5"
	"encoding/hex"
)

// 用于格式化文件名
// 信息摘要算法——保存文件名的信息摘要
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
