package main

import (
	"GoChat/server/store/types"
	"container/list"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tinode/chat/pbx"
	"github.com/tinode/chat/server/auth"
)

const sendTimeout = time.Millisecond * 7
const sendQueueLimit = 128
const deferredNotificationsTimeout = time.Second * 5

var minSupportedVersionValue = parseVersion(minSupportedVersion)

//用户传输消息的通信协议
type SessionProto int

//定义传输协议的所有蕾西
const (
	NONE SessionProto = iota
	WEBSOCK
	LPOLL
	GRPC
	PROXY
	MULTIPLEX
)

// Session represents a single WS connection or a long polling session. A user may have multiple
// sessions.
type Session struct {
	// protocol - NONE (unset), WEBSOCK, LPOLL, GRPC, PROXY, MULTIPLEX
	proto SessionProto

	// Session ID
	sid string

	// Websocket. Set only for websocket sessions.
	ws *websocket.Conn

	// Pointer to session's record in sessionStore. Set only for Long Poll sessions.
	lpTracker *list.Element

	// gRPC handle. Set only for gRPC clients.
	grpcnode pbx.Node_MessageLoopServer

	//TODO: 集群session暂时不考虑
	// Reference to the cluster node where the session has originated. Set only for cluster RPC sessions.
	//clnode *ClusterNode

	//TODO:多个proxy暂时不考虑
	// Reference to multiplexing session. Set only for proxy sessions.
	// multi        *Session
	// proxiedTopic string

	// IP address of the client. For long polling this is the IP of the last poll.
	remoteAddr string

	// User agent, a string provived by an authenticated client in {login} packet.
	userAgent string

	// Protocol version of the client: ((major & 0xff) << 8) | (minor & 0xff).
	ver int

	// Device ID of the client
	deviceID string
	// Platform: web, ios, android
	platf string
	// Human language of the client
	lang string
	// Country code of the client
	countryCode string

	// ID of the current user. Could be zero if session is not authenticated
	// or for multiplexing sessions.
	uid types.Uid

	// Authentication level - NONE (unset), ANON, AUTH, ROOT.
	authLvl auth.Level

	// Time when the long polling session was last refreshed
	lastTouched time.Time

	// Time when the session received any packer from client
	lastAction int64

	// Background session: subscription presence notifications and online status are delayed.
	background bool
	// Timer which triggers after some seconds to mark background session as foreground.
	bkgTimer *time.Timer

	// Number of subscribe/unsubscribe requests in flight.
	inflightReqs *sync.WaitGroup
	// Synchronizes access to session store in cluster mode:
	// subscribe/unsubscribe replies are asynchronous.
	sessionStoreLock sync.Mutex
	// Indicates that the session is terminating.
	// After this flag's been flipped to true, there must not be any more writes
	// into the session's send channel.
	// Read/written atomically.
	// 0 = false
	// 1 = true
	terminating int32

	// Outbound mesages, buffered.
	// The content must be serialized in format suitable for the session.
	send chan interface{}

	// Channel for shutting down the session, buffer 1.
	// Content in the same format as for 'send'
	stop chan interface{}

	// detach - channel for detaching session from topic, buffered.
	// Content is topic name to detach from.
	detach chan string

	// Map of topic subscriptions, indexed by topic name.
	// Don't access directly. Use getters/setters.
	subs map[string]*Subscription
	// Mutex for subs access: both topic go routines and network go routines access
	// subs concurrently.
	subsLock sync.RWMutex

	// Needed for long polling and grpc.
	lock sync.Mutex

	// Field used only in cluster mode by topic master node.

	//TODO:集群设置，暂时不考虑
	// Type of proxy to master request being handled.
	//proxyReq ProxyReqType
}

//Subscription 用于表示session所订阅的topic
type Subscription struct {
	broadcast chan<- *ServerComMessage
	//当session取消订阅，向所有的topic发送一个消息
	done chan<- *sessionLeave
	meta chan<- *metaReq
	//当session更新时，向topic发送ping包
	supd chan<- *sessionUpdate
}
