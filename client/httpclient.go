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
	MethodNotAllowedInProxy = errors.New("the reverse proxy does not support request methods other than get and post")
)

var yuShaHttpClient = &yuShaClient{http.Client{Timeout: time.Second * time.Duration(config.Yusha.Timeout)}}

type yuShaClient struct {
	http.Client
}

func Proxy(r *http.Request) (resp *http.Response, err error) {
	switch r.Method {
	case http.MethodGet:
		return yuShaHttpClient.Get(r)
	case http.MethodPost:
		return yuShaHttpClient.Post(r)
	default:
		return nil, MethodNotAllowedInProxy
	}
}

func (ysc *yuShaClient) Get(r *http.Request) (resp *http.Response, err error) {
	return ysc.Do(r)
}

func (ysc *yuShaClient) Post(r *http.Request) (resp *http.Response, err error) {
	return ysc.Do(r)
}
