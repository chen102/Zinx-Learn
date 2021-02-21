package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	Connections map[uint32]ziface.IConnection
	ConnLock    sync.RWMutex //保护连接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
}
func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()
	cm.Connections[conn.GetConnID()] = conn
	fmt.Println("connID=", conn.GetConnID(), " add to ConnManager successfully:conn num =", cm.Len())
}
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()
	delete(cm.Connections, conn.GetConnID())
	fmt.Println("connID=", conn.GetConnID(), " remove from ConnManager successfully:conn num =", cm.Len())
}
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.ConnLock.RLock()
	defer cm.ConnLock.RUnlock()
	if conn, ok := cm.Connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not Found!")
	}
}
func (cm *ConnManager) Len() int {
	return len(cm.Connections)

}
func (cm *ConnManager) ClearConn() {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()
	for connID, conn := range cm.Connections {
		conn.Stop()
		delete(cm.Connections, connID)
	}
	fmt.Println("Clear All Connections Succ! conn num =", cm.Len())

}
