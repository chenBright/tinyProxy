package gateway

import (
	"net"
	"log"
	"time"
	"tinyProxy/config"
	"tinyProxy/proxy"
	"tinyProxy/util"
)

// 代理服务器
type ProxyServer struct {
	host		string			// 服务器 host
	port		uint16			// 服务器端口
	beattime	int				// 心跳周期
	listener	net.Listener	// 监听对象
	proxy		proxy.Proxy		// 代理对象
	on			bool			// 是否可用
}

// 初始化
func (server *ProxyServer) Init(config *config.Config) {
	server.on = false
	server.host = config.Host
	server.port = config.Port
	server.beattime = config.Heartbeat
	server.setProxy(config)
}

// 设置代理信息
func (server *ProxyServer) setProxy(config *config.Config) {
	server.proxy = new(proxy.TinyProxy)
	server.proxy.Init(config)
}

// 获取代理服务器 url
func (server *ProxyServer) Address() string {
	return util.HostPortToAddress(server.host, server.port)
}

// 启动代理服务器
func (server *ProxyServer) Start() {
	local, err := net.Listen("tcp", server.Address())

	if err != nil {
		log.Panic("proxy server start error:", err)
	} else {
		log.Println("tiny proxy server start ok")
		server.listener = local
		server.on = true
		server.heartBeat()
		defer server.listener.Close()
	}

	// 等待连接到来
	for server.on {
		connection, err := server.listener.Accept()
		if err == nil {
			go server.proxy.Dispatch(connection)
		} else {
			log.Println("client connect server error: ", err)
		}
	}
}

// 启动心跳机制
func (server *ProxyServer) heartBeat() {
	// 周期地往 ticker 通道发送数据
	ticker := time.NewTicker(time.Second * time.Duration(server.beattime))
	go func() {
		for {
			select {
			// 周期地从 ticker 读取数据，然后检查所有 server 是否可用
			case <-ticker.C:
				server.proxy.Check()
			}
		}
	}()
}

// 关闭服务器
func (server *ProxyServer) Stop() {
	_ = server.listener.Close()
	server.proxy.Close()
	server.on = false

	log.Println("tiny proxy server stop")

}
