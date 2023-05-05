module ddfs-go/internal/tracker

go 1.16

replace ddfs-go => ../../../ddfs-go

replace ddfs-go/internal/comm/errno => ../../../ddfs-go/internal/comm/errno

replace ddfs-go/internal/comm/utils => ../../../ddfs-go/internal/comm/utils

require (
	ddfs-go/internal/comm/utils v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.9.0
	github.com/sirupsen/logrus v1.8.1
	github.com/ugorji/go v1.1.7 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
