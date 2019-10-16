package structure

import (
	"net"
)

const (
	ChannelPairNum = 2	// Channel 对象的成员个数
)

type Channel struct {
	SrcConnection	net.Conn	// 与 client 的连接
	DestConnection	net.Conn	// 与 server 的连接
}

// 获取 client 的 uul
func (channel *Channel) SrcUrl() string {
	return channel.SrcConnection.RemoteAddr().String()
}

// 获取 server 的 url
func (channel *Channel) DestUrl() string {
	return channel.DestConnection.RemoteAddr().String()
}

// 关闭 client 和 server 的连接
func (channel *Channel) Close() {
	channel.SrcConnection.Close()
	channel.DestConnection.Close()
}

