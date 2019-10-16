package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"tinyProxy/config"
	"tinyProxy/gateway"
	"tinyProxy/log"
	"tinyProxy/util"
	"tinyProxy/web"
)

const (
	DefaultConfigFile = "config_file/default.json"
	DefaultLogFile = "tiny_proxy.log"
)

// 服务器
type TinyServer struct {
	webServer 	*web.WebServer			// 数据服务器
	proxyServer	*gateway.ProxyServer	// 代理服务器
}

// 创建服务器
func CreateTinyServer() *TinyServer {
	return &TinyServer{
		webServer:   new(web.WebServer),
		proxyServer: new(gateway.ProxyServer),
	}
}

// 初始化服务器
func (tinyServer *TinyServer) Init(config *config.Config) {
	tinyServer.webServer.Init(config)
	tinyServer.proxyServer.Init(config)
}

// 启动服务器
func (tinyServer *TinyServer) Start() {
	tinyServer.webServer.Start()
	tinyServer.proxyServer.Start()
}

// 关闭服务器
func (tinyServer *TinyServer) Stop() {
	tinyServer.proxyServer.Stop()
}

// 设置捕获 stop 信号。
// 在捕获 stop 信号后，关闭服务器。
func (tinyServer *TinyServer) CatchStopSignal() {
	sig := make(chan os.Signal, 1)
	// 当捕获到 SIGINT、SIGQUIT 或者 SIGKILL 信号时，发送信号到 sig 通道。
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	// 启动一个协程。
	// 一开始，协程阻塞在 sig 通道处。
	// 当有 sig 可读时，表示捕获到 stop 信号了，则立即关闭服务器。
	go func() {
		<-sig
		tinyServer.Stop()
	}()
}

func main() {
	log.Init(DefaultLogFile) // 初始日志系统

	homePath := util.HomePath() // 获取当前项目的绝对路径
	// 加载配置，并得到配置对象
	config_, err := config.Load(filepath.Join(homePath, DefaultConfigFile))

	if err == nil {
		// 设置 GOMAXPROCS
		runtime.GOMAXPROCS(config_.MaxProcessor)

		tinyServer := CreateTinyServer()
		tinyServer.Init(config_)
		tinyServer.CatchStopSignal()
		tinyServer.Start()
	}
}