package schedule

import (
	"tinyProxy/util"
)

type IpHash struct {

}

func (strategy *IpHash) Init() {

}

// 选择 server 服务器
func (strategy *IpHash) Choose(client string, servers []string) string {
	// 将 IP 转化成数字，再对 servers 数量取模，得出所选 server 的索引。
	ip := util.UrlToHost(client)
	ipInt := util.IP4ToInt(ip)
	length := len(servers)
	url := servers[ipInt % length]

	return url
}
