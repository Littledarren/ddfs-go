module ddfs-go/internal/tracker

go 1.16

replace ddfs-go => ../../../ddfs-go

replace ddfs-go/internal/comm/errno => ../../../ddfs-go/internal/comm/errno

replace ddfs-go/internal/comm/utils => ../../../ddfs-go/internal/comm/utils

require (
	ddfs-go/internal/comm/utils v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.2
	github.com/sirupsen/logrus v1.8.1
)
