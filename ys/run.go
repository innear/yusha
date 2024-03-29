package ys

import (
	"net/http"
	"strconv"
	"yusha/config"
	"yusha/filesys"
	"yusha/logger"
	"yusha/proxy"
)

var port string

// golang 加载机制触发 init 函数
func init() {
	// 初始化文件管理系统
	filesys.NewAndInitFileControl()
	// 反向代理机制初始化
	proxy.NewAndInitProxy()
	port = ":" + strconv.Itoa(int(config.Yusha.Port))
}

// Run 运行函数
func Run() {
	defer logger.CheckLogChan()
	// 判断是否需要 TLS
	if config.Yusha.CertFile != "" && config.Yusha.KeyFile != "" {
		err := http.ListenAndServeTLS(port, config.Yusha.CertFile, config.Yusha.KeyFile, nil)
		if err != nil {
			logger.ERROR(err.Error())
			panic(err)
		}
	}
	// 监听
	err := http.ListenAndServe(port, nil)
	if err != nil {
		logger.ERROR(err.Error())
		panic(err)
	}
}
