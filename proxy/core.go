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

type YuShaProxyInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	ProxyTo()
}

type YuShaProxy struct {
	addr string
	api  string
}

func NewAndInitProxy() {
	if config.Yusha.ProxyAddr != "" && config.Yusha.ProxyPort != 0 && config.Yusha.ProxyApi != "" && config.Yusha.ProxyApi != "/" {
		ysp := &YuShaProxy{}
		ysp.api = config.Yusha.ProxyApi
		http.Handle(ysp.api, ysp)
		return
	}
	logger.LogOut(logger.WARN, "Proxy mode is not enabled, there may be a problem with the parameters in the configuration file")
}

func (ysp *YuShaProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("into proxy")
	log.Println(r.URL.Path)
	resp, err := http.Post("https://"+config.Yusha.ProxyAddr+":"+strconv.Itoa(int(config.Yusha.ProxyPort))+"/"+strings.Replace(r.URL.Path, ysp.api, "", 1), "application/json;charset=UTF-8", nil)

	if err != nil || resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(body)
}

func (ysp *YuShaProxy) ProxyTo() {

}
