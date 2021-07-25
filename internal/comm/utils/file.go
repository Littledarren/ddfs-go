package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// IsFileExist 检查文件是否存在
func IsFileExist(path string) error {
	f, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) || f.IsDir() {
		return err
	}
	return nil
}

// Md5s 计算块hash
func Md5s(blk []byte) string {
	r := md5.Sum(blk)
	return hex.EncodeToString(r[:])
}
