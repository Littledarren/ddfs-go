module ddfs-go/internal/storage

go 1.16

replace ddfs-go => ../../../ddfs-go

replace ddfs-go/internal/comm/errno => ../../../ddfs-go/internal/comm/errno

replace ddfs-go/internal/comm/utils => ../../../ddfs-go/internal/comm/utils

require (
	ddfs-go/internal/comm/errno v0.0.0-00010101000000-000000000000
	ddfs-go/internal/comm/utils v0.0.0-00010101000000-000000000000
	github.com/0xAX/notificator v0.0.0-20191016112426-3962a5ea8da1 // indirect
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/codegangsta/gin v0.0.0-20171026143024-cafe2ce98974 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/mattn/go-shellwords v1.0.12 // indirect
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
