package schedule

import (
	"time"
)

type Random struct {

}

func (stratery *Random) Init() {

}

func (strategy *Random) Choose(client string, servers []string) string {
	length := len(servers)
	// 使用时间戳对 server 数量取模，作为所选 server 的索引。
	url := servers[int(time.Now().UnixNano()) % length]

	return url
}
