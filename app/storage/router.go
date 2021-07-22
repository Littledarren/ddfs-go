package main

import (
	"ddfs-go/internal/comm/errno"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	ret := gin.Default()

	ret.POST("/blk", func(c *gin.Context) {
		blk, err := ioutil.ReadAll(c.Request.Body)
		if err != nil || len(blk) > int(GetConf().BlkSize) {
			c.String(errno.RetInvalidParam, "非法参数")
			return
		}
		if err = blkMgr.Set(c.Query("hash"), blk); err != nil {
			c.String(errno.RetInvalidParam, fmt.Sprintf("内部错误: %+v", err))
			return
		}
		c.String(200, "")
	})
	ret.DELETE("/blk", func(c *gin.Context) {
		if err := blkMgr.Del(c.Query("hash")); err != nil {
			c.String(errno.RetSetFailed, fmt.Sprintf("内部错误: %+v", err))
			return
		}
		c.String(200, "")
	})
	ret.GET("/blk", func(c *gin.Context) {
		blk, err := blkMgr.Get(c.Query("hash"))
		if err != nil {
			c.String(200, fmt.Sprintf("内部错误: %+v", err))
			return
		}
		c.Data(200, "blk", blk)
	})
	ret.GET("/", func(c *gin.Context) {
		c.String(200, "blk storage works fine😅")
	})
	return ret
}
