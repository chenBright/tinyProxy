package structure

import (
	"tinyProxy/util"
)

// server 对象
type Backend struct {
	Host string `json:"host"`	// host
	Port uint16 `json:"port"`	// 端口号
}

// 获取 server url
func (backend Backend) Url() string {
	return util.HostPortToAddress(backend.Host, backend.Port)
}