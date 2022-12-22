// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ys

import (
	"net/http"
	"strconv"
	"yusha/config"
	"yusha/filesys"
	"yusha/proxy"
)

var port string

func init() {
	filesys.NewAndInitFileControl()
	proxy.NewAndInitProxy()
	port = ":" + strconv.Itoa(int(config.Yusha.Port))
}

func Run() {
	if config.Yusha.CertFile != "" && config.Yusha.KeyFile != "" {
		err := http.ListenAndServeTLS(port, config.Yusha.CertFile, config.Yusha.KeyFile, nil)
		if err != nil {
			panic(err)
		}
	}
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
