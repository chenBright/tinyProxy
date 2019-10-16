package proxy

import (
	"sync"
	"tinyProxy/config"
	"tinyProxy/structure"
)

// 代理数据
type ProxyData struct {
	Service 		string							// 服务名
	Host 			string							// 代理服务器 host
	Port			uint16							// 代理服务器 端口号
	Backends 		map[string]structure.Backend	// 在线 server: url -> server 对象
	Deads 			map[string]structure.Backend	// 不在线 server: url -> server 对象
	ChannelManager	*structure.ChannelManager		// Channel 管理对象
	mutex 			*sync.RWMutex					// 互斥锁
}

// 初始化
func (proxyData *ProxyData) Init(config *config.Config) {
	proxyData.Service = config.Service
	proxyData.Host = config.Host
	proxyData.Port = config.Port
	proxyData.setBackends(config.Backends)
	proxyData.ChannelManager = new(structure.ChannelManager)
	proxyData.ChannelManager.Init()
	proxyData.mutex = new(sync.RWMutex)
}

//
func (proxyData *ProxyData) setBackends(backends []structure.Backend) {
	proxyData.Backends = make(map[string]structure.Backend)
	for _, backend := range backends {
		proxyData.Backends[backend.Url()] = backend
	}

	proxyData.Deads = make(map[string]structure.Backend)
}

// 获取所有在线的 server url
func (proxyData *ProxyData) BackendUrls() []string {
	proxyData.mutex.RLock()
	defer proxyData.mutex.RUnlock()

	backends := proxyData.Backends
	keys := make([]string, 0, len(backends))
	for key := range backends {
		keys = append(keys, key)
	}

	return keys
}

// 将 Backends 的某一 server 移动到 Deads 中
func (proxyData *ProxyData) DeleteBackend(url string) {
	proxyData.mutex.Lock()
	defer proxyData.mutex.Unlock()

	proxyData.Deads[url] = proxyData.Backends[url]
	delete(proxyData.Backends, url)
}

// 从 Deads 中删除 server
func (proxyData *ProxyData) DeleteDeadend(url string) {
	proxyData.mutex.Lock()
	proxyData.Backends[url] = proxyData.Deads[url]
	delete(proxyData.Deads, url)

	defer proxyData.mutex.Unlock()
}

// 清空 backends
func clearBackendMap(backends map[string]structure.Backend) {
	for backend := range backends {
		delete(backends, backend)
	}
}

// 清空所有代理数据
func (proxyData *ProxyData) Clear()  {
	clearBackendMap(proxyData.Backends)
	clearBackendMap(proxyData.Deads)
	proxyData.ChannelManager.Clear()
}
