package main

import (
	"ddfs-go/internal/comm/utils"
	"mime/multipart"

	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	TrackerFile = "tracker.json"
)

type Tracker struct {
	CliTable  map[string][]*StorageCli    `json:"cli_table,omitempty"`  // hash=>存储块的storage
	FileTable map[string]*DistFileLocator `json:"file_table,omitempty"` // filename->hash列表
	CliList   []*StorageCli               `json:"-"`                    // 存储可用的storage列表
	cliLock   sync.Mutex                  `json:"-"`
	fileLock  sync.Mutex                  `json:"-"`
	config    *Config                     `json:"-"`
}

func (t *Tracker) Do(list []*StorageCli, f func(cli *StorageCli) error) error {
	var wg sync.WaitGroup
	var err error
	var once sync.Once
	for i := range list {
		i := i
		wg.Add(1)
		go func() {
			ierr := f(list[i])
			if ierr != nil {
				once.Do(func() {
					err = ierr
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return err
}

// AllocateCli 分配n个存储
func (t *Tracker) AllocateCli(n int) []*StorageCli {
	t.cliLock.Lock()
	defer t.cliLock.Unlock()
	// 小于n，无法分配
	if len(t.CliList) < n {
		return nil
	}
	rand.Shuffle(len(t.CliList), func(i, j int) {
		t.CliList[i], t.CliList[j] = t.CliList[j], t.CliList[i]
	})
	// 块分配到哪个storage? 随机
	return t.CliList[:n]
}

// getClient 根据hash，获取一个存储该块的客户端
func (t *Tracker) getClient(hash string) *StorageCli {
	t.cliLock.Lock()
	defer t.cliLock.Unlock()
	item := t.CliTable[hash]
	if len(item) == 0 {
		return nil
	}
	// 调度算法: 随机
	cli := item[rand.Intn(len(item))]
	return cli
}

// GetFileSize 读取文件大小
func (t *Tracker) GetFileSize(filename string) (int, error) {
	t.fileLock.Lock()
	defer t.fileLock.Unlock()
	dfl := t.FileTable[filename]
	if dfl == nil {
		return 0, fmt.Errorf("file not exist")
	}
	return dfl.Size, nil
}

// GetFileReader 获取一个文件的reader
func (t *Tracker) GetFileReader(filename string) (io.Reader, error) {
	t.fileLock.Lock()
	defer t.fileLock.Unlock()
	dfl := t.FileTable[filename]
	if dfl == nil {
		return nil, fmt.Errorf("file not exist")
	}
	dfl.Reset()
	f := func(p []byte) (int, error) {
		if len(p) == 0 {
			return 0, io.ErrShortBuffer
		}
		dfl.Lock.Lock()
		defer dfl.Lock.Unlock()

		hash, start, err := dfl.GetLocator()
		if err != nil {
			return 0, err
		}
		buffer, err := t.getClient(hash).Get(hash)
		if err != nil {
			return 0, err
		}
		buffer = buffer[start:]
		cnt := 0
		for i := range p {
			if i >= len(buffer) || dfl.IsEOF() == io.EOF {
				// len(p) > len(buffer)
				// or eof
				break // 本次读取完毕
			}
			p[i] = buffer[i]
			cnt++
			dfl.HasRead(1)
		}
		return cnt, dfl.IsEOF()
	}
	return ReaderFunc(f), nil
}

func (t *Tracker) SaveFile(file *multipart.FileHeader) error {
	t.fileLock.Lock()
	defer t.fileLock.Unlock()

	dfl := t.FileTable[file.Filename]
	if dfl == nil {
		// 文件不存在
		t.FileTable[file.Filename] = &DistFileLocator{
			BlkSize: int(t.config.BlkSize),
		}
	}
	t.FileTable[file.Filename].Lock.Lock()
	defer t.FileTable[file.Filename].Lock.Unlock()
	return t.updateFile(file)
}
func (t *Tracker) DeleteFile(fileName string) error {
	t.fileLock.Lock()
	defer t.fileLock.Unlock()

	dfl := t.FileTable[fileName]
	if dfl == nil {
		return nil
	}
	for _, blk := range dfl.HashList {
		err := t.Do(t.CliTable[blk], func(cli *StorageCli) error {
			return cli.Delete(blk)
		})
		if err != nil {
			logrus.Error(err)
		}
	}
	delete(t.FileTable, fileName)
	return nil
}

func (t *Tracker) updateFile(file *multipart.FileHeader) error {
	dfl := t.FileTable[file.Filename]
	dfl.Size = int(file.Size) // 更新size
	buffer := make([]byte, t.config.BlkSize)
	// 共有这么多个块需要写入
	blkNum := (dfl.Size + dfl.BlkSize - 1) / dfl.BlkSize // 文件块数量
	f, err := file.Open()
	if err != nil {
		return err
	}
	for i := 0; i < blkNum; i++ {
		// 读取文件块
		cnt, ierr, n := 0, error(nil), 0
		for ierr == nil && cnt < len(buffer) {
			n, ierr = f.Read(buffer[cnt:])
			cnt += n
		}
		if ierr != nil && ierr != io.EOF {
			// 读取文件发生错误
			return ierr
		}
		hash := utils.Md5s(buffer)
		if len(dfl.HashList) > i && dfl.HashList[i] == hash {
			// 说明文件块没有变化，跳过
			continue
		}
		if len(dfl.HashList) <= i {
			// 文件扩长
			dfl.HashList = append(dfl.HashList, hash)
		} else {
			// 文件修改
			oldHash := dfl.HashList[i]
			dfl.HashList[i] = hash
			// 删除旧hash
			if t.Do(t.CliTable[oldHash], func(cli *StorageCli) error {
				return cli.Delete(oldHash)
			}) != nil {
				// 删除失败咋整。没办法，会产生碎片垃圾
				logrus.Error(err)
			}
		}

		cliList := t.CliTable[hash]
		if len(cliList) == 0 {
			// 说明没存过, 分配三个cli
			cliList = t.AllocateCli(t.config.BackupNum)
		}
		t.CliTable[hash] = cliList
		if err := t.Do(cliList, func(cli *StorageCli) error {
			return cli.Post(hash, buffer)
		}); err != nil {
			logrus.Error(err)
			// 写失败咋整。没办法，犯错重写吧？
			return err
		}
	}
	dfl.HashList = dfl.HashList[:blkNum]
	return nil
}
