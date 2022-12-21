// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import (
	"github.com/vrbyte/yusha/config"
	"github.com/vrbyte/yusha/filesys"
	"log"
	"net/http"
	"os"
)

type YuShaAgentInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type YuShaAgent struct {
	Fc filesys.FileControlInter
}

func (ysa *YuShaAgent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		goto agent
	}
	if _, err := os.Stat(config.Yusha.Root + r.URL.Path); err == nil {
		ysa.Fc.ServeHTTP(w, r)
		return
	}
agent:
	log.Println("启用代理")
}
