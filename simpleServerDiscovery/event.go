package serviceDiscovery

type EventType int

/*
事件类型
ETAdd : 增加一个新端点
ETReconnect : 端点重新连接
ETDisconnect : 端点断开连接
*/
const (
	ETAdd = iota + 1
	ETReconnect
	ETDisconnect
)

type EndpointEvent struct {
	Endpoint Endpoint
	Et       EventType
}
