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

type FileControlInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	initFileSys()
	checkUrl(string, int) bool
	initIndexHtmlUrl()
}

type fileControl struct {
	h            http.Handler
	root         string
	indexHtmlUrl string
}

func NewAndInitFileControl() FileControlInter {
	fc := &fileControl{nil, config.Yusha.Root, ""}
	fc.initFileSys()
	http.Handle("/", fc)
	return fc
}

func (fc *fileControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("FileControlInter HttpRequest Message ====> url: " + r.URL.Path + ", method: " + r.Method)
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

func (fc *fileControl) initFileSys() {
	fi, err := os.Stat(fc.root)
	if err != nil || !fi.IsDir() {
		panic("No directory exists in the current path : " + fc.root)
	}
	fc.initIndexHtmlUrl()
	fc.h = http.FileServer(http.Dir(fc.root))
}

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

func (fc *fileControl) initIndexHtmlUrl() {
	if strings.HasSuffix(fc.root, "/") {
		fc.indexHtmlUrl = fc.root + "index.html"
		return
	}
	fc.indexHtmlUrl = fc.root + "/index.html"
}
