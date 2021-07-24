package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StorageSeverceImpl struct {
	ServerConfigPath string // 服务配置文件
	blkMgr           Storage
	config           *Config
	server           *http.Server
}

func NewStorageServiceImpl() *StorageSeverceImpl {
	s := &StorageSeverceImpl{
		ServerConfigPath: "./storage.yaml",
	}
	// 1. 读取配置
	s.config = &Config{}
	s.config.LoadConf(s.ServerConfigPath)

	// 2. 初始化块管理器
	s.blkMgr = NewStorageProxy(s.config)

	// 3. 注册http服务
	router := s.newRouter()
	s.server = s.newHttpServer(router)

	// 4. 注册优雅退出
	s.RegisterOnShutdown(s.blkMgr.OnExit)

	return s
}

// RegisterOnShutdown 注册优雅退出函数
func (s *StorageSeverceImpl) RegisterOnShutdown(f func()) {
	s.server.RegisterOnShutdown(f)
}

func (s *StorageSeverceImpl) newHttpServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    s.config.GetPort(),
		Handler: router,
	}
}
func (s *StorageSeverceImpl) ListenAndServe() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()
	s.waitForExit()
}

func (s *StorageSeverceImpl) waitForExit() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-ctx.Done():
			if err := s.server.Shutdown(context.Background()); err != nil {
				logrus.Error(err)
			}
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
