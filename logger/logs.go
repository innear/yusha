// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"log"
	"os"
	"time"
)

// 日志总线参数声明
var (
	logChan   = make(chan *yuShaLog, 300)
	infoFile  *os.File
	warnFile  *os.File
	errorFile *os.File
	infoLog   *log.Logger
	warnLog   *log.Logger
	errorLog  *log.Logger
)

// 日志结构体模型
type yuShaLog struct {
	t int
	v string
}

// 日志服务总线
func logServer() {
	for {
		l := <-logChan
		switch l.t {
		case INFO_:
			infoLog.Println(l.v)
		case WARN_:
			warnLog.Println(l.v)
		case ERROR_:
			errorLog.Println(l.v)
		default:
			infoLog.Println(l.v)
		}
	}
}

func init() {
	go logServer()
	_, err := os.Stat("./log")
	if err != nil {
		os.Mkdir("log", 0777)
	}
	// 初始化日志服务参数
	infoFile, _ = os.OpenFile("./log/yusha-info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	warnFile, _ = os.OpenFile("./log/yusha-warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	errorFile, _ = os.OpenFile("./log/yusha-error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	infoLog = log.New(infoFile, "[INFO] ", 3)
	warnLog = log.New(warnFile, "[WARN] ", 3)
	errorLog = log.New(errorFile, "[ERROR] ", 3)
}

// 对外暴露的日志使用方法

func INFO(val string) {
	logChan <- &yuShaLog{t: INFO_, v: val}
}

func WARN(val string) {
	logChan <- &yuShaLog{t: WARN_, v: val}
}

func ERROR(val string) {
	logChan <- &yuShaLog{t: ERROR_, v: val}
}

// CheckLogChan 检查日志管道内消息是否全部消费完毕
func CheckLogChan() {
	// 睡眠 1 秒, 确保所有日志都进入管道中
	time.Sleep(time.Second)
	for {
		// 日志管道中所有消息都被消费完毕后结束死循环
		if len(logChan) == 0 {
			break
		}
	}
}
