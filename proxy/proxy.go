package proxy

import (
	"net"
	"tinyProxy/config"
)

// 代理对象接口
type Proxy interface {
	Init(ProxyConfig *config.Config)
	Check()
	Delete(url string)
	Recover(url string)
	Dispatch(connection net.Conn)
	Close()
}