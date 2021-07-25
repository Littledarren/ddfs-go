package main

import "github.com/gin-gonic/gin"

func (s *TrackerServiceImpl) newRouter() *gin.Engine {
	r := gin.Default()
	// 上传，下载，删除, 获取
	r.POST("/file", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(800, err.Error())
			return
		}
		if err := s.tracker.SaveFile(file); err != nil {
			c.String(801, err.Error())
			return
		}
		c.String(200, "upload success")

	})
	r.GET("/file", func(c *gin.Context) {
		fileName := c.Query("file")
		size, err := s.tracker.GetFileSize(fileName)
		if err != nil {
			c.String(200, err.Error())
			return
		}
		fr, err := s.tracker.GetFileReader(fileName)
		if err != nil {
			c.String(200, err.Error())
			return
		}
		c.DataFromReader(200, int64(size), "file", fr, nil)
	})
	r.DELETE("/file", func(c *gin.Context) {
		fileName := c.Query("file")
		if err := s.tracker.DeleteFile(fileName); err != nil {
			c.String(302, err.Error())
			return
		}
		c.String(200, "delete success")
	})
	r.GET("/", func(c *gin.Context) {
		c.String(200, "🤺🤺🤺go go go💃💃💃")
	})

	return r
}
