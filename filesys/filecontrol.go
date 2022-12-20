// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesys

import (
	"log"
	"net/http"
	"os"
)

type FileControlInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	initFileSys()
	checkUrl(string) bool
	searchIndexHtml()
}

type fileControl struct {
	h             http.Handler
	haveIndexHtml bool
	path          string
}

func NewFileControl(path string) http.Handler {
	fc := &fileControl{nil, false, path}
	fc.initFileSys()
	return fc
}

func (fc *fileControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("HttpRequest Message ====> url: " + r.URL.Path + ", method:" + r.Method)
	if r.Method != http.MethodGet {
		return
	}
	fc.h.ServeHTTP(w, r)
}

func (fc *fileControl) initFileSys() {
	fi, err := os.Stat(fc.path)
	if err != nil || !fi.IsDir() {
		panic("No directory named html exists in the current path")
	}
	fc.searchIndexHtml()
	fc.h = http.FileServer(http.Dir(fc.path))
}

func (fc *fileControl) checkUrl(url string) (sign bool) {
	return
}

func (fc *fileControl) searchIndexHtml() {
	f, _ := os.Open(fc.path)
	defer f.Close()
	fis, _ := f.Readdir(-1)
	for _, fi := range fis {
		if fi.Name() == "index.html" {
			fc.haveIndexHtml = true
			break
		}
	}
}
