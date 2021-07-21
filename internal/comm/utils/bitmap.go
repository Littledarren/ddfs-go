package utils

import "sync"

/*
	参考了网上的代码，简单位图实现
*/

type BitMap struct {
	Bits []uint64 `json:"bits"`
	lock sync.Mutex
}

func NewBitMap() *BitMap {
	bits := make([]uint64, 0)
	ret := &BitMap{
		Bits: bits,
	}
	return ret
}

// FirstZero 第一个0的索引
func (b *BitMap) FirstZero() int {
	var ret int
	var tmp uint64

	for _, p := range b.Bits {
		if p != ^uint64(0) {
			tmp = p
			break
		}
		ret += 64
	}
	i := firstZero(tmp)
	if i == -1 {
		return ret + 65
	}
	return ret + i
}

func (b *BitMap) FindAndSet() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	index := b.FirstZero()
	b.Set(index)
	return index
}

func (b *BitMap) Set(loc int) {
	loc--
	i := loc / 64
	mask := uint64(1) << (loc % 64)
	if i >= len(b.Bits) {
		b.Bits = append(b.Bits, make([]uint64, i+1-len(b.Bits))...)
	}
	b.Bits[i] |= mask
}

func (b *BitMap) Unset(loc int) {
	loc--
	i := loc / 64
	mask := uint64(1) << (loc % 64)
	if i >= len(b.Bits) {
		return
	}
	b.Bits[i] &= ^mask
}

func firstZero(word uint64) int {
	for i := 0; i < 64; i++ {
		if word&(1<<i) == 0 {
			return (i + 1)
		}
	}
	return -1
}
