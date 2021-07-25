package main

import (
	"io"
	"sync"
)

// DistFileLocator 一个分布式文件定位器
type DistFileLocator struct {
	HashList []string   `json:"hash_list"`
	BlkSize  int        `json:"blk_size"` // 一个块的大小
	Offset   int        `json:"offset"`   // 偏移量 byte
	Size     int        `json:"size"`     // 文件大小 byte
	Lock     sync.Mutex `json:"-"`        // 读写都加锁
}

func (dfl *DistFileLocator) Reset() {
	dfl.Offset = 0
}
func (dfl *DistFileLocator) IsEOF() error {
	if dfl.Offset == dfl.Size {
		return io.EOF
	}
	return nil
}
func (dfl *DistFileLocator) GetLocator() (string, int, error) {
	i := dfl.Offset / dfl.BlkSize
	if i >= len(dfl.HashList) {
		return "", 0, io.EOF
	}
	return dfl.HashList[i], dfl.Offset % dfl.BlkSize, nil
}
func (dfl *DistFileLocator) HasRead(n int) {
	dfl.Offset += n
}

type ReaderFunc func(p []byte) (int, error)

func (r ReaderFunc) Read(p []byte) (int, error) {
	return r(p)
}
