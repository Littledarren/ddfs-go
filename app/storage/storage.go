package main

import (
	"ddfs-go/internal/comm/errno"
	"ddfs-go/internal/comm/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// StorageFile 用于存储文件系统的信息
	StorageFile = "lfs.json"
)

// Storage  对存储的抽象，能够根据blk hash获取对应的内容
type Storage interface {
	Set(key string, blk []byte) error // 保存块到持久化存储
	Get(key string) ([]byte, error)   // 从持久化存储读取块
	Del(key string) error             // 删除块
	OnExit()                          // 退出时执行的函数
}

type localFileStorage struct {
	Table     Table        `json:"table"`
	Config    *Config      `json:"-"`
	BitMap    utils.BitMap `json:"bitmap"`
	TableLock sync.Mutex   `json:"-"`
}

func (l *localFileStorage) OnExit() {
	l.TableLock.Lock() // 不释放, 一直持有到进程结束
	l.SyncToFile(path.Join(l.Config.Root, StorageFile))
}

// SyncToFile 保存到本地
func (l *localFileStorage) SyncToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	// 这里需要复制一份防止出现同时读写的情况
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (l *localFileStorage) LoadFromFile(path string) error {
	l.TableLock.Lock()
	defer l.TableLock.Unlock()
	if err := utils.IsFileExist(path); err != nil {
		// 没有生成相应的文件，说明是第一次生成
		return l.SyncToFile(path)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, l); err != nil {
		return err
	}
	return nil
}

func NewStorageProxy(config *Config) Storage {
	lfs := &localFileStorage{
		Config: config,
		Table:  make(Table),
	}
	if err := lfs.LoadFromFile(path.Join(config.Root, StorageFile)); err != nil {
		logrus.Fatal(err)
	}
	// 启动一个协程自动同步文件
	go func() {
		for {
			time.Sleep(5 * time.Second)
			lfs.TableLock.Lock()
			lfs.SyncToFile(path.Join(config.Root, StorageFile))
			lfs.TableLock.Unlock()
		}
	}()
	return lfs
}

// Get 从文件获取块
func (s *localFileStorage) Get(key string) ([]byte, error) {
	s.TableLock.Lock()
	defer s.TableLock.Unlock()
	blk := s.Table[key]
	if blk == nil {
		return nil, errno.ErrBlkNotExist
	}
	return s.readBlk(blk)
}

// Set 保存块到文件
func (s *localFileStorage) Set(key string, blk []byte) error {
	blkLen := len(blk)
	if blkLen > int(s.Config.BlkSize) {
		return errno.ErrBlkSizeCheckFail
	}
	s.TableLock.Lock()
	defer s.TableLock.Unlock()
	item := s.Table[key]
	if item != nil && item.RefCount != 0 {
		item.RefCount++ // 去重
		return nil
	}
	// 第一次添加
	s.Table[key] = &BlkItem{
		RefCount: 1,
	}
	item = s.Table[key]

	// 分配一个块
	item.Index = s.BitMap.FindAndSet()

	postfix := item.Index / int(s.Config.BlkNumPerFile)
	item.Offset = int64(((item.Index - 1) % int(s.Config.BlkNumPerFile)) * int(s.Config.BlkSize))
	item.FileName = path.Join(s.Config.Root, genFileName(s.Config.Root, postfix))
	f, err := os.OpenFile(item.FileName, os.O_WRONLY|os.O_CREATE, 0664)
	defer func() {
		if err != nil {
			item.RefCount--
			s.BitMap.Unset(item.Index)
			item.Index = 0
		}
	}()
	if err != nil {
		return err
	}
	defer f.Close()

	cnt, err := f.WriteAt(blk, item.Offset)
	if cnt != blkLen {
		return err
	}
	return nil
}

// Del 删除块
func (s *localFileStorage) Del(key string) error {
	// 如果块不存在，直接返回nil
	// 如果块存在，refCount--
	// 如果refCount = 0， bitmap标记为0
	s.TableLock.Lock()
	defer s.TableLock.Unlock()
	item := s.Table[key]
	if item == nil || item.RefCount == 0 {
		return nil
	}
	item.RefCount--
	if item.RefCount == 0 {
		s.BitMap.Unset(item.Index)
		delete(s.Table, key)
	}
	return nil
}

// readBlk 读取块
func (s *localFileStorage) readBlk(bi *BlkItem) ([]byte, error) {
	fp, err := os.Open(bi.FileName)
	defer fp.Close()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, s.Config.BlkSize)
	cnt, err := fp.ReadAt(buf, bi.Offset)
	if cnt != int(s.Config.BlkSize) && err != io.EOF {
		return nil, errno.ErrBlkBroken
	}
	return buf, nil
}

// genFileName 生成文件名
func genFileName(path string, postfix int) string {
	return fmt.Sprintf("%s%d", path, postfix)
}
