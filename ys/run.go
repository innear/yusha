// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ys

import (
	"github.com/vrbyte/yusha/agent"
	"github.com/vrbyte/yusha/config"
	"github.com/vrbyte/yusha/filesys"
	"net/http"
	"strconv"
)

var port string

var fc filesys.FileControlInter

func init() {
	fc = filesys.NewAndInitFileControl()
	port = ":" + strconv.Itoa(int(config.Yusha.Port))
}

func Run() {
	if config.Yusha.Tls {
		err := http.ListenAndServeTLS(port, config.Yusha.CertFile, config.Yusha.KeyFile, nil)
		if err != nil {
			panic(err)
		}
	}
	err := http.ListenAndServe(port, &agent.YuShaAgent{Fc: fc})
	if err != nil {
		panic(err)
	}
}
