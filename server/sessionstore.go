package main

import (
	"GoChat/server/store"
	"container/list"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tinode/chat/server/logs"
)

// session management

//sessionstore用老存储所有活跃的session
type SessionStore struct {
	lock sync.Mutex

	//用来管理长轮询
	lru      *list.List
	lifeTime time.Duration

	//用map存储所有的session
	sessCache map[string]*Session
}

//创建一个新的session
func (ss *SessionStore) NewSession(conn interface{}, sid string) (*Session, int) {
	var s Session

	if sid == "" {
		s.sid = store.GetUidString()
	} else {
		s.sid = sid
	}

	ss.lock.Lock()
	if _, found := ss.sessCache[s.sid]; found {
		logs.Err.Fatalln("ERROR! duplicate session ID", s.sid)
	}
	ss.lock.Unlock()
	switch c := conn.(type) { //sessiond的类型，websocket或则为longpoll
	case *websocket.Conn:
		s.proto = WEBSOCK
		s.ws = c
	case http.ResponseWriter:
		//long polling不需要进行存储，每次请求将会改变
		s.proto = LPOLL
	default:
		logs.Err.Panic("session:unkonwn connection type", conn)
	}
	s.subs = make(map[string]*Subscription)
	s.send = make(chan interface{}, sendQueueLimit+32)
	s.stop = make(chan interface{}, 1)
	s.detach = make(chan string, 64) //topic的名称

	s.bkgTimer = time.NewTimer(time.Hour)
	s.bkgTimer.Stop()

	s.inflightReqs = &sync.WaitGroup{}

	s.lastTouched = time.Now()

	ss.lock.Lock()
	if s.proto == LPOLL {
		s.lpTracker = ss.lru.PushFront(&s)
	}
	ss.sessCache[s.sid] = &s

	ss.lock.Unlock()
}
