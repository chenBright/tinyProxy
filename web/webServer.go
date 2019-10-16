package web

import (
	"net/http"
	"log"
	_ "net/http/pprof"
	"tinyProxy/config"
	"tinyProxy/util"
)

// 数据服务器
type WebServer struct {
	host	string	// host
	port	uint16	// 端口号
}

// 初始化
func (server *WebServer) Init(config *config.Config) {
	server.host = config.Host
	server.port = config.WebPort
}

// 启动服务器
func (server *WebServer) Start() {
	// 在一个协程内启动服务器
	go func() {
		server.AddHandler() // 设置 handler，处理请求
		url := util.HostPortToAddress(server.host, server.port) // 获取 url 字符串
		err := http.ListenAndServe(url, nil) // 启动 http 服务器
		if err != nil {
			log.Println("create web server failed: ", err)
		} else {
			log.Println("create web server success")
		}
	}()
}

// 设置 handler，处理请求
func (server *WebServer) AddHandler() {
	http.HandleFunc(StatisticsUrl, Statistic)
}