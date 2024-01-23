package conn

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewManager)

type Manager struct {
	m      sync.RWMutex
	wsConn map[int64]map[int][]*Conn
	log    *log.Helper
}

func NewManager(logger log.Logger) *Manager {
	return &Manager{
		m:      sync.RWMutex{},
		wsConn: make(map[int64]map[int][]*Conn),
		log:    log.NewHelper(logger),
	}
}

func (m *Manager) AddConn(c *Conn) {
	m.m.Lock()
	defer m.m.Unlock()
	if oldConnMap, ok := m.wsConn[c.userId]; ok {
		if conns, ok := oldConnMap[c.platformId]; ok {
			conns = append(conns, c)
		} else {
			var conns []*Conn
			conns = append(conns, c)
			oldConnMap[c.platformId] = conns
		}
	} else {
		i := make(map[int][]*Conn)
		var conns []*Conn
		conns = append(conns, c)
		i[c.platformId] = conns
		m.wsConn[c.userId] = i
	}
	m.log.Info("add userId:%s platformId:%d", c.userId, c.platformId)
}

func (m *Manager) GetUserAllCons(uid int64) map[int][]*Conn {
	m.m.RLock()
	defer m.m.RUnlock()
	if connMap, ok := m.wsConn[uid]; ok {
		newConnMap := make(map[int][]*Conn)
		for k, v := range connMap {
			newConnMap[k] = v
		}
		return newConnMap
	}
	return nil
}

func (m *Manager) getUsers() []int64 {
	m.m.RLock()
	defer m.m.RUnlock()

	var userIds []int64
	for k, _ := range m.wsConn {
		userIds = append(userIds, k)
	}
	return userIds
}

func (m *Manager) DelConn(conn *Conn) {
	userAllCons := m.GetUserAllCons(conn.userId)

	m.m.Lock()
	defer m.m.Unlock()
	platform := conn.platformId
	for i, conns := range userAllCons {
		if platform == i {
			var newCoons []*Conn
			for _, c := range conns {
				if c != conn {
					newCoons = append(newCoons, c)
				}
			}
			if len(newCoons) != 0 {
				userAllCons[platform] = newCoons
			} else {
				delete(userAllCons, platform)
			}
		}
	}

	m.wsConn[conn.userId] = userAllCons
	if len(userAllCons) == 0 {
		delete(m.wsConn, conn.userId)
	}

	m.log.Info("delete userId:%s platformId:%d", conn.userId, conn.platformId)

	conn.GetWsConn().Close()
}
