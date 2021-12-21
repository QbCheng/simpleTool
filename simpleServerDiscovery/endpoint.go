package serviceDiscovery

type endpointType int

const (
	disconnect = iota + 1
	connect
)

// Endpoint 服务端点
type Endpoint struct {
	id      string
	state   endpointType
	content []byte
}

func (e *Endpoint) Id() string {
	return e.id
}

func (e *Endpoint) Connect() bool {
	return e.state == connect
}

func (e *Endpoint) Disconnect() bool {
	return e.state == disconnect
}

func (e *Endpoint) Content() []byte {
	return e.content
}

func (e *Endpoint) Change(state endpointType) {
	e.state = state
}

func (e *Endpoint) SetContent(content []byte) {
	e.content = content
}
