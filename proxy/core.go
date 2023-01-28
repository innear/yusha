package proxy

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"yusha/client"
	"yusha/config"
	"yusha/logger"
)

// YuShaProxyInter 代理模块的等级抽象接口
type YuShaProxyInter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	dealRequest(r *http.Request)
}

// yuShaProxy 代理模块等级抽象接口的具体实现等级父类
type yuShaProxy struct {
	addr string
	hp   string
	host string
	api  string
}

// NewAndInitProxy 代理模块对外暴露的初始化函数
func NewAndInitProxy() {
	// 判断是否需要开启代理模块
	if config.Yusha.ProxyAddr != "" && config.Yusha.ProxyPort != 0 && config.Yusha.ProxyApi != "/" {
		ysp := &yuShaProxy{}
		ysp.hp = "http"
		if config.Yusha.ProxyCertFile != "" && config.Yusha.ProxyKeyFile != "" {
			ysp.hp = "https"
		}
		ysp.host = config.Yusha.ProxyAddr + ":" + strconv.Itoa(int(config.Yusha.ProxyPort))
		ysp.api = strings.TrimSuffix(config.Yusha.ProxyApi, "/")
		http.Handle(ysp.api+"/", ysp)
		return
	}
	logger.WARN("Proxy mode is not enabled, there may be a problem with the parameters in the configuration file")
}

// 实现 golang 内部的 http.Handle 接口
func (ysp *yuShaProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 处理 Request
	ysp.dealRequest(r)

	// 代理调用
	resp, err := client.Proxy(r)
	defer resp.Body.Close()

	// 异常处理
	if err != nil && resp == nil {
		if err == client.MethodNotAllowedInProxy {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(err.Error()))
			logger.ERROR("Proxy to " + r.URL.Path + " failed : " + client.MethodNotAllowedInProxy.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		logger.ERROR("Proxy to " + r.URL.Path + " failed : " + http.StatusText(http.StatusInternalServerError))
		return
	}

	// 代理调用状态判断
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(http.StatusText(resp.StatusCode)))
		logger.ERROR("Proxy to " + r.URL.Path + " failed : " + http.StatusText(http.StatusInternalServerError))
		return
	}

	// 读取数据
	body, _ := io.ReadAll(resp.Body)

	// 同步 Response header 信息
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	// 回写数据
	w.Write(body)
	logger.INFO("Proxy to " + r.URL.Path + " success")
}

// 代理转发前对 Request 信息进行修改
func (ysp *yuShaProxy) dealRequest(r *http.Request) {
	r.URL.Host = ysp.host
	r.URL.Scheme = ysp.hp
	r.URL.Path = strings.Replace(r.URL.Path, ysp.api, "", 1)
	r.Host = ysp.host
	r.RequestURI = ""
	r.RemoteAddr = ""
}
