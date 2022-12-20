// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"log"
	"os"
)

var defaultProfilePath = "./conf/yusha.json"

type YuShaConf struct {
	Root string
}

var Yusha *YuShaConf

// 初始化
func init() {
	Yusha = &YuShaConf{
		Root: "./html",
	}
	_, err := os.Stat(defaultProfilePath)
	if err != nil {
		log.Println("No corresponding file found in the default configuration file path : " + defaultProfilePath)
		log.Println("Default configuration will be enabled in Yusha")
		return
	}
	b, _ := os.ReadFile(defaultProfilePath)
	err = json.Unmarshal(b, Yusha)
	if err != nil {
		log.Println("Failed to transfer the configuration file content to JSON")
		panic(err)
	}
}
