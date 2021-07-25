package main

import (
	"ddfs-go/internal/comm/utils"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// GenCli 生成cli
func (t *Tracker) GenCli(target string) {
	t.cliLock.Lock()
	defer t.cliLock.Unlock()
	t.CliList = append(t.CliList, &StorageCli{
		Target: target,
	})
}

// NewTrackerProxy 生成tracker
func NewTrackerProxy(config *Config) *Tracker {
	t := &Tracker{
		config:    config,
		CliTable:  map[string][]*StorageCli{},
		FileTable: map[string]*DistFileLocator{},
	}
	for _, target := range config.StorageList {
		t.GenCli(target)
	}
	if err := t.LoadFromFile(TrackerFile); err != nil {
		logrus.Fatal(err)
	}
	// 启动一个协程自动同步文件
	go func() {
		for {
			// 挺影响性能的
			time.Sleep(10 * time.Second)
			t.fileLock.Lock()
			t.cliLock.Lock()
			if err := t.SyncToFile(TrackerFile); err != nil {
				logrus.Error(err)
			}
			t.cliLock.Unlock()
			t.fileLock.Unlock()
		}
	}()
	return t
}

func (t *Tracker) OnExit() {
	t.cliLock.Lock()
	t.fileLock.Lock() // 不释放, 一直持有到进程结束
	if err := t.SyncToFile(TrackerFile); err != nil {
		logrus.Error(err)
	}
}

// SyncToFile 保存到本地
func (t *Tracker) SyncToFile(p string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}
func (t *Tracker) LoadFromFile(path string) error {
	t.cliLock.Lock()
	defer t.cliLock.Unlock()
	t.fileLock.Lock()
	defer t.fileLock.Unlock()
	if err := utils.IsFileExist(path); err != nil {
		// 没有生成相应的文件，说明是第一次生成
		return t.SyncToFile(path)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, t); err != nil {
		return err
	}
	return nil
}
