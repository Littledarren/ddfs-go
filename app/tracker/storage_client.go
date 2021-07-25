package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// StorageCli  想Storage请求
type StorageCli struct {
	Target string `json:"target"` // 请求目标
}

// Get 获取块
func (c *StorageCli) Get(hash string) ([]byte, error) {
	// 1. 设置请求
	reqURL, err := c.getURL(hash)
	if err != nil {
		return nil, err
	}

	// 2. 发起请求
	rsp, err := http.Get(reqURL)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rsp.Body.Close()
	// 3. 读取返回
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return body, nil
}

// Post 上传
func (c *StorageCli) Post(hash string, blk []byte) error {
	reqURL, err := c.getURL(hash)
	if err != nil {
		return err
	}

	contentType := "application/ddfs-go"
	rsp, err := http.Post(reqURL, contentType, bytes.NewReader(blk))
	if err != nil {
		logrus.Error(err)
		return err
	}
	rsp.Body.Close()
	return nil
}

func (c *StorageCli) Delete(hash string) error {
	reqURL, err := c.getURL(hash)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		logrus.Error(err)
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	rsp.Body.Close()
	return nil
}

// getURL 因为restful，所以url是公用的
func (c *StorageCli) getURL(hash string) (string, error) {
	reqURL, err := url.ParseRequestURI(fmt.Sprintf("http://%s/blk", c.Target))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	reqURL.Query().Add("hash", hash)
	return reqURL.String(), nil
}
