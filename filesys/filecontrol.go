// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesys

import (
	"github.com/vrbyte/yusha/config"
	"log"
	"net/http"
	"os"
	"strings"
)

// FileControlInter 文件系统抽象接口模型
type FileControlInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	initFileSys()
	checkUrl(string, int) bool
	initIndexHtmlUrl()
}

// 文件系统抽象模型的实现结构
type fileControl struct {
	h            http.Handler
	root         string
	indexHtmlUrl string
}

// NewAndInitFileControl 文件系统控制初始化方法
func NewAndInitFileControl() {
	fc := &fileControl{nil, config.Yusha.Root, ""}
	fc.initFileSys()
	http.Handle("/", fc)
}

// 实现 http.Handler 接口
func (fc *fileControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("FileControlInter HttpRequest Message ====> host: url: " + r.URL.Path + ", method: " + r.Method)
	// 非 GET 请求直接拒绝
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if fc.checkUrl(r.URL.Path, len(r.URL.Path)) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fc.h.ServeHTTP(w, r)
}

// 初始化文件系统
func (fc *fileControl) initFileSys() {
	fi, err := os.Stat(fc.root)
	if err != nil || !fi.IsDir() {
		panic("No directory exists in the current path : " + fc.root)
	}
	fc.initIndexHtmlUrl()
	fc.h = http.FileServer(http.Dir(fc.root))
}

// url 判断与检查
func (fc *fileControl) checkUrl(url string, l int) bool {
	if url[l-1] != '/' {
		return false
	}
	if l != 1 {
		return true
	}
	_, err := os.Stat(fc.indexHtmlUrl)
	if err != nil {
		return true
	}
	return false
}

// 生成 index.html 目标文件的 url
func (fc *fileControl) initIndexHtmlUrl() {
	if strings.HasSuffix(fc.root, "/") {
		fc.indexHtmlUrl = fc.root + "index.html"
		return
	}
	fc.indexHtmlUrl = fc.root + "/index.html"
}
