package schedule

const (
	PollName = "poll"
	IpHashName = "iphash"
	RandomName = "random"
)

var registry = make(map[string]Strategy)	// 代理策略名 -> 代理策略对象

type Strategy interface {
	// 初始化
	Init()
	// 选择 server 服务器，返回 server 服务器的 url
	// client string: client 的 url
	// servers []string: server 列表
	Choose(client string, servers []string) string
}

// 初始化
func init() {
	registry[PollName] = new(Poll)
	registry[IpHashName] = new(IpHash)
	registry[RandomName] = new(Random)
}

// 获取策略
func GetStrategy(name string) Strategy {
	return registry[name]
}