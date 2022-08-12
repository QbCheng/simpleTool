package serviceDiscovery

type ServiceDiscovery interface {
	NodeEvent() <-chan []EndpointEvent
	Close()
}
