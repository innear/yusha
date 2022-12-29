// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesys

import (
	"net/http"
	"os"
	"strings"
	"yusha/config"
	"yusha/logger"
)

// FileControlInter 静态资源代理(文件系统)顶级抽象接口
type FileControlInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	initFileSys()
	checkUrl(string, int) bool
	initIndexHtmlUrl()
}

// 静态资源代理(文件系统)顶级抽象接口的实现
type fileControl struct {
	h            http.Handler
	root         string
	indexHtmlUrl string
}

// NewAndInitFileControl 静态资源代理(文件系统)模块对外暴露的初始化接口
func NewAndInitFileControl() {
	fc := &fileControl{nil, config.Yusha.Root, ""}
	fc.initFileSys()
	http.Handle("/", fc)
}

// 实现 golang 内部的 http.Handle 接口
func (fc *fileControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.INFO("FileControlInter HttpRequest Message ====> url: " + r.URL.Path)
	// 静态资源访问模块非 GET 请求一律驳回
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

// 主要初始化函数
func (fc *fileControl) initFileSys() {
	defer logger.CheckLogChan()
	// 检测静态资源代理文件的顶级路径(默认为相对路径的 ./html 下)
	fi, err := os.Stat(fc.root)
	if err != nil || !fi.IsDir() {
		logger.ERROR("No directory exists in the current path : " + fc.root)
		panic("No directory exists in the current path : " + fc.root)
	}
	// 初始化 index.html 的路径
	fc.initIndexHtmlUrl()
	// 加载 golang 内部文件系统机制接口
	fc.h = http.FileServer(http.Dir(fc.root))
}

// url 判断机制
func (fc *fileControl) checkUrl(url string, l int) bool {
	// 规避重定向过来的文件目录访问, 防止资源泄露
	if url[l-1] != '/' {
		return false
	}
	// 放行 /* 路径
	if l != 1 {
		return true
	}
	// 能达到这一级判断的只有 / 根路径访问
	// 如果静态资源代理的顶级路径一级目录下有 index.html 文件, 则放行 golang 内部会将此访问关联至 /index.html
	// 如果不存在则拒绝, 规避根路径的文件目录访问, 防止资源泄露
	_, err := os.Stat(fc.indexHtmlUrl)
	if err != nil {
		return true
	}
	return false
}

// 初始化 index.html 文件路径
func (fc *fileControl) initIndexHtmlUrl() {
	if strings.HasSuffix(fc.root, "/") {
		fc.indexHtmlUrl = fc.root + "index.html"
		return
	}
	fc.indexHtmlUrl = fc.root + "/index.html"
}
