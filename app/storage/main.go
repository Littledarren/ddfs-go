package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	blkMgr Storage
	// 做优雅退出
	gctx  = context.Background()
	done  = make(chan struct{})
	count = 1
)

func NewHTTPServer(conf *Config) *http.Server {
	// 绑定路由
	router := NewRouter()
	router.TrustedProxies = append(router.TrustedProxies, "255.255.255.255")
	return &http.Server{
		Addr:    conf.GetPort(),
		Handler: router,
	}
}

// Run 监听执行
func Run(s *http.Server) {
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()
}

func main() {
	ctx, cancel := signal.NotifyContext(gctx, syscall.SIGINT, syscall.SIGTERM)
	gctx = ctx

	blkMgr = NewStorageProxy(GetConf())
	s := NewHTTPServer(GetConf())
	Run(s)

	defer cancel()
	for {
		select {
		case <-done:
			count--
			if count == 0 {
				if err := s.Shutdown(context.Background()); err != nil {
					logrus.Error(err)
				}
				return
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func init() {
	logrus.SetReportCaller(true)
}
