package main

import (
	"ddfs-go/internal/comm/errno"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter() *gin.Engine {
	ret := gin.Default()

	ret.POST("/blk", func(c *gin.Context) {
		blk, err := ioutil.ReadAll(c.Request.Body)
		logrus.Infof("%+v", blk)
		if err != nil || len(blk) > int(GetConf().BlkSize) {
			logrus.Info(len(blk))
			c.String(errno.RetInvalidParam, "非法参数")
			return
		}
		if err = blkMgr.Set(c.Query("hash"), blk); err != nil {
			c.String(errno.RetInvalidParam, "内部错误: %w", err)
			return
		}
		c.JSON(200, nil)
	})
	ret.DELETE("/blk", func(c *gin.Context) {
		if err := blkMgr.Del(c.Query("hash")); err != nil {
			c.String(errno.RetSetFailed, "内部错误: %w", err)
			return
		}
		c.JSON(200, nil)
	})
	ret.GET("/blk", func(c *gin.Context) {
		blk, err := blkMgr.Get(c.Query("hash"))
		if err != nil {
			logrus.Error(err)
			c.JSON(errno.RetGetFailed, nil)
			return
		}
		c.Data(200, "blk", blk)
	})
	ret.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	return ret
}
