package web

import (
	"net/http"
	"tinyProxy/proxy"
)

const (
	StatisticsUrl = "/statistic"
)

// 从代理服务器获取数据，并渲染页面
func Statistic(writer http.ResponseWriter, request *http.Request) {
	proxy.Record()
	Render(writer, "statistic", StatisticHtml, proxy.StatisticData())
}
