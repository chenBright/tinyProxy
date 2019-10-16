package proxy

import "tinyProxy/util"

var statistic = new(Statistic)	// 数据对象，单例

// 数据对象
type Statistic struct {
	Clients   map[string]Client // url -> client
	Services  map[string]Server // url -> service
	proxyData *ProxyData        // 代理数据
}

// client 对象
type Client struct {
	Host	string
	Count	int
}

// server 对象
type Server struct {
	Url 	string
	Count 	int
	Status 	string
}

// 初始化
func InitStatistic(proxyData *ProxyData) {
	statistic.proxyData = proxyData
}

// 获取数据对象
func StatisticData() *Statistic {
	return statistic
}

// 记录数据
func Record() {
	statistic.Clients = make(map[string]Client)
	statistic.Services = make(map[string]Server)
	// 记录在线的 server
	for _, server := range statistic.proxyData.Backends {
		statistic.Services[server.Url()] = Server{
			Url:    server.Url(),
			Count:  0,
			Status: "on",
		}
	}

	// 记录不在线的 server
	for _, server := range statistic.proxyData.Deads {
		statistic.Services[server.Url()] = Server{
			Url:    server.Url(),
			Count:  0,
			Status: "off",
		}
	}

	// 记录 client
	for _, channel := range statistic.proxyData.ChannelManager.GetChannel() {
		host := util.UrlToHost(channel.SrcUrl())
		serverUrl := channel.DestUrl()
		var client Client
		var service Server

		if _, ok := statistic.Clients[host]; ok {
			client = statistic.Clients[host]
			client.Count += 1
		} else {
			client = Client{
				Host:  host,
				Count: 1,
			}
		}

		service = statistic.Services[serverUrl]
		service.Count += 1

		statistic.Clients[host] = client
		statistic.Services[serverUrl] = service
	}
}