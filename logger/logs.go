// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"log"
)

var logChan = make(chan *yuShaLog, 200)

/**
后续日志模块的功能在这包下实现
具体实现还要规划一下
*/

// 日志结构体模型
type yuShaLog struct {
	t string
	v string
}

// 日志服务总线
func logServer() {
	for {
		l := <-logChan
		log.Println(l.t + l.v)
	}
}

func init() {
	go logServer()
}

func INFO(v string) {
	l := &yuShaLog{INFO_, v}
	logChan <- l
}

func WARN(v string) {
	l := &yuShaLog{WARN_, v}
	logChan <- l
}

func ERROR(v string) {
	l := &yuShaLog{ERROR_, v}
	logChan <- l
}
