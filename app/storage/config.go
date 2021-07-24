package main

import "fmt"

type Config struct {
	Root          string
	BlkNumPerFile int32  // 一个文件下能有多少块
	FilePrefix    string // 文件前缀
	BlkSize       int32  // 块大小 B
	Port          int32
}

var (
	DefaultConf *Config = &Config{
		Root:          "data",
		BlkNumPerFile: 100,
		FilePrefix:    "blk_",
		BlkSize:       4 * 1024, // 4 KB
		Port:          8000,
	}
)

func (c *Config) GetPort() string {
	return fmt.Sprintf(":%d", c.Port)
}

// LoadConf 读取配置
func (c *Config) LoadConf(path string) error {
	*c = *DefaultConf
	return nil
}
