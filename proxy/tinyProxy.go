package proxy

import (
	"io"
	"log"
	"net"
	"time"
	"tinyProxy/config"
	"tinyProxy/proxy/schedule"
	"tinyProxy/structure"
)

const (
	DefaultTimeout = 3	// 响应超时时间
)

// 代理对象
type TinyProxy struct {
	data		*ProxyData			// 代理数据
	strategy	schedule.Strategy	// 代理策略
}

// 初始化
func (proxy *TinyProxy) Init(config *config.Config) {
	proxy.data = new(ProxyData)
	proxy.data.Init(config)
	proxy.setStrategy(config.Strategy)
	InitStatistic(proxy.data)

}

// 设置代理策略：poll、iphash、random
func (proxy *TinyProxy) setStrategy(name string) {
	proxy.strategy = schedule.GetStrategy(name)
	proxy.strategy.Init()
}

// 检查 Backends 和 Deads 中 server。
// 如果 Backends 中的 server 不可用，则将其移动到 Deads 中。
// 如果 Deads 中的 server 可用，则将其移回到 Backends 中。
func (proxy *TinyProxy) Check() {
	for _, backend := range proxy.data.Backends {
		_, err := net.Dial("tcp", backend.Url())
		if err != nil {
			proxy.Delete(backend.Url())
		}
	}

	for _, deadend := range proxy.data.Deads {
		_, err := net.Dial("tcp", deadend.Url())
		if err != nil {
			proxy.Recover(deadend.Url())
		}
	}
}

// 检查 Backends 是否有可用 server
func (proxy *TinyProxy) isBackendAvaiable() bool {
	return len(proxy.data.Backends) > 0
}

// 派发连接
func (proxy *TinyProxy) Dispatch(connection net.Conn) {
	if proxy.isBackendAvaiable() {
		servers := proxy.data.BackendUrls()
		// 服务器 url 地址
		url := proxy.strategy.Choose(connection.RemoteAddr().String(), servers)
		proxy.transfer(connection, url)
	} else {
		connection.Close()
		log.Println("no backends available now,please check your server!")
	}
}

// 转发数据
func (proxy *TinyProxy) safeCopy(from net.Conn, to net.Conn, syncChan chan int) {
	// Copy 会一直运行，直到其中一个 net.Conn 关闭为止
	_, _ = io.Copy(to, from)
	defer from.Close()

	// 连接关闭，向通道发送数据
	syncChan <- 1
}

// 添加 Channel
func (proxy *TinyProxy) putChannel(channel *structure.Channel) {
	proxy.data.ChannelManager.PutChannel(channel)
}

// 关闭 Channel
func (proxy *TinyProxy) closeChannel(channel *structure.Channel, syncChan chan int) {
	// 当从 syncChan 读取了 ChannelPairNum 个，即 2 个值，则表示该连接已关闭，删除所在 Channel 即可。
	for i := 0; i < structure.ChannelPairNum; i++ {
		<-syncChan
	}
	proxy.data.ChannelManager.DeleteChannel(channel)
}

// 为 client 和 server 建立连接
func (proxy *TinyProxy) transfer(clientConnection net.Conn, remote string) {
	remoteConnection, err := net.DialTimeout("tcp", remote, DefaultTimeout * time.Second)
	// server 不可用
	if err != nil {
		clientConnection.Close()
		proxy.Delete(remote)
		log.Println("connection backend error: %s", err)

		return
	}

	syncChan := make(chan int, 1)
	channel := structure.Channel{
		SrcConnection:  clientConnection,
		DestConnection: remoteConnection,
	}

	proxy.putChannel(&channel)
	go proxy.safeCopy(clientConnection, remoteConnection, syncChan)
	go proxy.safeCopy(remoteConnection, clientConnection, syncChan)
	go proxy.closeChannel(&channel, syncChan)
}

// 移除 server：将 Backends 中的 server 移动到 Deads 中。
func (proxy *TinyProxy) Delete(url string) {
	proxy.data.DeleteBackend(url)
}

// 恢复 server：将 Deads 中的 server 移回到 Backends 中。
func (proxy *TinyProxy) Recover(url string) {
	proxy.data.DeleteDeadend(url)
}

// 关闭服务器，清空所有代理数据
func (proxy *TinyProxy) Close() {
	proxy.data.Clear()
}
