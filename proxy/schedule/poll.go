package schedule

import (
	"sync"
)

const CycleCount = 1 << 21	// 用于取模

// 轮询对象
type Poll struct {
	counter Counter
}

// 计数器
type Counter struct {
	count int
	mutex *sync.Mutex
}

func (counter *Counter) Increase() {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	counter.count = (counter.count + 1) % CycleCount
}

func (counter *Counter) Get() int {
	return counter.count
}

func (strategy *Poll) Init()  {
	strategy.counter = Counter{
		count: 0,
		mutex: new(sync.Mutex)}
}

func (strategy *Poll) Choose(client string, servers []string) string {
	strategy.counter.Increase()
	length := len(servers)
	url := servers[strategy.counter.Get() % length]

	return url
}
