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
	GlobalConfig *Config = &Config{
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

func GetConf() *Config {
	return GlobalConfig
}
