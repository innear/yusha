// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"yusha/logger"
)

var defaultProfilePath = "./conf/yusha.json"

type YuShaConf struct {
	Root      string
	Port      uint16
	CertFile  string
	KeyFile   string
	ProxyAddr string
	ProxyPort uint16
	ProxyApi  string
}

var Yusha *YuShaConf

// 初始化
func init() {
	Yusha = &YuShaConf{
		Root: "./html",
		Port: 8100,
	}
	_, err := os.Stat(defaultProfilePath)
	if err != nil {
		logger.LogOut(logger.WARN, "No corresponding file found in the default configuration file path : "+defaultProfilePath)
		logger.LogOut(logger.WARN, "Default configuration will be enabled in Yusha")
		return
	}
	b, _ := os.ReadFile(defaultProfilePath)
	err = json.Unmarshal(b, Yusha)
	if err != nil {
		log.Println("Failed to transfer the configuration file content to JSON")
		panic(err)
	}
	if Yusha.ProxyApi != "" && !strings.HasPrefix(Yusha.ProxyApi, "/") {
		Yusha.ProxyApi = "/" + Yusha.ProxyApi
	}
}
