// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"yusha/config"
	"yusha/logger"
)

// YuShaProxyInter 代理模块的等级抽象接口
type YuShaProxyInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	ProxyGet()
	ProxyPost()
}

// YuShaProxy 代理模块等级抽象接口的具体实现
type YuShaProxy struct {
	addr string
	api  string
}

// NewAndInitProxy 代理模块对外暴露的初始化函数
func NewAndInitProxy() {
	// 判断是否需要开启代理模块
	if config.Yusha.ProxyAddr != "" && config.Yusha.ProxyPort != 0 && config.Yusha.ProxyApi != "" && config.Yusha.ProxyApi != "/" {
		ysp := &YuShaProxy{}
		ysp.api = config.Yusha.ProxyApi
		http.Handle(ysp.api, ysp)
		return
	}
	logger.LogOut(logger.WARN, "Proxy mode is not enabled, there may be a problem with the parameters in the configuration file")
}

// 实现 golang 内部的 http.Handle 接口
func (ysp *YuShaProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("into proxy")
	log.Println(r.URL.Path)

	/**
	这里仅仅是简单模拟了一下代理调用
	后续要处理很多细节问题, 例如用户自定义 header 的传递; 远程调用返回结果携带 header 的处理等
	*/

	resp, err := http.Get("https://" + config.Yusha.ProxyAddr + ":" + strconv.Itoa(int(config.Yusha.ProxyPort)) + "/" + strings.Replace(r.URL.Path, ysp.api, "", 1))

	if err != nil || resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		// w.Write([]byte(http.StatusText(resp.StatusCode)))
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(body)
}

func (ysp *YuShaProxy) ProxyTo() {

}
