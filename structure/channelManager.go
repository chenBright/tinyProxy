package structure

import (
	"errors"
	"sync"
	"tinyProxy/util"
)

//  Channel 管理对象
type ChannelManager struct {
	channels	[]Channel			// channel 列表
	SrcMap		map[string]*Channel	// url -> Channel
	DestMap		map[string]*Channel	// url -> Channel
	mutex		*sync.Mutex			// 互斥量
}

// 初始化
func (channeManager *ChannelManager) Init()  {
	channeManager.channels = make([]Channel, 0)
	channeManager.SrcMap = make(map[string]*Channel)
	channeManager.DestMap = make(map[string]*Channel)
	channeManager.mutex = new(sync.Mutex)
}

// 添加 Channel
func (channeManager *ChannelManager) PutChannel(channel *Channel) {
	channeManager.mutex.Lock()
	defer channeManager.mutex.Unlock()

	channeManager.channels = append(channeManager.channels, *channel)
	channeManager.SrcMap[channel.SrcUrl()] = channel
	channeManager.DestMap[channel.DestUrl()] = channel
}

// 删除 Channel
func (channeManager *ChannelManager) DeleteChannel(channel *Channel) {
	channeManager.mutex.Lock()
	defer channeManager.mutex.Unlock()

	index := util.SliceIndex(channeManager.channels, *channel) // 获取 Channel 的索引
	if index >= 0 {
		// 删除 channels 中的 Channel
		channeManager.channels = append(
			channeManager.channels[:index], channeManager.channels[index + 1:]...)
		// 删除 Map 中的 Channel 映射
		deleteChannelInMap(channeManager.SrcMap, channel.SrcUrl())
		deleteChannelInMap(channeManager.DestMap, channel.DestUrl())
	}
}

// 获取 Channel列表
func (channeManager *ChannelManager) GetChannel() []Channel {
	return channeManager.channels
}

// 检查列表和 Map 中 Channel 的数量是否一致
func (channeManager *ChannelManager) Check() (error, error) {
	var srcErr, destErr error
	channelsLength := len(channeManager.channels)
	if len(channeManager.SrcMap) != channelsLength {
		srcErr = errors.New("client socket close maybe error")
	}
	
	if len(channeManager.DestMap) != channelsLength {
		srcErr = errors.New("server socket close maybe error")
	}
	
	return srcErr, destErr
}

// 删除 Map 中的 Channel
func deleteChannelInMap(map_ map[string]*Channel, url string) {
	_, ok := map_[url]
	if ok {
		delete(map_, url)
	}
}

// 清空所有 Channel
func (channeManager *ChannelManager) Clear() {
	for _, channel := range channeManager.channels {
		deleteChannelInMap(channeManager.SrcMap, channel.SrcUrl())
		deleteChannelInMap(channeManager.DestMap, channel.DestUrl())
		channel.Close()
	}

	channeManager.channels = channeManager.channels[:0]
}