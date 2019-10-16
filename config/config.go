package config

import (
	"encoding/json"
	"log"
	"io/ioutil"
	"tinyProxy/structure"
)

// 代理服务器配置信息
type Config struct {
	Service 		string `json:"service"`					// 服务名
	Host 			string `json:"host"`					// 代理服务器 host（域名或者 IP）
	Port 			uint16 `json:"port"`					// 代理服务器端口号
	WebPort 		uint16 `json:"webport"`					// web 服务器端口号
	Strategy		string `json:"strategy"`				// 代理模式（选择后端服务器的方式：轮询、随机、基于 IP 的哈希）
	Heartbeat 	 	int `json:"heartbeat"`					// 心跳机制的时间
	MaxProcessor 	int `json:"max_processor"`				// GOMAXPROCS 变量，决定会有多少个操作系统的线程同时执行Go的代码
	Backends 	 	[]structure.Backend `json:"backends"`	// 后端服务器数组，每个服务器由 host 和 port 组成
}

// 加载配置文件，并提取配置信息
func Load(filename string) (*Config, error) {
	var config Config
	file, err := ioutil.ReadFile(filename) // 读文件
	if err != nil {
		log.Println("load config file failed: ", err)
	} else {
		err = json.Unmarshal(file, &config) // 解析 json
		if err != nil {
			log.Println("xdecode json config failed: ", err)
		}
	}
	log.Println("success locad config file: ", filename)

	return &config, err
}