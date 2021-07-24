package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	s := NewStorageServiceImpl()
	s.ListenAndServe()
}

func init() {
	logrus.SetReportCaller(true)
}
