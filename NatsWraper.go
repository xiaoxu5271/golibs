package golibs

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Publish()：单纯的发布 Data 信息。
// PublishRequest()：带 Reply 的 Publish()。
// PublishMsg()：带 Header 的 PublishRequest()。
// Request()：同步版本的 PublishRequest()。
// RequestMsg()：高级版本的带 Header 的 Request()。
// Subscribe()：订阅感兴趣信息，收到信息后将调用回调函数。
// SubscribeSync()：同步版本的 Subscribe()。
// Respond()：通过 Reply 主题 Publish()。
// RespondMsg()：通过 Reply 主题 PublishMsg()。

// NATS 客户端.
type NatsWraper struct {
	requests Queue
	NConn    *nats.Conn
	client   *NClient
	NatsState

	server string
	mySubj string
	myName string

	recnnHdl ReconnectHandler
	recnnLck sync.Mutex
}

type NatsState struct {
	connected bool
	sync.RWMutex
}
type RequestMsg struct {
	Handler RespHandler
	Subj    string
	Data    StdHead
}

type ReconnectHandler func(client *NClient)
type RespHandler func(role, subj, topic string, content interface{}, err error)

func (n *NatsState) Write(state bool) {
	n.RWMutex.Lock()
	defer n.RWMutex.Unlock()
	n.connected = state
}

func (n *NatsState) Read() bool {
	n.RWMutex.RLock()
	defer n.RWMutex.RUnlock()
	return n.connected
}

// NatsWraper 构造方法.
// 返回 NatsWraper 对象指针.
func (n *NatsWraper) Init(server, mySubj, myName string, capacity int, client *NClient) *NatsWraper {
	n.NatsState.Write(false)

	if myName == "" {
		n.myName = "nclient"
	} else {
		n.myName = myName
	}
	if server == "" {
		n.server = "localhost:6000"
	} else {
		n.server = server
	}
	n.mySubj = mySubj

	n.client = client
	n.requests.Init(capacity)
	return n
}

func (n *NatsWraper) SetRecnnHandler(recnnHdl ReconnectHandler) {
	n.recnnLck.Lock()
	defer n.recnnLck.Unlock()
	n.recnnHdl = recnnHdl
}

func (n *NatsWraper) connect() {
	var e error

	for {
		// AllowReconnect:     true,
		// MaxReconnect:       DefaultMaxReconnect,
		// ReconnectWait:      DefaultReconnectWait,
		// ReconnectJitter:    DefaultReconnectJitter,
		// ReconnectJitterTLS: DefaultReconnectJitterTLS,
		// Timeout:            DefaultTimeout,
		// PingInterval:       DefaultPingInterval,
		// MaxPingsOut:        DefaultMaxPingOut,
		// SubChanLen:         DefaultMaxChanLen,
		// ReconnectBufSize:   1MB,
		// DrainTimeout:       DefaultDrainTimeout,
		n.NConn, e = nats.Connect(
			n.server,
			nats.ClosedHandler(n.onClose),
			nats.DisconnectHandler(n.onDisconnect),
			nats.ErrorHandler(n.onErr),
			nats.ReconnectHandler(n.onReconnect),
			nats.ReconnectBufSize(1024*1024),
			nats.Name(n.myName))
		if e == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	n.NatsState.Write(true)
}

// 关闭连接.
func (n *NatsWraper) Close() {
	if !n.NatsState.Read() {
		return
	}
	n.NatsState.Write(false)
	n.NConn.Close()
}

func (n *NatsWraper) IsConnected() bool {
	return n.NatsState.Read()
}

// 参数 subj 为接收使用的主题
func (n *NatsWraper) Subscribe(subj string) {
	_, e := n.NConn.Subscribe(subj, n.onSubject)
	if e != nil {
		panic(e)
	}
}

func (n *NatsWraper) SubscribeDefault() {
	n.Subscribe(n.mySubj)
}

func (n *NatsWraper) onSubject(m *nats.Msg) {
	n.client.TopicFilt.onSubject(m)
	e := m.Respond(m.Data)
	if e != nil && e != nats.ErrMsgNotBound && e != nats.ErrMsgNoReply {
		panic(e)
	}
}

// 请求 NATS 服务, 并等待返回处理结果.
// 参数 subj 为请求使用的主题.
// 参数 req 为解析请求数据.
// 参数 waiton 为等待的事件通道.
// 参数 timeout 为等待超时时间.
// 返回 rsp 为响应数据.

func (n *NatsWraper) Request(subj, topic string, content interface{}, hdl RespHandler, timeout time.Duration) bool {
	e := n.RequestAs(n.client.GetMySubj(), subj, topic, content, hdl, timeout)
	return e
}

func (n *NatsWraper) RequestAs(role, subj, topic string, content interface{}, hdl RespHandler, timeout time.Duration) bool {
	if !n.NatsState.Read() {
		return false
	}

	pub := StdHead{
		Src:     role,
		Topic:   topic,
		Content: content,
	}
	// 将待发送消息添加到 chan 内
	var req RequestMsg
	req.Subj = subj
	req.Data = pub
	req.Handler = hdl
	r := n.requests.EnQueue(req, timeout)
	return r
}

func (n *NatsWraper) requestService() {
	for {
		msgOut := n.requests.DeQueueBlock()
		req := msgOut.(RequestMsg)
		data, e := json.Marshal(req.Data)
		if e == nil {
			_, e = n.NConn.Request(req.Subj, data, 100*time.Millisecond)
		}
		hdl := req.Handler
		if hdl != nil {
			hdl(req.Data.Src, req.Subj, req.Data.Topic, req.Data.Content, e)
		}
	}
}

// 发布 NATS 消息.
// 参数 subj 为发布信息使用的主题.
// 参数 topic 为发布话题.
// 参数 content 为发布内容
func (n *NatsWraper) Publish(subj, topic string, content interface{}) bool {
	return n.PublishAs(n.client.GetMySubj(), subj, topic, content)
}

// 以 role 身份发布 NATS 消息.
// 参数 role 为指定的发布者身份.
// 参数 subj 为发布信息使用的主题.
// 参数 topic 为发布话题.
// 参数 content 为发布内容
func (n *NatsWraper) PublishAs(role, subj, topic string, content interface{}) bool {
	if !n.NatsState.Read() {
		return false
	}

	pub := StdHead{
		Src:     role,
		Topic:   topic,
		Content: content,
	}

	data, e := json.Marshal(pub)
	if e != nil {
		return false
	}

	e = n.NConn.Publish(subj, data)
	return e == nil
}

func (n *NatsWraper) onErr(*nats.Conn, *nats.Subscription, error) {
	// Do nothing
}

// 重连超时则断开连接
func (n *NatsWraper) onClose(nc *nats.Conn) {
	n.NatsState.Write(false)
	n.connect()

	n.recnnLck.Lock()
	defer n.recnnLck.Unlock()
	if n.recnnHdl != nil {
		n.recnnHdl(n.client)
	}
}

// 响应断开连接事件.
func (n *NatsWraper) onDisconnect(nc *nats.Conn) {
	n.NatsState.Write(false)
}

// 响应重连接事件.
func (n *NatsWraper) onReconnect(nc *nats.Conn) {
	n.NatsState.Write(true)

	n.recnnLck.Lock()
	defer n.recnnLck.Unlock()
	if n.recnnHdl != nil {
		n.recnnHdl(n.client)
	}
}
