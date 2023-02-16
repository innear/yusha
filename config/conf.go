package config

import (
	"encoding/json"
	"os"
	"strings"
	"yusha/logger"
)

// 默认配置文件路径参数(默认配置文件路径为 ./conf/yusha.json)
var defaultProfilePath = "./conf/yusha.json"

// YuShaConf 配置参数结构体
/**
Root 静态资源代理根路径(默认根路径 ./html)
Port 监听端口(默认端口 8100)
CertFile TLS 加密需要的证书文件路径
KeyFile TLS 加密需要的密钥文件路径
ProxyAddr 代理地址(可以为ip或者域名)
ProxyApi 代理接口 api 前缀标识
ProxyCertFile 代理接口加密需要的证书文件路径
ProxyKeyFile  代理接口加密需要的密钥文件路径
TimeOut http 请求代理转发超时时间参数(单位秒)
*/
type YuShaConf struct {
	Root          string
	Port          uint16
	CertFile      string
	KeyFile       string
	ProxyAddr     string
	ProxyPort     uint16
	ProxyApi      string
	ProxyCertFile string
	ProxyKeyFile  string
	Timeout       int
}

// Yusha 全局配置参数
var Yusha *YuShaConf

// 初始化
func init() {
	defer logger.CheckLogChan()
	Yusha = &YuShaConf{
		Root:     "./html",
		Port:     8100,
		ProxyApi: "/api/",
		Timeout:  3,
	}
	_, err := os.Stat(defaultProfilePath)
	if err != nil {
		logger.WARN("No corresponding file found in the default configuration file path : " + defaultProfilePath)
		logger.WARN("Default configuration will be enabled in Yusha")
		return
	}
	// 读取 JSON 配置文件
	b, _ := os.ReadFile(defaultProfilePath)
	err = json.Unmarshal(b, Yusha)
	if err != nil {
		logger.ERROR("Failed to transfer the configuration file content to JSON")
		panic(err)
	}

	if Yusha.CertFile != "" {
		_, err := os.Stat(Yusha.CertFile)
		if err != nil {
			panic(err)
		}
	}

	if Yusha.KeyFile != "" {
		_, err := os.Stat(Yusha.KeyFile)
		if err != nil {
			panic(err)
		}
	}

	if Yusha.ProxyCertFile != "" {
		_, err := os.Stat(Yusha.ProxyCertFile)
		if err != nil {
			panic(err)
		}
	}

	if Yusha.ProxyKeyFile != "" {
		_, err := os.Stat(Yusha.ProxyKeyFile)
		if err != nil {
			panic(err)
		}
	}

	if !strings.HasPrefix(Yusha.ProxyApi, "/") {
		Yusha.ProxyApi = "/" + Yusha.ProxyApi
	}

	if !strings.HasSuffix(Yusha.ProxyApi, "/") {
		Yusha.ProxyApi += "/"
	}
}
