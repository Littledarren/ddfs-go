package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// TrackerServiceImpl  提供文件上传下载服务
type TrackerServiceImpl struct {
	ServerConfigPath string
	server           *http.Server
	tracker          *Tracker
	config           *Config
}

// RegisterOnShutdown 注册优雅退出函数
func (s *TrackerServiceImpl) RegisterOnShutdown(f func()) {
	s.server.RegisterOnShutdown(f)
}

func NewTrackerServiceImpl() *TrackerServiceImpl {
	s := &TrackerServiceImpl{
		ServerConfigPath: "./tracker.yaml",
	}
	s.config = &Config{}
	if err := s.config.LoadConf(s.ServerConfigPath); err != nil {
		logrus.Error(err)
	}

	s.tracker = NewTrackerProxy(s.config)

	s.server = &http.Server{
		Addr:    s.config.GetPort(),
		Handler: s.newRouter(),
	}
	s.RegisterOnShutdown(s.tracker.OnExit)

	return s
}

func (s *TrackerServiceImpl) ListenAndServe() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()
	s.waitForExit()
}

func (s *TrackerServiceImpl) waitForExit() {
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
