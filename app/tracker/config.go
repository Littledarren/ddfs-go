package main

import "fmt"

type Config struct {
	Port        int32
	BlkSize     int32
	StorageList []string
	BackupNum   int
}

var (
	DefaultConf = &Config{
		Port:        8888,
		BlkSize:     4 * 1024,
		StorageList: []string{"127.0.0.1:8080"},
		BackupNum:   1,
	}
)

func (c *Config) GetPort() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *Config) LoadConf(path string) error {
	*c = *DefaultConf
	return nil
}
