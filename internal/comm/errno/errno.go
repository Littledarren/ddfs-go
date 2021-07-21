package errno

import (
	"errors"
)

const (
	RetInvalidParam = 450 // 非法参数

	RetSetFailed = 550 // Set failed
	RetGetFailed = 551 // Get failed
)

var (
	ErrBlkNotExist      = errors.New("blk not exist")       // 块不存在
	ErrBlkBroken        = errors.New("blk file broke")      // 文件损坏
	ErrBlkSizeCheckFail = errors.New("blk size check fail") // 块大小校验失败
	ErrInvalidParam     = errors.New("invalid params")      // 参数错误
)
