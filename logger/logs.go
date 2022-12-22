// Copyright 2022 The Yusha Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"log"
)

/**
后续日志模块的功能在这包下实现
具体实现还要规划一下
*/

func LogOut(t int, v any) {
	switch t {
	case INFO:
		log.Println("INFO", v)
	case WARN:
		log.Println("WARN", v)
	case ERROR:
		log.Println("ERROR", v)
	}
}
