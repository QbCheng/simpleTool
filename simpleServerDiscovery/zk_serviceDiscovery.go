package serviceDiscovery

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	defaultHeartbeatTimeout = time.Second * 30
	defaultRootPath         = "/simpleServerDiscovery/dynamicInstance"

	logPre = "[serviceDiscovery] [endpointId : %s] "
)

type cb func() error

type event struct {
	event zk.Event
	f     cb
}

type Option func(*serviceDiscovery)

// WithOpenDetailInfo 开启 go-zookeeper 所有日志信息. false 只开启错误日志
func WithOpenDetailInfo(openAll bool) Option {
	return func(sd *serviceDiscovery) {
		sd.zkLoggerAllOpen = openAll
	}
}

// WithLogger 指定日志
func WithLogger(logger zk.Logger) Option {
	return func(sd *serviceDiscovery) {
		sd.logger = logger
	}
}

// WithContext 指定上下文
func WithContext(ctx context.Context) Option {
	return func(sd *serviceDiscovery) {
		sd.ctx = ctx
	}
}

// WithHeartbeatTimeout 指定心跳
func WithHeartbeatTimeout(heartbeatTimeout time.Duration) Option {
	return func(sd *serviceDiscovery) {
		sd.heartbeatTimeout = heartbeatTimeout
	}
}

// WithEndpointEventNotify 接收事件通知
func WithEndpointEventNotify(nodeEventNotify bool) Option {
	return func(sd *serviceDiscovery) {
		sd.nodeEventNotify = nodeEventNotify
	}
}

// WithNodeContent 设置节点内容
func WithNodeContent(content string) Option {
	return func(sd *serviceDiscovery) {
		sd.nodeContent = content
	}
}

// WithRootPath 设置根地址
func WithRootPath(rootPath string) Option {
	return func(sd *serviceDiscovery) {
		sd.rootPath = rootPath
	}
}

type serviceDiscovery struct {
	conn      *zk.Conn
	eventChan <-chan zk.Event // zk 内部事件接收通道

	connUrl    []string // zk 连接地址
	endpointId string

	eventsLock sync.Mutex
	events     map[zk.EventType]map[string][]*event // zk.EventType 和 zk path 决定

	endpointsLock sync.RWMutex
	endpoints     map[string]*Endpoint
	nodeEventChan chan []EndpointEvent // 端点事件

	ctx              context.Context
	logger           zk.Logger
	zkLoggerAllOpen  bool
	heartbeatTimeout time.Duration
	nodeEventNotify  bool
	nodeContent      string
	rootPath         string
}

func NewServiceDiscovery(url []string, endpointId string, opts ...Option) (ServiceDiscovery, error) {
	c := &serviceDiscovery{
		connUrl:          url,
		heartbeatTimeout: defaultHeartbeatTimeout,
		events:           make(map[zk.EventType]map[string][]*event),
		nodeEventChan:    make(chan []EndpointEvent),
		endpointId:       endpointId,

		endpoints: map[string]*Endpoint{},

		ctx:             context.Background(),
		logger:          NewServerDiscoveryLogger(),
		zkLoggerAllOpen: false,
		nodeEventNotify: false,
		rootPath:        defaultRootPath,
	}
	for i := range opts {
		opts[i](c)
	}
	var err error
	c.conn, c.eventChan, err = zk.Connect(
		c.connUrl,
		c.heartbeatTimeout,
		zk.WithLogInfo(c.zkLoggerAllOpen),
		zk.WithLogger(c.logger),
	)
	if err != nil {
		return nil, err
	}

	// 将节点刷新事件进行注册
	c.register(&event{
		event: zk.Event{
			Type: zk.EventNodeChildrenChanged,
			Path: c.rootPath,
		},
		f: c.triggerAfterRegisterRefreshNodeEvent,
	})

	if err := c.advancedCreate(c.rootPath, 0, zk.WorldACL(zk.PermAll)); err != nil {
		return nil, err
	}

	go c.monitor(c.eventChan)
	return c, nil
}

func (s *serviceDiscovery) printf(err bool, layout string, str ...interface{}) {
	if s.zkLoggerAllOpen || err {
		s.logger.Printf(logPre+layout, append([]interface{}{s.endpointId}, str...)...)
	}
}

func (s *serviceDiscovery) error(layout string, str ...interface{}) {
	s.printf(true, layout, str...)
}

func (s *serviceDiscovery) info(layout string, str ...interface{}) {
	s.printf(false, layout, str...)
}

func (s *serviceDiscovery) monitor(eventChan <-chan zk.Event) {
	for {
		select {
		case event, ok := <-eventChan: // ZooKeeper连接事件
			if !ok {
				// <-eventChan 关闭
				return
			}
			if event.Path != "" {
				// 非连接事件
				s.handler(event)
			} else {
				// 连接事件
				if event.Type == zk.EventSession {
					if event.State == zk.StateHasSession {
						// 触发连接成功事件.
						err := s.connectSucceedEvent()
						if err != nil {
							s.error("Connect init fail. {err : %v}", err)
						}
					} else {
						s.info("{conn event trigger. non handler. {event : %#v}", event)
					}
				}
			}
		case <-s.ctx.Done():
			// 关闭连接
			s.Close()
			return
		}
	}
}

// 注册一个事件
func (s *serviceDiscovery) register(e *event) bool {
	s.eventsLock.Lock()
	defer s.eventsLock.Unlock()
	if s.events[e.event.Type] == nil {
		s.events[e.event.Type] = map[string][]*event{}
	}
	s.events[e.event.Type][e.event.Path] = append(s.events[e.event.Type][e.event.Path], e)
	return true
}

// 事件分发
func (s *serviceDiscovery) handler(e zk.Event) {
	if len(s.events) <= 0 {
		return
	}
	s.info("trigger event. {len : %v}. { event, %#v }", len(s.events[e.Type][e.Path]), e)
	s.eventsLock.Lock()
	defer s.eventsLock.Unlock()
	for i := range s.events[e.Type][e.Path] {
		err := s.events[e.Type][e.Path][i].f()
		if err != nil {
			s.error("trigger event. {id : %v}. { event, %#v }, {err : %v}", i, e, err)
		}
	}
}

// advancedCreate 能对传入的 zk path 进行层级创建
func (s *serviceDiscovery) advancedCreate(path string, flags int32, acl []zk.ACL) error {
	t := strings.Trim(path, "/")
	tArr := strings.Split(t, "/")

	tp := "/"
	for i := 0; i < len(tArr); i++ {
		tp += tArr[i]
		var err error
		_, err = s.conn.Create(tp, []byte{}, flags, acl)
		if err != nil {
			if !errors.Is(zk.ErrNodeExists, err) {
				return fmt.Errorf("s.conn.Create(%v), err : %v ", tp, err)
			}
		}
		tp += "/"
	}
	return nil
}

// 连接成功, 进行初始化
func (s *serviceDiscovery) connectSucceedEvent() error {
	err := s.registerSelf()
	if err != nil {
		return err
	}
	err = s.triggerAfterRegisterRefreshNodeEvent()
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceDiscovery) registerSelf() error {
	selfZkPath := s.rootPath + "/" + s.endpointId + "|FlagSequence|"
	var content []byte
	if s.nodeContent != "" {
		content = []byte(s.nodeContent)
	}
	_, err := s.conn.Create(selfZkPath, content, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

// 触发后. 重新注册事件
func (s *serviceDiscovery) triggerAfterRegisterRefreshNodeEvent() error {
	// 重新拉取全部 节点
	dynamicChildren, _, _, err := s.conn.ChildrenW(s.rootPath)
	if err != nil {
		return err
	}

	s.refreshNode(dynamicChildren)
	return nil
}

// 刷新在线服务.
func (s *serviceDiscovery) refreshNode(dynamicChildren []string) {
	s.endpointsLock.Lock()
	defer s.endpointsLock.Unlock()

	var endpointEvents []EndpointEvent
	active := map[string]interface{}{}
	for _, zkPath := range dynamicChildren {
		id := strings.Split(zkPath, "|FlagSequence|")[0]
		active[id] = struct{}{}
		if endpoint, ok := s.endpoints[id]; !ok {
			content, _, err := s.conn.Get(s.rootPath + "/" + zkPath)
			if err != nil {
				s.error("endpoint content failure to get. {path : %v} {err : %v}", zkPath, err)
				continue
			}
			s.endpoints[id] = &Endpoint{
				id:      id,
				state:   connect,
				content: content,
			}
			endpointEvents = append(endpointEvents, EndpointEvent{
				Endpoint: *s.endpoints[id],
				Et:       ETAdd,
			})
		} else {
			if endpoint.Connect() {
				continue
			}
			content, _, err := s.conn.Get(s.rootPath + "/" + zkPath)
			if err != nil {
				s.error("endpoint content failure to get. {path : %v}", zkPath)
				continue
			}
			// 开启节点
			endpoint.Change(connect)
			endpoint.SetContent(content)
			endpointEvents = append(endpointEvents, EndpointEvent{
				Endpoint: *s.endpoints[id],
				Et:       ETReconnect,
			})
		}
	}

	for _, v := range s.endpoints {
		if _, ok := active[v.id]; !ok {
			if v.Disconnect() {
				continue
			}
			v.Change(disconnect)
			endpointEvents = append(endpointEvents, EndpointEvent{
				Endpoint: *v,
				Et:       ETDisconnect,
			})
		}
	}

	if len(endpointEvents) > 0 && s.nodeEventNotify {
		s.nodeEventChan <- endpointEvents
	}

}

func (s *serviceDiscovery) NodeEvent() <-chan []EndpointEvent {
	return s.nodeEventChan
}

func (s *serviceDiscovery) Close() {
	s.eventsLock.Lock()
	defer s.eventsLock.Unlock()
	// 关闭连接
	s.conn.Close()
	// 清理全部自定义事件
	s.events = make(map[zk.EventType]map[string][]*event)
	// 关闭接收通道
	close(s.nodeEventChan)
	s.info("close success")
}
