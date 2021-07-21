package main

// BlkItem  存储一个块的具体位置与引用次数
type BlkItem struct {
	Index    int    `json:"index"`
	FileName string `json:"file_name"` // 实际存储的文件名
	Offset   int64  `json:"offset"`    // 块偏移量
	RefCount int64  `json:"ref_count"` // 引用计数
}

// Table 保存hash=>块位置的映射
type Table map[string]*BlkItem
