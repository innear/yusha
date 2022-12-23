// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"yusha/client"
	"yusha/config"
	"yusha/logger"
)

// YuShaProxyInter 代理模块的等级抽象接口
type YuShaProxyInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	DealRequest(r *http.Request)
}

// YuShaProxy 代理模块等级抽象接口的具体实现等级父类
type YuShaProxy struct {
	addr string
	hp   string
	host string
	api  string
}

// NewAndInitProxy 代理模块对外暴露的初始化函数
func NewAndInitProxy() {
	// 判断是否需要开启代理模块
	if config.Yusha.ProxyAddr != "" && config.Yusha.ProxyPort != 0 && config.Yusha.ProxyApi != "/" {
		ysp := &YuShaProxy{}
		ysp.hp = "https"
		if config.Yusha.CertFile != "" && config.Yusha.KeyFile != "" {
			ysp.hp = "https"
		}
		ysp.host = config.Yusha.ProxyAddr + ":" + strconv.Itoa(int(config.Yusha.ProxyPort))
		ysp.api = strings.TrimSuffix(config.Yusha.ProxyApi, "/")
		http.Handle(ysp.api+"/", ysp)
		return
	}
	logger.LogOut(logger.WARN, "Proxy mode is not enabled, there may be a problem with the parameters in the configuration file")
}

// 实现 golang 内部的 http.Handle 接口
func (ysp *YuShaProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ysp.DealRequest(r)

	resp, err := client.Proxy(r)
	defer resp.Body.Close()

	if err != nil {
		if err == client.MethodNotAllowed {
			w.WriteHeader(http.StatusMethodNotAllowed)
			goto end
		}
		if resp != nil {
			w.WriteHeader(resp.StatusCode)
			goto end
		}
	end:
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(http.StatusText(resp.StatusCode)))
		return
	}

	body, _ := io.ReadAll(resp.Body)

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.Write(body)
}

func (ysp *YuShaProxy) DealRequest(r *http.Request) {
	r.URL.Host = ysp.host
	r.URL.Scheme = ysp.hp
	r.URL.Path = strings.Replace(r.URL.Path, ysp.api, "", 1)
	r.Host = ysp.host
	r.RequestURI = ""
	r.RemoteAddr = ""
}
