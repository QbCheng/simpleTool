package simpleServerDiscovery

type ServiceDiscovery interface {
	NodeEvent() <-chan []EndpointEvent
	Close()
}
