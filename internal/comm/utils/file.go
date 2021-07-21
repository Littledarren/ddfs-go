package utils

import "os"

// IsFileExist 检查文件是否存在
func IsFileExist(path string) error {
	f, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) || f.IsDir() {
		return err
	}
	return nil
}
