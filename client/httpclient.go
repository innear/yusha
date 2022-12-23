// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"errors"
	"net/http"
	"time"
	"yusha/config"
)

var (
	MethodNotAllowed = errors.New("method not support")
)

var YuShaHttpClient = &YuShaClient{http.Client{Timeout: time.Second * time.Duration(config.Yusha.Timeout)}}

type YuShaClient struct {
	http.Client
}

func Proxy(r *http.Request) (resp *http.Response, err error) {
	switch r.Method {
	case http.MethodGet:
		return YuShaHttpClient.Get(r)
	case http.MethodPost:
		return YuShaHttpClient.Post(r)
	default:
		return nil, MethodNotAllowed
	}
}

func (ysc *YuShaClient) Get(r *http.Request) (resp *http.Response, err error) {
	return ysc.Do(r)
}

func (ysc *YuShaClient) Post(r *http.Request) (resp *http.Response, err error) {
	return ysc.Do(r)
}
